package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Azure/go-autorest/autorest"
	"github.com/manicminer/hamilton/auth"
)

type Authorizer struct {
	auth.Authorizer
}

// WithAuthorization implements the autorest.Authorizer interface
func (c *Authorizer) WithAuthorization() autorest.PrepareDecorator {
	return func(p autorest.Preparer) autorest.Preparer {
		return autorest.PreparerFunc(func(req *http.Request) (*http.Request, error) {
			var err error
			req, err = p.Prepare(req)
			if err == nil {
				token, err := c.Token()
				if err != nil {
					return nil, err
				}

				req, err = autorest.Prepare(req, autorest.WithHeader("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken)))
				if err != nil {
					return req, err
				}

				auxTokens, err := c.AuxiliaryTokens()
				if err != nil {
					return req, err
				}

				auxTokenList := make([]string, 0)
				for _, a := range auxTokens {
					if a != nil && a.AccessToken != "" {
						auxTokenList = append(auxTokenList, fmt.Sprintf("%s %s", a.TokenType, a.AccessToken))
					}
				}

				return autorest.Prepare(req, autorest.WithHeader("x-ms-authorization-auxiliary", strings.Join(auxTokenList, ", ")))
			}

			return req, err
		})
	}
}

// BearerAuthorizerCallback is a helper that returns an *autorest.BearerAuthorizerCallback for use in data plane API clients in the Azure SDK
func (c *Authorizer) BearerAuthorizerCallback() *autorest.BearerAuthorizerCallback {
	return autorest.NewBearerAuthorizerCallback(nil, func(_, resource string) (*autorest.BearerAuthorizer, error) {
		token, err := c.Token()
		if err != nil {
			return nil, fmt.Errorf("obtaining token: %v", err)
		}

		return autorest.NewBearerAuthorizer(&servicePrincipalTokenWrapper{
			tokenType:  "Bearer",
			tokenValue: token.AccessToken,
		}), nil
	})
}
