package client

import (
	"github.com/Azure/azure-sdk-for-go/services/managedservices/mgmt/2019-06-01/managedservices"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	LighthouseDefinitionsClient *managedservices.RegistrationDefinitionsClient
	LighthouseAssignmentsClient *managedservices.RegistrationAssignmentsClient
}

func NewClient(o *common.ClientOptions) *Client {
	LighthouseDefinitionsClient := managedservices.NewRegistrationDefinitionsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&LighthouseDefinitionsClient.Client, o.ResourceManagerAuthorizer)

	LighthouseAssignmentsClient := managedservices.NewRegistrationAssignmentsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&LighthouseAssignmentsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		LighthouseDefinitionsClient: &LighthouseDefinitionsClient,
		LighthouseAssignmentsClient: &LighthouseAssignmentsClient,
	}
}
