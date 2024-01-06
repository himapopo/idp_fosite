package client

import "github.com/ory/fosite"

type Client struct {
	ID            string
	Name          string
	Secret        string
	RedirectURIs  []string
	GrantTypes    fosite.Arguments
	ResponseTypes fosite.Arguments
	Scopes        fosite.Arguments
	AuthMethod    string
}

var (
	Client1 = &Client{
		ID:           "1",
		Name:         "太紀のアプリケーション",
		Secret:       "client1",
		RedirectURIs: []string{"http://localhost:3846/callback"},
		Scopes:       fosite.Arguments{"openid", "offline"},
		GrantTypes:   fosite.Arguments{"authorization_code", "refresh_token"},
	}
	Client2 = &Client{
		ID:           "2",
		Name:         "クライアントクレデンシャルズアプリケーション",
		Secret:       "client2",
		RedirectURIs: []string{"http://localhost:3846/callback"},
		Scopes:       fosite.Arguments{"sample", "offline"},
		GrantTypes:   fosite.Arguments{"client_credentials"},
	}

	Clients = []*Client{Client1, Client2}
)

func GetClient(id string) *Client {
	for _, c := range Clients {
		if id == c.ID {
			return c
		}
	}
	return nil
}
