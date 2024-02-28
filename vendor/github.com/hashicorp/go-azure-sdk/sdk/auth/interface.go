// Copyright (c) HashiCorp Inc. All rights reserved.
// Licensed under the MPL-2.0 License. See NOTICE.txt in the project root for license information.

package auth

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
)

// Authorizer is anything that can return an access token for authorizing API connections
type Authorizer interface {
	// Token obtains a new access token for the configured tenant
	Token(ctx context.Context, request *http.Request) (*oauth2.Token, error)

	// AuxiliaryTokens obtains new access tokens for the configured auxiliary tenants
	AuxiliaryTokens(ctx context.Context, request *http.Request) ([]*oauth2.Token, error)
}

// CachingAuthorizer implements Authorizer whilst caching access tokens and offering a way to intentionally invalidate them
type CachingAuthorizer interface {
	Authorizer

	// InvalidateCachedTokens invalidates any cached access tokens, so that new tokens are automatically
	// retrieved from the authorization service on the next call to Token or AuxiliaryTokens.
	InvalidateCachedTokens() error
}

// HTTPClient is an HTTP client used for sending authentication requests and obtaining tokens
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
