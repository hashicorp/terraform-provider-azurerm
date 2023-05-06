// Copyright (c) HashiCorp Inc. All rights reserved.
// Licensed under the MPL-2.0 License. See NOTICE.txt in the project root for license information.

package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"golang.org/x/oauth2"
)

type GitHubOIDCAuthorizerOptions struct {
	// Api describes the Azure API being used
	Api environments.Api

	// ClientId is the client ID used when authenticating
	ClientId string

	// Environment is the Azure environment/cloud being targeted
	Environment environments.Environment

	// TenantId is the tenant to authenticate against
	TenantId string

	// AuxiliaryTenantIds lists additional tenants to authenticate against, currently only
	// used for Resource Manager when auxiliary tenants are needed.
	// e.g. https://learn.microsoft.com/en-us/azure/azure-resource-manager/management/authenticate-multi-tenant
	AuxiliaryTenantIds []string

	// IdTokenRequestUrl is the URL for the OIDC provider from which to request an ID token.
	// Usually exposed via the ACTIONS_ID_TOKEN_REQUEST_URL environment variable when running in GitHub Actions
	IdTokenRequestUrl string

	// IdTokenRequestToken is the bearer token for the request to the OIDC provider.
	// Usually exposed via the ACTIONS_ID_TOKEN_REQUEST_TOKEN environment variable when running in GitHub Actions
	IdTokenRequestToken string
}

// NewGitHubOIDCAuthorizer returns an authorizer which acquires a client assertion from a GitHub endpoint, then uses client assertion authentication to obtain an access token.
func NewGitHubOIDCAuthorizer(ctx context.Context, options GitHubOIDCAuthorizerOptions) (Authorizer, error) {
	scope, err := environments.Scope(options.Api)
	if err != nil {
		return nil, fmt.Errorf("determining scope for %q: %+v", options.Api.Name(), err)
	}

	conf := gitHubOIDCConfig{
		Environment:         options.Environment,
		TenantID:            options.TenantId,
		AuxiliaryTenantIDs:  options.AuxiliaryTenantIds,
		ClientID:            options.ClientId,
		IDTokenRequestURL:   options.IdTokenRequestUrl,
		IDTokenRequestToken: options.IdTokenRequestToken,
		Scopes: []string{
			*scope,
		},
	}

	return conf.TokenSource(ctx)
}

var _ Authorizer = &GitHubOIDCAuthorizer{}

type GitHubOIDCAuthorizer struct {
	conf *gitHubOIDCConfig
}

func (a *GitHubOIDCAuthorizer) githubAssertion(ctx context.Context, _ *http.Request) (*string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.conf.IDTokenRequestURL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("githubAssertion: failed to build request: %+v", err)
	}

	query, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		return nil, fmt.Errorf("githubAssertion: cannot parse URL query")
	}

	if query.Get("audience") == "" {
		query.Set("audience", "api://AzureADTokenExchange")
		req.URL.RawQuery = query.Encode()
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.conf.IDTokenRequestToken))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("githubAssertion: cannot request token: %v", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("githubAssertion: cannot parse response: %v", err)
	}

	if c := resp.StatusCode; c < 200 || c > 299 {
		return nil, fmt.Errorf("githubAssertion: received HTTP status %d with response: %s", resp.StatusCode, body)
	}

	var tokenRes struct {
		Count *int    `json:"count"`
		Value *string `json:"value"`
	}
	if err := json.Unmarshal(body, &tokenRes); err != nil {
		return nil, fmt.Errorf("githubAssertion: cannot unmarshal response: %v", err)
	}

	return tokenRes.Value, nil
}

func (a *GitHubOIDCAuthorizer) tokenSource(ctx context.Context, req *http.Request) (Authorizer, error) {
	assertion, err := a.githubAssertion(ctx, req)
	if err != nil {
		return nil, err
	}
	if assertion == nil {
		return nil, fmt.Errorf("GitHubOIDCAuthorizer: nil JWT assertion received from GitHub")
	}

	conf := clientCredentialsConfig{
		Environment:        a.conf.Environment,
		TenantID:           a.conf.TenantID,
		AuxiliaryTenantIDs: a.conf.AuxiliaryTenantIDs,
		ClientID:           a.conf.ClientID,
		FederatedAssertion: *assertion,
		Scopes:             a.conf.Scopes,
		TokenURL:           a.conf.TokenURL,
		Audience:           a.conf.Audience,
	}

	source, err := conf.TokenSource(ctx, clientCredentialsAssertionType)
	if err != nil {
		return nil, fmt.Errorf("GitHubOIDCAuthorizer: building Authorizer: %+v", err)
	}
	return source, nil
}

func (a *GitHubOIDCAuthorizer) Token(ctx context.Context, req *http.Request) (*oauth2.Token, error) {
	source, err := a.tokenSource(ctx, req)
	if err != nil {
		return nil, err
	}
	return source.Token(ctx, req)
}

func (a *GitHubOIDCAuthorizer) AuxiliaryTokens(ctx context.Context, req *http.Request) ([]*oauth2.Token, error) {
	source, err := a.tokenSource(ctx, req)
	if err != nil {
		return nil, err
	}
	return source.AuxiliaryTokens(ctx, req)
}

type gitHubOIDCConfig struct {
	// Environment is the national cloud environment to use
	Environment environments.Environment

	// TenantID is the required tenant ID for the primary token
	TenantID string

	// AuxiliaryTenantIDs is an optional list of tenant IDs for which to obtain additional tokens
	AuxiliaryTenantIDs []string

	// ClientID is the application's ID.
	ClientID string

	// IDTokenRequestURL is the URL for GitHub's OIDC provider.
	IDTokenRequestURL string

	// IDTokenRequestToken is the bearer token for the request to the OIDC provider.
	IDTokenRequestToken string

	// Scopes specifies a list of requested permission scopes (used for v2 tokens)
	Scopes []string

	// TokenURL is the clientCredentialsToken endpoint, which overrides the default endpoint constructed from a tenant ID
	TokenURL string

	// Audience optionally specifies the intended audience of the
	// request.  If empty, the value of TokenURL is used as the
	// intended audience.
	Audience string
}

func (c *gitHubOIDCConfig) TokenSource(ctx context.Context) (Authorizer, error) {
	return NewCachedAuthorizer(&GitHubOIDCAuthorizer{
		conf: c,
	})
}
