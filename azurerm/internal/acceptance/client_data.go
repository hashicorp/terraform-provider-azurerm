package acceptance

import "os"

type ClientData struct {
	// ClientID is the UUID of the Service Principal being used to connect to Azure
	ClientID string

	// ClientSecret is the Client Secret being used to connect to Azure
	ClientSecret string

	// SubscriptionID is the UUID of the Azure Subscription where tests are being run
	SubscriptionID string

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
		ClientID:           os.Getenv("ARM_CLIENT_ID"),
		ClientSecret:       os.Getenv("ARM_CLIENT_SECRET"),
		IsServicePrincipal: true,
		SubscriptionID:     os.Getenv("ARM_SUBSCRIPTION_ID"),
		TenantID:           os.Getenv("ARM_TENANT_ID"),
	}
}
