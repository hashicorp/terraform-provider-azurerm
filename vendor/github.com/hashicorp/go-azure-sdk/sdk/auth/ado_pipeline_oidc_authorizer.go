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

const (
	adoPipelineOIDCAPIVersion = "7.1"
)

type ADOPipelineOIDCAuthorizerOptions struct {
	// Api describes the Azure API being used
	Api environments.Api

	// ClientId is the client ID used when authenticating
	ClientId string

	// ServiceConnectionId is the ADO service connection ID used when authenticating
	ServiceConnectionId string

	// Environment is the Azure environment/cloud being targeted
	Environment environments.Environment

	// TenantId is the tenant to authenticate against
	TenantId string

	// AuxiliaryTenantIds lists additional tenants to authenticate against, currently only
	// used for Resource Manager when auxiliary tenants are needed.
	// e.g. https://learn.microsoft.com/en-us/azure/azure-resource-manager/management/authenticate-multi-tenant
	AuxiliaryTenantIds []string

	// IdTokenRequestUrl is the URL for the OIDC provider from which to request an ID token.
	// Usually exposed via the SYSTEM_OIDCREQUESTURI environment variable when running in ADO Pipelines
	IdTokenRequestUrl string

	// IdTokenRequestToken is the bearer token for the request to the OIDC provider.
	// Usually exposed via the SYSTEM_ACCESSTOKEN environment variable when running in ADO Pipelines
	IdTokenRequestToken string
}

// NewADOPipelineOIDCAuthorizer returns an authorizer which acquires a client assertion from a ADO endpoint, then uses client assertion authentication to obtain an access token.
func NewADOPipelineOIDCAuthorizer(ctx context.Context, options ADOPipelineOIDCAuthorizerOptions) (Authorizer, error) {
	scope, err := environments.Scope(options.Api)
	if err != nil {
		return nil, fmt.Errorf("determining scope for %q: %+v", options.Api.Name(), err)
	}

	conf := adoPipelineOIDCConfig{
		Environment:         options.Environment,
		TenantID:            options.TenantId,
		AuxiliaryTenantIDs:  options.AuxiliaryTenantIds,
		ClientID:            options.ClientId,
		ServiceConnectionID: options.ServiceConnectionId,
		IDTokenRequestURL:   options.IdTokenRequestUrl,
		IDTokenRequestToken: options.IdTokenRequestToken,
		Scopes: []string{
			*scope,
		},
	}

	return conf.TokenSource(ctx)
}

var _ Authorizer = &ADOPipelineOIDCAuthorizer{}

type ADOPipelineOIDCAuthorizer struct {
	conf *adoPipelineOIDCConfig
}

func (a *ADOPipelineOIDCAuthorizer) adoPipelineAssertion(ctx context.Context, _ *http.Request) (*string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.conf.IDTokenRequestURL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("adoPipelineAssertion: failed to build request: %+v", err)
	}

	query, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		return nil, fmt.Errorf("adoPipelineAssertion: cannot parse URL query")
	}
	if query.Get("api-version") == "" {
		query.Add("api-version", adoPipelineOIDCAPIVersion)
	}
	if query.Get("serviceConnectionId") == "" {
		query.Add("serviceConnectionId", a.conf.ServiceConnectionID)
	}
	if query.Get("audience") == "" {
		query.Add("audience", "api://AzureADTokenExchange")
	}
	req.URL.RawQuery = query.Encode()

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.conf.IDTokenRequestToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("adoPipelineAssertion: cannot request token: %v", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("adoPipelineAssertion: cannot parse response: %v", err)
	}

	if c := resp.StatusCode; c < 200 || c > 299 {
		return nil, fmt.Errorf("adoPipelineAssertion: received HTTP status %d with response: %s", resp.StatusCode, body)
	}

	var tokenRes struct {
		Value *string `json:"oidcToken"`
	}
	if err := json.Unmarshal(body, &tokenRes); err != nil {
		return nil, fmt.Errorf("adoPipelineAssertion: cannot unmarshal response: %v", err)
	}

	return tokenRes.Value, nil
}

func (a *ADOPipelineOIDCAuthorizer) tokenSource(ctx context.Context, req *http.Request) (Authorizer, error) {
	assertion, err := a.adoPipelineAssertion(ctx, req)
	if err != nil {
		return nil, err
	}
	if assertion == nil {
		return nil, fmt.Errorf("ADOPipelineOIDCAuthorizer: nil JWT assertion received from ADOPipeline")
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
		return nil, fmt.Errorf("ADOPipelineOIDCAuthorizer: building Authorizer: %+v", err)
	}
	return source, nil
}

func (a *ADOPipelineOIDCAuthorizer) Token(ctx context.Context, req *http.Request) (*oauth2.Token, error) {
	source, err := a.tokenSource(ctx, req)
	if err != nil {
		return nil, err
	}
	return source.Token(ctx, req)
}

func (a *ADOPipelineOIDCAuthorizer) AuxiliaryTokens(ctx context.Context, req *http.Request) ([]*oauth2.Token, error) {
	source, err := a.tokenSource(ctx, req)
	if err != nil {
		return nil, err
	}
	return source.AuxiliaryTokens(ctx, req)
}

type adoPipelineOIDCConfig struct {
	// Environment is the national cloud environment to use
	Environment environments.Environment

	// TenantID is the required tenant ID for the primary token
	TenantID string

	// AuxiliaryTenantIDs is an optional list of tenant IDs for which to obtain additional tokens
	AuxiliaryTenantIDs []string

	// ClientID is the application's ID.
	ClientID string

	// ServiceConnectionID is the ADO service connection ID used when authenticating
	ServiceConnectionID string

	// IDTokenRequestURL is the URL for ADO Pipeline's OIDC provider.
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

func (c *adoPipelineOIDCConfig) TokenSource(ctx context.Context) (Authorizer, error) {
	return NewCachedAuthorizer(&ADOPipelineOIDCAuthorizer{
		conf: c,
	})
}
