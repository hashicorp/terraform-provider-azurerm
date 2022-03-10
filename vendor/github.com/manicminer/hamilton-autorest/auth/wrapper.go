package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/manicminer/hamilton/auth"
	"golang.org/x/oauth2"
)

// NewAuthorizerWrapper returns a Hamilton auth.Authorizer that sources tokens from a supplied autorest.BearerAuthorizer
func NewAuthorizerWrapper(autorestAuthorizer autorest.Authorizer) (auth.Authorizer, error) {
	return &AuthorizerWrapper{authorizer: autorestAuthorizer}, nil
}

// AuthorizerWrapper is a Hamilton auth.Authorizer which sources tokens from an autorest.Authorizer
// Fully supports:
// - autorest.BearerAuthorizer
// - autorest.MultiTenantBearerAuthorizer
// For other types that satisfy autorest.Authorizer, the Authorization and X-Ms-Authorization-Auxiliary headers
// are parsed for access token values, but additional metadata such as refresh tokens and expiry are not provided.
type AuthorizerWrapper struct {
	authorizer autorest.Authorizer
}

type servicePrincipalTokenWrapper struct {
	tokenType  string
	tokenValue string
}

func (s *servicePrincipalTokenWrapper) OAuthToken() string {
	return s.tokenValue
}

func (s *servicePrincipalTokenWrapper) Token() adal.Token {
	return adal.Token{
		AccessToken: s.tokenValue,
		Type:        s.tokenType,
	}
}

type ServicePrincipalToken interface {
	Token() adal.Token
}

func (a *AuthorizerWrapper) tokenProviders() (tokenProviders []adal.OAuthTokenProvider, err error) {
	if authorizer, ok := a.authorizer.(*autorest.BearerAuthorizer); ok && authorizer != nil {
		// autorest.BearerAuthorizer provides a single token for the specified tenant
		tokenProviders = append(tokenProviders, authorizer.TokenProvider())
	} else if authorizer, ok := a.authorizer.(*autorest.MultiTenantBearerAuthorizer); ok && authorizer != nil {
		// autorest.MultiTenantBearerAuthorizer provides tokens for the primary specified
		// tenant plus any specified auxiliary tenants
		if multiTokenProvider := authorizer.TokenProvider(); multiTokenProvider != nil {
			if m, ok := multiTokenProvider.(*adal.MultiTenantServicePrincipalToken); ok && m != nil {
				tokenProviders = append(tokenProviders, m.PrimaryToken)
				for _, aux := range m.AuxiliaryTokens {
					tokenProviders = append(tokenProviders, aux)
				}
			}
		}
	} else {
		// a generic autorest.Authorizer only supplies HTTP headers, so we'll have
		// to parse those to obtain a token
		req, err := autorest.Prepare(&http.Request{}, a.authorizer.WithAuthorization())
		if err != nil {
			return nil, err
		}

		// first parse out the Authorization header to get a token type and value
		if authorization := strings.SplitN(req.Header.Get("Authorization"), " ", 2); len(authorization) == 2 {
			tokenProviders = append(tokenProviders, &servicePrincipalTokenWrapper{
				tokenType:  authorization[0],
				tokenValue: authorization[1],
			})
		}

		// next parse out any comma-separated auxiliary tokens to get their token type and value
		if authorizationAux := strings.Split(req.Header.Get("X-Ms-Authorization-Auxiliary"), ","); len(authorizationAux) > 0 {
			for _, authorizationRaw := range authorizationAux {
				if authorization := strings.SplitN(strings.TrimSpace(authorizationRaw), " ", 2); len(authorization) == 2 {
					tokenProviders = append(tokenProviders, &servicePrincipalTokenWrapper{
						tokenType:  authorization[0],
						tokenValue: authorization[1],
					})
				}
			}
		}
	}

	for _, tokenProvider := range tokenProviders {
		if refresher, ok := tokenProvider.(adal.Refresher); ok {
			if err = refresher.EnsureFresh(); err != nil {
				return
			}
		} else if refresher, ok := tokenProvider.(adal.RefresherWithContext); ok {
			if err = refresher.EnsureFreshWithContext(context.Background()); err != nil {
				return
			}
		}
	}

	return
}

// Token returns an access token using an autorest.BearerAuthorizer struct
func (a *AuthorizerWrapper) Token() (*oauth2.Token, error) {
	tokenProviders, err := a.tokenProviders()
	if err != nil {
		return nil, err
	}
	if len(tokenProviders) == 0 {
		return nil, fmt.Errorf("no token providers returned")
	}

	var adalToken adal.Token
	if spToken, ok := tokenProviders[0].(ServicePrincipalToken); ok && spToken != nil {
		adalToken = spToken.Token()
	}
	if adalToken.AccessToken == "" {
		return nil, fmt.Errorf("could not obtain access token from token provider")
	}

	return &oauth2.Token{
		AccessToken:  adalToken.AccessToken,
		TokenType:    adalToken.Type,
		RefreshToken: adalToken.RefreshToken,
		Expiry:       adalToken.Expires(),
	}, nil
}

// AuxiliaryTokens returns additional tokens for auxiliary tenant IDs, sourced from an
// autorest.MultiTenantBearerAuthorizer, for use in multi-tenant scenarios
func (a *AuthorizerWrapper) AuxiliaryTokens() ([]*oauth2.Token, error) {
	tokenProviders, err := a.tokenProviders()
	if err != nil {
		return nil, err
	}

	var auxTokens []*oauth2.Token
	for i := 1; i < len(tokenProviders); i++ {
		var adalToken adal.Token

		if spToken, ok := tokenProviders[i].(ServicePrincipalToken); ok && spToken != nil {
			adalToken = spToken.Token()
		}

		if adalToken.AccessToken == "" {
			return nil, fmt.Errorf("could not obtain access token from token providers")
		}

		auxTokens = append(auxTokens, &oauth2.Token{
			AccessToken:  adalToken.AccessToken,
			TokenType:    adalToken.Type,
			RefreshToken: adalToken.RefreshToken,
			Expiry:       adalToken.Expires(),
		})
	}

	return auxTokens, nil
}
