package client

import (
	"github.com/Azure/azure-sdk-for-go/services/managedservices/mgmt/2019-06-01/managedservices"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	DefinitionsClient *managedservices.RegistrationDefinitionsClient
	AssignmentsClient *managedservices.RegistrationAssignmentsClient
}

func NewClient(o *common.ClientOptions) *Client {
	DefinitionsClient := managedservices.NewRegistrationDefinitionsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&DefinitionsClient.Client, o.ResourceManagerAuthorizer)

	AssignmentsClient := managedservices.NewRegistrationAssignmentsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&AssignmentsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		DefinitionsClient: &DefinitionsClient,
		AssignmentsClient: &AssignmentsClient,
	}
}
