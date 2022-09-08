package authentication

import (
	"context"
	"fmt"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-multierror"
	authWrapper "github.com/manicminer/hamilton-autorest/auth"
	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/environments"
)

type oidcAuth struct {
	auxiliaryTenantIds  []string
	clientId            string
	environment         string
	idToken             string
	idTokenRequestToken string
	idTokenRequestUrl   string
	tenantId            string
}

func (a oidcAuth) build(b Builder) (authMethod, error) {
	method := oidcAuth{
		auxiliaryTenantIds:  b.AuxiliaryTenantIDs,
		clientId:            b.ClientID,
		environment:         b.Environment,
		idToken:             b.IDToken,
		idTokenRequestUrl:   b.IDTokenRequestURL,
		idTokenRequestToken: b.IDTokenRequestToken,
		tenantId:            b.TenantID,
	}
	return method, nil
}

func (a oidcAuth) isApplicable(b Builder) bool {
	return b.SupportsOIDCAuth && b.UseMicrosoftGraph && (b.IDToken != "" || (b.IDTokenRequestURL != "" && b.IDTokenRequestToken != ""))
}

func (a oidcAuth) name() string {
	return "OIDC"
}

func (a oidcAuth) getADALToken(_ context.Context, _ autorest.Sender, _ *OAuthConfig, _ string) (autorest.Authorizer, error) {
	return nil, fmt.Errorf("ADAL tokens are not supported for OIDC authentication")
}

func (a oidcAuth) getMSALToken(ctx context.Context, api environments.Api, _ autorest.Sender, _ *OAuthConfig, _ string) (autorest.Authorizer, error) {
	environment, err := environments.EnvironmentFromString(a.environment)
	if err != nil {
		return nil, fmt.Errorf("environment config error: %v", err)
	}

	if a.idToken == "" {
		conf := auth.GitHubOIDCConfig{
			Environment:         environment,
			TenantID:            a.tenantId,
			AuxiliaryTenantIDs:  a.auxiliaryTenantIds,
			ClientID:            a.clientId,
			IDTokenRequestURL:   a.idTokenRequestUrl,
			IDTokenRequestToken: a.idTokenRequestToken,
			Scopes:              []string{api.DefaultScope()},
		}
		return &authWrapper.Authorizer{Authorizer: conf.TokenSource(ctx)}, nil
	}

	conf := auth.ClientCredentialsConfig{
		Environment:        environment,
		TenantID:           a.tenantId,
		AuxiliaryTenantIDs: a.auxiliaryTenantIds,
		ClientID:           a.clientId,
		FederatedAssertion: a.idToken,
		Scopes:             []string{api.DefaultScope()},
		TokenVersion:       auth.TokenVersion2,
	}

	return &authWrapper.Authorizer{Authorizer: conf.TokenSource(ctx, auth.ClientCredentialsAssertionType)}, nil
}

func (a oidcAuth) populateConfig(c *Config) error {
	c.AuthenticatedAsAServicePrincipal = true
	c.AuthenticatedViaOIDC = true
	c.GetAuthenticatedObjectID = buildServicePrincipalObjectIDFunc(c)
	return nil
}

func (a oidcAuth) validate() error {
	var err *multierror.Error

	fmtErrorMessage := "a %s must be configured when authenticating with OIDC"

	if a.tenantId == "" {
		err = multierror.Append(err, fmt.Errorf(fmtErrorMessage, "Tenant ID"))
	}

	if a.clientId == "" {
		err = multierror.Append(err, fmt.Errorf(fmtErrorMessage, "Client ID"))
	}

	if a.idTokenRequestUrl == "" && a.idToken == "" {
		err = multierror.Append(err, fmt.Errorf(fmtErrorMessage, "ID Token or ID Token Request URL"))
	}

	if a.idTokenRequestToken == "" && a.idToken == "" {
		err = multierror.Append(err, fmt.Errorf(fmtErrorMessage, "ID Token or ID Token Request Token"))
	}

	return err.ErrorOrNil()
}
