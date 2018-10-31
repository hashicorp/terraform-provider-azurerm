package authentication

import (
	"fmt"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
)

// GetAuthorizationToken returns an authentication token for the current authentication method
func GetAuthorizationToken(c *Config, oauthConfig *adal.OAuthConfig, endpoint string) (*autorest.BearerAuthorizer, error) {
	if c.usingClientSecret {
		spt, err := adal.NewServicePrincipalToken(*oauthConfig, c.ClientID, c.clientSecret, endpoint)
		if err != nil {
			return nil, err
		}

		auth := autorest.NewBearerAuthorizer(spt)
		return auth, nil
	}

	if c.usingManagedServiceIdentity {
		spt, err := adal.NewServicePrincipalTokenFromMSI(c.msiEndpoint, endpoint)
		if err != nil {
			return nil, err
		}
		auth := autorest.NewBearerAuthorizer(spt)
		return auth, nil
	}

	if c.usingCloudShell {
		// load the refreshed tokens from the CloudShell Azure CLI credentials
		err := c.LoadTokensFromAzureCLI()
		if err != nil {
			return nil, fmt.Errorf("Error loading the refreshed CloudShell tokens: %+v", err)
		}
	}

	spt, err := adal.NewServicePrincipalTokenFromManualToken(*oauthConfig, c.ClientID, endpoint, *c.accessToken)
	if err != nil {
		return nil, err
	}

	err = spt.Refresh()

	if err != nil {
		return nil, fmt.Errorf("Error refreshing Service Principal Token: %+v", err)
	}

	auth := autorest.NewBearerAuthorizer(spt)
	return auth, nil
}
