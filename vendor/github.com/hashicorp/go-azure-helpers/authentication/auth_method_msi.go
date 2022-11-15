package authentication

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/hashicorp/go-multierror"
	"github.com/manicminer/hamilton/environments"
)

type managedServiceIdentityAuth struct {
	msiEndpoint string
	clientID    string
}

func (a managedServiceIdentityAuth) build(b Builder) (authMethod, error) {
	msiEndpoint := b.MsiEndpoint
	if msiEndpoint == "" {
		//nolint:SA1019
		ep, err := adal.GetMSIVMEndpoint()
		if err != nil {
			return nil, fmt.Errorf("determining MSI Endpoint: ensure the VM has MSI enabled, or configure the MSI Endpoint. Error: %s", err)
		}
		msiEndpoint = ep
	}

	log.Printf("[DEBUG] Using MSI msiEndpoint %q", msiEndpoint)

	auth := managedServiceIdentityAuth{
		msiEndpoint: msiEndpoint,
		clientID:    b.ClientID,
	}
	return auth, nil
}

func (a managedServiceIdentityAuth) isApplicable(b Builder) bool {
	// Per the Azure SDK: if the Endpoint and Sender are present this is App Service/Function Apps
	// which we intentionally don't support at this time
	isAppService := os.Getenv("MSI_ENDPOINT") != "" && os.Getenv("MSI_SECRET") != ""
	return b.SupportsManagedServiceIdentity && !isAppService
}

func (a managedServiceIdentityAuth) name() string {
	return "Managed Service Identity"
}

func (a managedServiceIdentityAuth) getADALToken(_ context.Context, sender autorest.Sender, oauthConfig *OAuthConfig, endpoint string) (autorest.Authorizer, error) {
	log.Printf("[DEBUG] getADALToken with MSI msiEndpoint %q, ClientID %q for msiEndpoint %q", a.msiEndpoint, a.clientID, endpoint)

	if oauthConfig.OAuth == nil {
		return nil, fmt.Errorf("getting Authorization Token for MSI auth: an OAuth token wasn't configured correctly; please file a bug with more details")
	}

	var spt *adal.ServicePrincipalToken
	var err error
	if a.clientID == "" {
		//nolint:SA1019
		spt, err = adal.NewServicePrincipalTokenFromMSI(a.msiEndpoint, endpoint)
		if err != nil {
			return nil, err
		}
	} else {
		//nolint:SA1019
		spt, err = adal.NewServicePrincipalTokenFromMSIWithUserAssignedID(a.msiEndpoint, endpoint, a.clientID)
		if err != nil {
			return nil, fmt.Errorf("failed to get an oauth token from MSI for user assigned identity from MSI endpoint %q with client ID %q for endpoint %q: %v", a.msiEndpoint, a.clientID, endpoint, err)
		}
	}

	spt.SetSender(sender)
	auth := autorest.NewBearerAuthorizer(spt)
	return auth, nil
}

func (a managedServiceIdentityAuth) getMSALToken(ctx context.Context, _ environments.Api, sender autorest.Sender, oauthConfig *OAuthConfig, endpoint string) (autorest.Authorizer, error) {
	// auth tokens come directly from the MSI endpoint, so we'll pass through to the existing method for continuity
	return a.getADALToken(ctx, sender, oauthConfig, endpoint)
}

func (a managedServiceIdentityAuth) populateConfig(c *Config) error {
	// nothing to populate back
	return nil
}

func (a managedServiceIdentityAuth) validate() error {
	var err *multierror.Error

	if a.msiEndpoint == "" {
		err = multierror.Append(err, fmt.Errorf("An MSI Endpoint must be configured"))
	}

	return err.ErrorOrNil()
}
