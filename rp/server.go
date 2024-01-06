package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"text/template"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
	"github.com/ory/fosite"
)

var stateDB = &sync.Map{}

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

type TokenRes struct {
	Client
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	Expire       int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

var (
	Client1 = &Client{
		ID:           "1",
		Name:         "太紀のアプリケーション",
		Secret:       "client1",
		RedirectURIs: []string{"http://localhost:3846/callback"},
		Scopes:       fosite.Arguments{"openid"},
	}
	Client2 *Client

	Clients = []*Client{Client1, Client2}
)

func getClient(id string) *Client {
	for _, c := range Clients {
		if id == c.ID {
			return c
		}
	}
	return nil
}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/client", func(r chi.Router) {
		// ホーム画面
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				c := getClient(chi.URLParam(r, "id"))
				if c == nil {
					panic("client not found")
				}
				tmp := template.Must(template.ParseFiles("home.html"))
				if err := tmp.Execute(w, *c); err != nil {
					panic(err)
				}
			})

			// 認可コードリクエスト
			// トークン受け取って表示
			r.Get("/auth_request", func(w http.ResponseWriter, r *http.Request) {
				c := getClient(chi.URLParam(r, "id"))
				if c == nil {
					panic("client not found")
				}
				uid := uuid.New().String()
				stateDB.Store(uid, chi.URLParam(r, "id"))

				url := fmt.Sprintf("http://localhost:3000/oauth2/auth?response_type=code&client_id=%s&state=%s&redirect_uri=%s&scope=%s", c.ID, uid, c.RedirectURIs[0], strings.Join(c.Scopes, " "))

				res, _ := http.Get(url)
				// if err != nil {
				// 	panic(err)
				// }

				if res == nil {
					return
				}
				defer res.Body.Close()

				b, err := io.ReadAll(res.Body)
				if err != nil {
					panic(err)
				}

				var tokenRes TokenRes
				if err := json.Unmarshal(b, &tokenRes); err != nil {
					panic(err)
				}

				tokenRes.ID = c.ID
				tokenRes.Name = c.Name

				tmp := template.Must(template.ParseFiles("token.html"))
				if err := tmp.Execute(w, tokenRes); err != nil {
					panic(err)
				}
				return
			})
		})
	})

	// 内部でトークンリクエスト
	r.Get("/callback", func(w http.ResponseWriter, r *http.Request) {
		b := tokenRequest(r)
		w.Write(b)
		return
	})

	http.ListenAndServe(":3846", r)
}

func tokenRequest(r *http.Request) []byte {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	cid, _ := stateDB.Load(state)
	id, _ := cid.(string)
	c := getClient(id)
	if c == nil {
		panic("client not found")
	}

	data := url.Values{}
	data.Add("grant_type", "authorization_code")
	data.Add("code", code)
	data.Add("redirect_uri", "http://localhost:3846/callback")
	data.Add("client_id", c.ID)
	data.Add("client_secret", c.Secret)

	req, err := http.NewRequest("POST", "http://localhost:3000/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)

	return b
}
