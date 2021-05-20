package authentication

import (
	"context"
	"fmt"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure/cli"
	"github.com/hashicorp/go-multierror"
)

type azureCliTokenMultiTenantAuth struct {
	profile                      *azureCLIProfileMultiTenant
	servicePrincipalAuthDocsLink string
}

func (a azureCliTokenMultiTenantAuth) build(b Builder) (authMethod, error) {
	auth := azureCliTokenMultiTenantAuth{
		profile: &azureCLIProfileMultiTenant{
			clientId:           b.ClientID,
			environment:        b.Environment,
			subscriptionId:     b.SubscriptionID,
			tenantId:           b.TenantID,
			auxiliaryTenantIDs: b.AuxiliaryTenantIDs,
		},
		servicePrincipalAuthDocsLink: b.ClientSecretDocsLink,
	}
	profilePath, err := cli.ProfilePath()
	if err != nil {
		return nil, fmt.Errorf("Error loading the Profile Path from the Azure CLI: %+v", err)
	}

	profile, err := cli.LoadProfile(profilePath)
	if err != nil {
		return nil, fmt.Errorf("Azure CLI Authorization Profile was not found. Please ensure the Azure CLI is installed and then log-in with `az login`.")
	}

	auth.profile.profile = profile

	// Authenticating as a Service Principal doesn't return all of the information we need for authentication purposes
	// as such Service Principal authentication is supported using the specific auth method
	if authenticatedAsAUser := auth.profile.verifyAuthenticatedAsAUser(); !authenticatedAsAUser {
		return nil, fmt.Errorf(`Authenticating using the Azure CLI is only supported as a User (not a Service Principal).

To authenticate to Azure using a Service Principal, you can use the separate 'Authenticate using a Service Principal'
auth method - instructions for which can be found here: %s

Alternatively you can authenticate using the Azure CLI by using a User Account.`, auth.servicePrincipalAuthDocsLink)
	}

	err = auth.profile.populateFields()
	if err != nil {
		return nil, fmt.Errorf("Error retrieving the Profile from the Azure CLI: %s Please re-authenticate using `az login`.", err)
	}

	err = auth.profile.populateClientId()
	if err != nil {
		return nil, fmt.Errorf("Error populating Client ID from the Azure CLI: %+v", err)
	}

	return auth, nil
}

func (a azureCliTokenMultiTenantAuth) isApplicable(b Builder) bool {
	return b.SupportsAzureCliToken && b.SupportsAuxiliaryTenants && (len(b.AuxiliaryTenantIDs) > 0)
}

func (a azureCliTokenMultiTenantAuth) getAuthorizationToken(sender autorest.Sender, oauth *OAuthConfig, endpoint string) (autorest.Authorizer, error) {
	if oauth.MultiTenantOauth == nil {
		return nil, fmt.Errorf("Error getting Authorization Token for cli auth: an MultiTenantOauth token wasn't configured correctly; please file a bug with more details")
	}

	m := adal.MultiTenantServicePrincipalToken{
		AuxiliaryTokens: make([]*adal.ServicePrincipalToken, len(a.profile.auxiliaryTenantIDs)),
	}

	// the Azure CLI appears to cache these, so to maintain compatibility with the interface this method is intentionally not on the pointer
	primaryToken, err := obtainAuthorizationTokenByTenant(endpoint, a.profile.tenantId)
	if err != nil {
		return nil, fmt.Errorf("Error obtaining Authorization Token from the Azure CLI: %s", err)
	}

	adalToken, err := primaryToken.ToADALToken()
	if err != nil {
		return nil, fmt.Errorf("Error converting Authorization Token to an ADAL Token: %s", err)
	}

	spt, err := adal.NewServicePrincipalTokenFromManualToken(*oauth.OAuth, a.profile.clientId, endpoint, adalToken)
	if err != nil {
		return nil, err
	}

	var refreshFunc adal.TokenRefresh = func(ctx context.Context, resource string) (*adal.Token, error) {
		token, err := obtainAuthorizationToken(resource, a.profile.subscriptionId, "")
		if err != nil {
			return nil, err
		}

		adalToken, err := token.ToADALToken()
		if err != nil {
			return nil, err
		}

		return &adalToken, nil
	}

	spt.SetCustomRefreshFunc(refreshFunc)

	m.PrimaryToken = spt
	for t := range a.profile.auxiliaryTenantIDs {
		token, err := obtainAuthorizationTokenByTenant(endpoint, a.profile.auxiliaryTenantIDs[t])
		if err != nil {
			return nil, fmt.Errorf("Error obtaining Authorization Token from the Azure CLI: %s", err)
		}

		adalToken, err := token.ToADALToken()
		if err != nil {
			return nil, fmt.Errorf("Error converting Authorization Token to an ADAL Token: %s", err)
		}

		aux, err := adal.NewServicePrincipalTokenFromManualToken(*oauth.OAuth, a.profile.clientId, endpoint, adalToken)
		if err != nil {
			return nil, err
		}

		aux.SetCustomRefreshFunc(refreshFunc)

		m.AuxiliaryTokens[t] = aux
	}

	auth := autorest.NewMultiTenantServicePrincipalTokenAuthorizer(&m)
	return auth, nil
}

func (a azureCliTokenMultiTenantAuth) name() string {
	return "Obtaining a Multi-tenant token from the Azure CLI"
}

func (a azureCliTokenMultiTenantAuth) populateConfig(c *Config) error {
	c.ClientID = a.profile.clientId
	c.TenantID = a.profile.tenantId
	c.Environment = a.profile.environment
	c.SubscriptionID = a.profile.subscriptionId

	c.GetAuthenticatedObjectID = func(ctx context.Context) (string, error) {
		objectId, err := obtainAuthenticatedObjectID()
		if err != nil {
			return "", err
		}

		return objectId, nil
	}

	return nil
}

func (a azureCliTokenMultiTenantAuth) validate() error {
	var err *multierror.Error

	errorMessageFmt := "A %s was not found in your Azure CLI Credentials.\n\nPlease login to the Azure CLI again via `az login`"

	if a.profile == nil {
		return fmt.Errorf("Azure CLI Profile is nil - this is an internal error and should be reported.")
	}

	if a.profile.clientId == "" {
		err = multierror.Append(err, fmt.Errorf(errorMessageFmt, "Client ID"))
	}

	if a.profile.subscriptionId == "" {
		err = multierror.Append(err, fmt.Errorf(errorMessageFmt, "Subscription ID"))
	}

	if a.profile.tenantId == "" {
		err = multierror.Append(err, fmt.Errorf(errorMessageFmt, "Tenant ID"))
	}

	if len(a.profile.auxiliaryTenantIDs) == 0 {
		err = multierror.Append(err, fmt.Errorf("Aux Tenant IDs missing from Multi Tenant configuration"))
	}

	return err.ErrorOrNil()
}

func obtainAuthorizationTokenByTenant(endpoint string, tenantId string) (*cli.Token, error) {
	var token cli.Token
	err := jsonUnmarshalAzCmd(&token, "account", "get-access-token", "--resource", endpoint, "--tenant", tenantId, "--only-show-errors", "-o=json")
	if err != nil {
		return nil, fmt.Errorf("Error parsing json result from the Azure CLI: %v", err)
	}

	return &token, nil
}
