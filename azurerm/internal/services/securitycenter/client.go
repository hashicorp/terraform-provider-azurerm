package securitycenter

import "github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v1.0/security"

type Client struct {
	ContactsClient  security.ContactsClient
	PricingClient   security.PricingsClient
	WorkspaceClient security.WorkspaceSettingsClient
}
