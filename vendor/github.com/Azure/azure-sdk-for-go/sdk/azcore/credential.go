// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"context"
	"time"
)

// AuthenticationPolicyOptions contains various options used to create a credential policy.
type AuthenticationPolicyOptions struct {
	// Options contains the TokenRequestOptions that includes a scopes field which contains
	// the list of OAuth2 authentication scopes used when requesting a token.
	// This field is ignored for other forms of authentication (e.g. shared key).
	Options TokenRequestOptions
}

// Credential represents any credential type.
type Credential interface {
	// AuthenticationPolicy returns a policy that requests the credential and applies it to the HTTP request.
	AuthenticationPolicy(options AuthenticationPolicyOptions) Policy
}

// credentialFunc is a type that implements the Credential interface.
// Use this type when implementing a stateless credential as a first-class function.
type credentialFunc func(options AuthenticationPolicyOptions) Policy

// AuthenticationPolicy implements the Credential interface on credentialFunc.
func (cf credentialFunc) AuthenticationPolicy(options AuthenticationPolicyOptions) Policy {
	return cf(options)
}

// TokenCredential represents a credential capable of providing an OAuth token.
type TokenCredential interface {
	Credential
	// GetToken requests an access token for the specified set of scopes.
	GetToken(ctx context.Context, options TokenRequestOptions) (*AccessToken, error)
}

// AccessToken represents an Azure service bearer access token with expiry information.
type AccessToken struct {
	Token     string
	ExpiresOn time.Time
}

// TokenRequestOptions contain specific parameter that may be used by credentials types when attempting to get a token.
type TokenRequestOptions struct {
	// Scopes contains the list of permission scopes required for the token.
	Scopes []string
}
