package acceptance

import (
	"os"
)

type ClientAuthDetails struct {
	// ClientID is the UUID of the Service Principal being used to connect to Azure
	ClientID string

	// ClientSecret is the Client Secret being used to connect to Azure
	ClientSecret string
}

type ClientData struct {
	// Default is the Default Credentials being used to connect to Azure
	Default ClientAuthDetails

	// Alternate is an alternate set of Credentials being used to connect to Azure
	Alternate ClientAuthDetails

	// SubscriptionID is the UUID of the Azure Subscription where tests are being run
	SubscriptionID string

	// SubscriptionIDAlt is the UUID of the Alternate Azure Subscription where tests are being run
	SubscriptionIDAlt string

	// TenantID is the UUID of the Azure Tenant where tests are being run
	TenantID string

	// Are we connected as a Service Principal
	// this exists for allowing folks to run the test suite via other
	// credentials in the future; but requires code changes first
	IsServicePrincipal bool
}

// Client returns a struct containing information about the Client being
// used to connect to Azure
func (td TestData) Client() ClientData {
	return ClientData{
		Default: ClientAuthDetails{
			ClientID:     os.Getenv("ARM_CLIENT_ID"),
			ClientSecret: os.Getenv("ARM_CLIENT_SECRET"),
		},
		Alternate: ClientAuthDetails{
			ClientID:     os.Getenv("ARM_CLIENT_ID_ALT"),
			ClientSecret: os.Getenv("ARM_CLIENT_SECRET_ALT"),
		},
		IsServicePrincipal: true,
		SubscriptionID:     os.Getenv("ARM_SUBSCRIPTION_ID"),
		SubscriptionIDAlt:  os.Getenv("ARM_SUBSCRIPTION_ID_ALT"),
		TenantID:           os.Getenv("ARM_TENANT_ID"),
	}
}
