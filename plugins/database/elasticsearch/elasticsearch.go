package elasticsearch

import (
	"context"
	"errors"
	"time"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/logical/database/dbplugin"
	"github.com/hashicorp/vault/plugins"
)

// Elasticsearch implements dbplugin's Database interface
type Elasticsearch struct {
	Logger hclog.Logger
	Client *Client
}

func New() (interface{}, error) {
	return &Elasticsearch{
		Logger: hclog.Default(), // TODO what if Vault is on Debug? Will this pick it up?
	}, nil
}

func Run(apiTLSConfig *api.TLSConfig) error {
	dbType, err := New()
	if err != nil {
		return err
	}
	plugins.Serve(dbType.(dbplugin.Database), apiTLSConfig)
	return nil
}

func (es *Elasticsearch) Type() (string, error) {
	return "elasticsearch", nil
}

func (es *Elasticsearch) Init(ctx context.Context, config map[string]interface{}, verifyConnection bool) (map[string]interface{}, error) {
	var esURL string
	if raw, ok := config["url"]; ok {
		esURL, ok = raw.(string)
		if !ok {
			return nil, errors.New(`"url" must be a string`)
		}
	} else {
		return nil, errors.New(`"url" must be provided`)
	}

	var username string
	if raw, ok := config["username"]; ok {
		username, ok = raw.(string)
		if !ok {
			return nil, errors.New(`"username" must be a string`)
		}
	} else {
		return nil, errors.New(`"username" must be provided`)
	}

	var password string
	if raw, ok := config["password"]; ok {
		password, ok = raw.(string)
		if !ok {
			return nil, errors.New(`"password" must be a string"`)
		}
	} else {
		return nil, errors.New(`"password" must be provided`)
	}

	tlsConfigProvided := false
	tlsConfig := &TLSConfig{}

	if raw, ok := config["ca_cert"]; ok {
		tlsConfig.CACert, ok = raw.(string)
		if !ok {
			return nil, errors.New(`"ca_cert" must be a string`)
		}
		tlsConfigProvided = true
	}
	if raw, ok := config["ca_path"]; ok {
		tlsConfig.CAPath, ok = raw.(string)
		if !ok {
			return nil, errors.New(`"ca_path" must be a string`)
		}
		tlsConfigProvided = true
	}
	if raw, ok := config["client_cert"]; ok {
		tlsConfig.ClientCert, ok = raw.(string)
		if !ok {
			return nil, errors.New(`"client_cert" must be a string`)
		}
		tlsConfigProvided = true
	}
	if raw, ok := config["client_key"]; ok {
		tlsConfig.ClientKey, ok = raw.(string)
		if !ok {
			return nil, errors.New(`"client_key" must be a string`)
		}
		tlsConfigProvided = true
	}
	if raw, ok := config["tls_server_name"]; ok {
		tlsConfig.TLSServerName, ok = raw.(string)
		if !ok {
			return nil, errors.New(`"tls_server_name" must be a string`)
		}
		tlsConfigProvided = true
	}
	if raw, ok := config["insecure"]; ok {
		tlsConfig.Insecure, ok = raw.(bool)
		if !ok {
			return nil, errors.New(`"insecure" must be a bool`)
		}
		tlsConfigProvided = true
	}

	if tlsConfigProvided {
		client, err := NewTLSClient(ctx.Done(), es.Logger, username, password, esURL, tlsConfig)
		if err != nil {
			return nil, err
		}
		es.Client = client
	} else {
		client, err := NewClient(ctx.Done(), es.Logger, username, password, esURL)
		if err != nil {
			return nil, err
		}
		es.Client = client
	}

	if verifyConnection {
		// Whether this role is found or unfound, if we're configured correctly there will
		// be no err from the client. However, if something is misconfigured, this will yield
		// an error response, which will be described in the returned error.
		if _, err := es.Client.GetRole("vault-test"); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (es *Elasticsearch) CreateUser(ctx context.Context, statements dbplugin.Statements, usernameConfig dbplugin.UsernameConfig, expiration time.Time) (username string, password string, err error) {
	// TODO
	return "", "", nil
}

func (es *Elasticsearch) RenewUser(ctx context.Context, statements dbplugin.Statements, username string, expiration time.Time) error {
	// TODO
	return nil
}

func (es *Elasticsearch) RevokeUser(ctx context.Context, statements dbplugin.Statements, username string) error {
	// TODO
	return nil
}

func (es *Elasticsearch) RotateRootCredentials(ctx context.Context, statements []string) (config map[string]interface{}, err error) {
	// TODO
	return nil, nil
}

// This gets called after calling roles or creds.
func (es *Elasticsearch) Close() error {
	// TODO
	return nil
}

// DEPRECATED, included for backward-compatibility until removal
func (es *Elasticsearch) Initialize(ctx context.Context, config map[string]interface{}, verifyConnection bool) error {
	_, err := es.Init(ctx, config, verifyConnection)
	return err
}
