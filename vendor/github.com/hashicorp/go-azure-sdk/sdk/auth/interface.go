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
	Token(ctx context.Context, request *http.Request) (*oauth2.Token, error)

	AuxiliaryTokens(ctx context.Context, request *http.Request) ([]*oauth2.Token, error)
}

// HTTPClient is an HTTP client used for sending authentication requests and obtaining tokens
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
