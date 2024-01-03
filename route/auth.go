package route

import (
	"idp_fosite/oauth"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/token/jwt"
)

func authenticate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ar, err := oauth.OAuth2Provider.NewAuthorizeRequest(ctx, r)
	if err != nil {
		oauth.OAuth2Provider.WriteAuthorizeError(ctx, w, ar, err)
		return
	}

	sess := &openid.DefaultSession{
		Claims: &jwt.IDTokenClaims{
			Issuer:      "",
			Subject:     "taiki",
			Audience:    []string{""},
			ExpiresAt:   time.Now().Add(time.Hour * 6),
			IssuedAt:    time.Now(),
			RequestedAt: time.Now(),
			AuthTime:    time.Now(),
		},
		Headers: &jwt.Headers{
			Extra: make(map[string]interface{}),
		},
	}

	// MEMO: fositeでやってくれてもいい気がする
	scopes := strings.Split(r.URL.Query().Get("scope"), " ")
	for _, s := range scopes {
		if len(s) > 0 {
			ar.GrantScope(s)
		}
	}

	response, err := oauth.OAuth2Provider.NewAuthorizeResponse(ctx, ar, sess)
	if err != nil {
		oauth.OAuth2Provider.WriteAuthorizeError(ctx, w, ar, err)
		return
	}

	oauth.OAuth2Provider.WriteAuthorizeResponse(ctx, w, ar, response)
	return
}

func token(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sess := &openid.DefaultSession{
		Claims: &jwt.IDTokenClaims{
			Issuer:      "",
			Subject:     "taiki",
			Audience:    []string{""},
			ExpiresAt:   time.Now().Add(time.Hour * 6),
			IssuedAt:    time.Now(),
			RequestedAt: time.Now(),
			AuthTime:    time.Now(),
		},
		Headers: &jwt.Headers{
			Extra: make(map[string]interface{}),
		},
	}

	accessRequest, err := oauth.OAuth2Provider.NewAccessRequest(ctx, r, sess)

	if err != nil {
		log.Printf("Error occurred in NewAccessRequest: %+v", err)
		oauth.OAuth2Provider.WriteAccessError(ctx, w, accessRequest, err)
		return
	}

	if accessRequest.GetGrantTypes().ExactOne("client_credentials") {
		for _, scope := range accessRequest.GetRequestedScopes() {
			accessRequest.GrantScope(scope)
		}
	}

	response, err := oauth.OAuth2Provider.NewAccessResponse(ctx, accessRequest)
	if err != nil {
		log.Printf("Error occurred in NewAccessResponse: %+v", err)
		oauth.OAuth2Provider.WriteAccessError(ctx, w, accessRequest, err)
		return
	}

	oauth.OAuth2Provider.WriteAccessResponse(ctx, w, accessRequest, response)
}

func introspect(w http.ResponseWriter, r *http.Request) {

}

func logout(w http.ResponseWriter, r *http.Request) {

}

func userInfo(w http.ResponseWriter, r *http.Request) {

}

func revoke(w http.ResponseWriter, r *http.Request) {

}
