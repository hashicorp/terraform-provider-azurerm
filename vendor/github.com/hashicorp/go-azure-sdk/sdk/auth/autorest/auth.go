// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package autorest

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/hashicorp/go-azure-sdk/sdk/auth"
)

func AutorestAuthorizer(authorizer auth.Authorizer) *Authorizer {
	return &Authorizer{Authorizer: authorizer}
}

type Authorizer struct {
	auth.Authorizer
}

// WithAuthorization implements the autorest.Authorizer interface
func (c *Authorizer) WithAuthorization() autorest.PrepareDecorator {
	return func(p autorest.Preparer) autorest.Preparer {
		return autorest.PreparerFunc(func(req *http.Request) (*http.Request, error) {
			ctx := req.Context()
			var err error
			req, err = p.Prepare(req)
			if err == nil {
				token, err := c.Token(ctx, req)
				if err != nil {
					return nil, err
				}

				req, err = autorest.Prepare(req, autorest.WithHeader("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken)))
				if err != nil {
					return nil, fmt.Errorf("preparing request: %+v", err)
				}

				auxTokens, err := c.AuxiliaryTokens(ctx, req)
				if err != nil {
					return nil, fmt.Errorf("preparing auxiliary tokens for request: %+v", err)
				}
				if len(auxTokens) > 0 {
					auxTokenList := make([]string, 0)
					for _, a := range auxTokens {
						if a != nil && a.AccessToken != "" {
							auxTokenList = append(auxTokenList, fmt.Sprintf("%s %s", a.TokenType, a.AccessToken))
						}
					}

					if len(auxTokenList) > 0 {
						return autorest.Prepare(req, autorest.WithHeader("x-ms-authorization-auxiliary", strings.Join(auxTokenList, ", ")))
					}
				}

				return req, nil
			}

			return req, err
		})
	}
}

// BearerAuthorizerCallback is a helper that returns an *autorest.BearerAuthorizerCallback for use in data plane API clients in the Azure SDK
func (c *Authorizer) BearerAuthorizerCallback() *autorest.BearerAuthorizerCallback {
	return autorest.NewBearerAuthorizerCallback(nil, func(_, resource string) (*autorest.BearerAuthorizer, error) {
		token, err := c.Token(context.TODO(), &http.Request{})
		if err != nil {
			return nil, fmt.Errorf("obtaining token: %v", err)
		}

		return autorest.NewBearerAuthorizer(&adalTokenProvider{
			tokenType:  "Bearer",
			tokenValue: token.AccessToken,
		}), nil
	})
}

type adalTokenProvider struct {
	tokenType  string
	tokenValue string
}

func (s *adalTokenProvider) OAuthToken() string {
	return s.tokenValue
}

func (s *adalTokenProvider) Token() adal.Token {
	return adal.Token{
		AccessToken: s.tokenValue,
		Type:        s.tokenType,
	}
}
