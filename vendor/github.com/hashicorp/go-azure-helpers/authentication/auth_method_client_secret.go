package authentication

import (
	"context"
	"fmt"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/hashicorp/go-multierror"
	authWrapper "github.com/manicminer/hamilton-autorest/auth"
	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/environments"
)

type servicePrincipalClientSecretAuth struct {
	clientId       string
	clientSecret   string
	environment    string
	subscriptionId string
	tenantId       string
	tenantOnly     bool
}

func (a servicePrincipalClientSecretAuth) build(b Builder) (authMethod, error) {
	method := servicePrincipalClientSecretAuth{
		clientId:       b.ClientID,
		clientSecret:   b.ClientSecret,
		environment:    b.Environment,
		subscriptionId: b.SubscriptionID,
		tenantId:       b.TenantID,
		tenantOnly:     b.TenantOnly,
	}
	return method, nil
}

func (a servicePrincipalClientSecretAuth) isApplicable(b Builder) bool {
	return b.SupportsClientSecretAuth && b.ClientSecret != ""
}

func (a servicePrincipalClientSecretAuth) name() string {
	return "Service Principal / Client Secret"
}

func (a servicePrincipalClientSecretAuth) getADALToken(_ context.Context, sender autorest.Sender, oauthConfig *OAuthConfig, endpoint string) (autorest.Authorizer, error) {
	if oauthConfig.OAuth == nil {
		return nil, fmt.Errorf("getting Authorization Token for client secret auth: an OAuth token wasn't configured correctly; please file a bug with more details")
	}

	spt, err := adal.NewServicePrincipalToken(*oauthConfig.OAuth, a.clientId, a.clientSecret, endpoint)
	if err != nil {
		return nil, err
	}
	spt.SetSender(sender)

	return autorest.NewBearerAuthorizer(spt), nil
}

func (a servicePrincipalClientSecretAuth) getMSALToken(ctx context.Context, api environments.Api, _ autorest.Sender, _ *OAuthConfig, _ string) (autorest.Authorizer, error) {
	environment, err := environments.EnvironmentFromString(a.environment)
	if err != nil {
		return nil, fmt.Errorf("environment config error: %v", err)
	}

	conf := auth.ClientCredentialsConfig{
		Environment:  environment,
		TenantID:     a.tenantId,
		ClientID:     a.clientId,
		ClientSecret: a.clientSecret,
		Scopes:       []string{api.DefaultScope()},
		TokenVersion: auth.TokenVersion2,
	}

	return &authWrapper.Authorizer{Authorizer: conf.TokenSource(ctx, auth.ClientCredentialsSecretType)}, nil
}

func (a servicePrincipalClientSecretAuth) populateConfig(c *Config) error {
	c.AuthenticatedAsAServicePrincipal = true
	c.GetAuthenticatedObjectID = buildServicePrincipalObjectIDFunc(c)
	return nil
}

func (a servicePrincipalClientSecretAuth) validate() error {
	var err *multierror.Error

	fmtErrorMessage := "A %s must be configured when authenticating as a Service Principal using a Client Secret."

	if !a.tenantOnly && a.subscriptionId == "" {
		err = multierror.Append(err, fmt.Errorf(fmtErrorMessage, "Subscription ID"))
	}
	if a.clientId == "" {
		err = multierror.Append(err, fmt.Errorf(fmtErrorMessage, "Client ID"))
	}
	if a.clientSecret == "" {
		err = multierror.Append(err, fmt.Errorf(fmtErrorMessage, "Client Secret"))
	}
	if a.tenantId == "" {
		err = multierror.Append(err, fmt.Errorf(fmtErrorMessage, "Tenant ID"))
	}

	return err.ErrorOrNil()
}
