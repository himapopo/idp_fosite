package oauth

import (
	"context"
	"log"

	"github.com/ory/fosite"
)

type Client struct {
	ID            string
	HashedSecret  []byte
	RedirectURIs  []string
	GrantTypes    fosite.Arguments
	ResponseTypes fosite.Arguments
	Scopes        fosite.Arguments
	AuthMethod    string
}

var (
	Client1 *Client
	Client2 *Client
)

func init() {
	hasher := &fosite.BCrypt{
		Config: config,
	}
	secret1, err := hasher.Hash(context.Background(), []byte("client1"))
	if err != nil {
		log.Fatalln(err)
	}
	Client1 = &Client{
		ID:            "1",
		HashedSecret:  secret1,
		RedirectURIs:  []string{"http://localhost:3846/callback"},
		GrantTypes:    fosite.Arguments{"authorization_code", "refresh_token"},
		ResponseTypes: fosite.Arguments{"code"},
		Scopes:        fosite.Arguments{"openid", "offline"},
		AuthMethod:    "client_secret_post",
	}

	secret2, _ := hasher.Hash(context.Background(), []byte("client2"))
	Client2 = &Client{
		ID:            "2",
		HashedSecret:  secret2,
		RedirectURIs:  []string{"http://localhost:3846/callback"},
		GrantTypes:    fosite.Arguments{"client_credentials"},
		ResponseTypes: fosite.Arguments{},
		Scopes:        fosite.Arguments{"sample", "offline"},
		AuthMethod:    "client_secret_post",
	}

}

func (c *Client) GetID() string {
	return c.ID
}

func (c *Client) IsPublic() bool {
	return false
}

func (c *Client) GetAudience() fosite.Arguments {
	return nil
}

func (c *Client) GetHashedSecret() []byte {
	return c.HashedSecret
}

func (c *Client) GetRedirectURIs() []string {
	return c.RedirectURIs
}

func (c *Client) GetGrantTypes() fosite.Arguments {
	return c.GrantTypes
}

func (c *Client) GetResponseTypes() fosite.Arguments {
	return c.ResponseTypes
}

func (c *Client) GetScopes() fosite.Arguments {
	return c.Scopes
}

func (c *Client) GetTokenEndpointAuthMethod() string {
	return c.AuthMethod
}
