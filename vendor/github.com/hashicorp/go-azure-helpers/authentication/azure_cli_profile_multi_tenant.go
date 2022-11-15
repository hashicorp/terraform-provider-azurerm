package authentication

import (
	"strings"

	"github.com/Azure/go-autorest/autorest/azure/cli"
	"github.com/hashicorp/go-version"
)

type azureCLIProfileMultiTenant struct {
	profile cli.Profile

	environment        string
	subscriptionId     string
	tenantId           string
	auxiliaryTenantIDs []string

	azVersion version.Version
}

func (a *azureCLIProfileMultiTenant) populateFields() error {
	// Ensure we know the Subscription ID - since it's needed for everything else
	if a.subscriptionId == "" {
		err := a.populateSubscriptionID()
		if err != nil {
			return err
		}
	}

	// Ensure we know the Tenant ID - since it's needed for everything else
	// Note that order matters that subscription id should be populated first, and then tenant id.
	if a.tenantId == "" {
		err := a.populateTenantID()
		if err != nil {
			return err
		}
	}

	// always pull the environment from the Azure CLI, since the Access Token's associated with it
	return a.populateEnvironment()
}

func (a *azureCLIProfileMultiTenant) verifyAuthenticatedAsAUser() bool {
	for _, subscription := range a.profile.Subscriptions {
		if subscription.User == nil {
			continue
		}

		authenticatedAsAUser := strings.EqualFold(subscription.User.Type, "user")
		if authenticatedAsAUser {
			return true
		}
	}

	return false
}
