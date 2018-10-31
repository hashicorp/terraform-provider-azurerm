package authentication

import (
	"fmt"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/hashicorp/go-multierror"
)

type servicePrincipalClientSecretAuth struct {
	clientId       string
	clientSecret   string
	subscriptionId string
	tenantId       string
	environment    string
}

func newServicePrincipalClientSecretAuth(b Builder) authMethod {
	return servicePrincipalClientSecretAuth{
		clientId:       b.ClientID,
		clientSecret:   b.ClientSecret,
		environment:    b.Environment,
		subscriptionId: b.SubscriptionID,
		tenantId:       b.TenantID,
	}
}

func (a servicePrincipalClientSecretAuth) getAuthorizationToken(c *Config, oauthConfig *adal.OAuthConfig, endpoint string) (*autorest.BearerAuthorizer, error) {
	spt, err := adal.NewServicePrincipalToken(*oauthConfig, c.ClientID, c.clientSecret, endpoint)
	if err != nil {
		return nil, err
	}

	auth := autorest.NewBearerAuthorizer(spt)
	return auth, nil
}

func (a servicePrincipalClientSecretAuth) validate() error {
	var err *multierror.Error

	fmtErrorMessage := "A %s must be configured when authenticating as a Service Principal using a Client Secret."

	if a.subscriptionId == "" {
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
	if a.environment == "" {
		err = multierror.Append(err, fmt.Errorf(fmtErrorMessage, "Environment"))
	}

	return err.ErrorOrNil()
}
