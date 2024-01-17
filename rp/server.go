package main

import (
	"client/client"
	"encoding/base64"
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
)

var stateDB = &sync.Map{}

type TokenRes struct {
	client.Client
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	Expire       int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/client", func(r chi.Router) {
		// ホーム画面
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				c := client.GetClient(chi.URLParam(r, "id"))
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
				c := client.GetClient(chi.URLParam(r, "id"))
				if c == nil {
					panic("client not found")
				}
				uid := uuid.New().String()
				stateDB.Store(uid, chi.URLParam(r, "id"))

				url := fmt.Sprintf("http://localhost:3000/oauth2/auth?response_type=code&client_id=%s&state=%s&redirect_uri=%s&scope=%s", c.ID, uid, c.RedirectURIs[0], strings.Join(c.Scopes, "+"))

				res, _ := http.Get(url)

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

				renderHTML(w, "token.html", tokenRes)
			})

			// クライアントクレデンシャルズフロー
			r.Get("/client_credentials_request", func(w http.ResponseWriter, r *http.Request) {
				c := client.GetClient(chi.URLParam(r, "id"))
				if c == nil {
					panic("client not found")
				}

				data := url.Values{}
				data.Add("grant_type", "client_credentials")
				data.Add("scope", "sample offline")

				b := tokenRequest(r, c, data)

				var tokenRes TokenRes
				if err := json.Unmarshal(b, &tokenRes); err != nil {
					panic(err)
				}

				tokenRes.ID = c.ID
				tokenRes.Name = c.Name

				renderHTML(w, "token.html", tokenRes)
			})
		})
	})

	// 内部でトークンリクエスト
	r.Get("/callback", func(w http.ResponseWriter, r *http.Request) {

		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")

		cid, _ := stateDB.Load(state)
		stateDB.Delete(state)

		id, _ := cid.(string)
		c := client.GetClient(id)
		if c == nil {
			panic("client not found")
		}

		data := url.Values{}
		data.Add("grant_type", "authorization_code")
		data.Add("code", code)
		data.Add("redirect_uri", "http://localhost:3846/callback")

		b := tokenRequest(r, c, data)

		w.Write(b)
		return
	})

	http.ListenAndServe(":3846", r)
}

func tokenRequest(r *http.Request, c *client.Client, values url.Values) []byte {
	req, err := http.NewRequest("POST", "http://localhost:3000/oauth2/token", strings.NewReader(values.Encode()))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil
	}

	src := []byte(c.ID + ":" + c.Secret)
	enc := base64.StdEncoding.EncodeToString(src)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+enc)

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

func renderHTML(w http.ResponseWriter, file string, data any) {
	tmp := template.Must(template.ParseFiles(file))
	if err := tmp.Execute(w, data); err != nil {
		panic(err)
	}
}
