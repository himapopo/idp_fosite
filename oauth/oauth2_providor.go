package oauth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"

	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
	"github.com/ory/fosite/token/jwt"
)

var secret = []byte("my super secret signing password")

var privateKey, _ = rsa.GenerateKey(rand.Reader, 2048)

var config = &fosite.Config{
	GlobalSecret: secret,
}

var keyGetter = func(context.Context) (interface{}, error) {
	return privateKey, nil
}

var strategy = &compose.CommonStrategy{
	CoreStrategy:               compose.NewOAuth2HMACStrategy(config),
	OpenIDConnectTokenStrategy: compose.NewOpenIDConnectStrategy(keyGetter, config),
	Signer:                     &jwt.DefaultSigner{GetPrivateKey: keyGetter},
}

var OAuth2Provider = compose.Compose(
	config,
	OAuth2ProvidorApp{},
	strategy,
	compose.OAuth2AuthorizeExplicitFactory,
	compose.OAuth2ClientCredentialsGrantFactory,
	compose.OpenIDConnectExplicitFactory,
	compose.OAuth2TokenIntrospectionFactory,
)
