package client

import (
	"github.com/Azure/azure-sdk-for-go/services/aad/mgmt/2017-04-01/aad"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	DiagnosticSettingsClient *aad.DiagnosticSettingsClient
}

func NewClient(o *common.ClientOptions) *Client {
	diagnosticsSettingsClient := aad.NewDiagnosticSettingsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&diagnosticsSettingsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		DiagnosticSettingsClient: &diagnosticsSettingsClient,
	}
}
