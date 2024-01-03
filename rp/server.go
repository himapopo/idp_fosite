package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/callback", func(w http.ResponseWriter, r *http.Request) {
		token := tokenReqest(r)
		fmt.Println("-------token------")
		fmt.Println(token)
		w.Write([]byte(token))
		return
	})

	http.ListenAndServe(":3846", r)
}

func tokenReqest(r *http.Request) string {
	code := r.URL.Query().Get("code")
	// 送信するデータを作成
	data := url.Values{}
	data.Add("grant_type", "authorization_code")
	data.Add("code", code)
	data.Add("redirect_uri", "http://localhost:3846/callback")
	data.Add("client_id", "1")
	data.Add("client_secret", "client1")

	req, err := http.NewRequest("POST", "http://localhost:3000/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return ""
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// HTTPクライアントを作成してリクエストを送信
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return ""
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)

	return string(b)
}
