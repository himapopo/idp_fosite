package oauth

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/openid"
)

var db sync.Map

type OAuth2ProvidorApp struct {
}

func (a OAuth2ProvidorApp) GetClient(ctx context.Context, id string) (fosite.Client, error) {
	if id == "1" {
		return Client1, nil
	}

	if id == "2" {
		return Client2, nil
	}

	return nil, errors.New("client not found")
}

func (a OAuth2ProvidorApp) CreateOpenIDConnectSession(ctx context.Context, authorizeCode string, requester fosite.Requester) error {
	sess, ok := requester.GetSession().(*openid.DefaultSession)
	if !ok {
		return errors.New("not openid.DefaultSession")
	}
	db.Store(sess.Claims.Subject, sess)
	db.Store(sess.Claims.Subject+"authorizeCode", authorizeCode)
	return nil
}

func (a OAuth2ProvidorApp) GetOpenIDConnectSession(ctx context.Context, authorizeCode string, requester fosite.Requester) (fosite.Requester, error) {
	// TODO: 実装
	return nil, nil
}

func (a OAuth2ProvidorApp) DeleteOpenIDConnectSession(ctx context.Context, authorizeCode string) error {
	// TODO: 実装
	return nil
}

func (a OAuth2ProvidorApp) CreateAuthorizeCodeSession(ctx context.Context, code string, request fosite.Requester) error {
	sess, ok := request.GetSession().(*openid.DefaultSession)
	if !ok {
		return errors.New("not openid.DefaultSession")
	}
	db.Store(sess.Claims.Subject, sess)
	db.Store(sess.Claims.Subject+"authorizeCode", code)
	return nil
}

func (a OAuth2ProvidorApp) GetAuthorizeCodeSession(ctx context.Context, code string, session fosite.Session) (fosite.Requester, error) {
	// TODO: 実装
	return nil, nil
}

func (a OAuth2ProvidorApp) InvalidateAuthorizeCodeSession(ctx context.Context, code string) error {
	// TODO: 実装
	return nil
}

func (a OAuth2ProvidorApp) CreateAccessTokenSession(ctx context.Context, signature string, request fosite.Requester) error {
	// TODO: 実装
	return nil
}

func (a OAuth2ProvidorApp) CreateRefreshTokenSession(ctx context.Context, signature string, request fosite.Requester) error {
	// TODO: 実装
	return nil
}

func (a OAuth2ProvidorApp) GetAccessTokenSession(ctx context.Context, signature string, session fosite.Session) (fosite.Requester, error) {
	// TODO: 実装
	return nil, nil
}

func (a OAuth2ProvidorApp) DeleteAccessTokenSession(ctx context.Context, signature string) error {
	// TODO: 実装
	return nil
}

func (a OAuth2ProvidorApp) GetRefreshTokenSession(ctx context.Context, signature string, session fosite.Session) (fosite.Requester, error) {
	// TODO: 実装
	return nil, nil
}

func (a OAuth2ProvidorApp) DeleteRefreshTokenSession(ctx context.Context, signature string) error {
	// TODO: 実装
	return nil
}

// implements fosite.ClientManager
func (a OAuth2ProvidorApp) ClientAssertionJWTValid(_ context.Context, jti string) error {
	// TODO: jtiの値を見て信頼できるものかを返す。
	return nil
}

// SetClientAssertionJWT
// implements fosite.ClientManager
func (a OAuth2ProvidorApp) SetClientAssertionJWT(_ context.Context, jti string, exp time.Time) error {
	// TODO: jtiの有効期限を保存する。
	return nil
}

func (a OAuth2ProvidorApp) RevokeRefreshToken(ctx context.Context, requestID string) error {
	// TODO: 実装
	return nil
}

func (a OAuth2ProvidorApp) RevokeRefreshTokenMaybeGracePeriod(ctx context.Context, requestID string, signature string) error {
	// TODO: 実装
	return nil
}

func (a OAuth2ProvidorApp) RevokeAccessToken(ctx context.Context, requestID string) error {
	// TODO: 実装
	return nil
}
