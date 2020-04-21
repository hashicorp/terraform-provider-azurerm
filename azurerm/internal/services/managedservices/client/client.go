package client

import (
	"github.com/Azure/azure-sdk-for-go/services/managedservices/mgmt/2019-06-01/managedservices"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	RegistrationDefinitionsClient *managedservices.RegistrationDefinitionsClient
	RegistrationAssignmentsClient *managedservices.RegistrationAssignmentsClient
}

func NewClient(o *common.ClientOptions) *Client {
	RegistrationDefinitionsClient := managedservices.NewRegistrationDefinitionsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&RegistrationDefinitionsClient.Client, o.ResourceManagerAuthorizer)

	RegistrationAssignmentsClient := managedservices.NewRegistrationAssignmentsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&RegistrationAssignmentsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		RegistrationDefinitionsClient: &RegistrationDefinitionsClient,
		RegistrationAssignmentsClient: &RegistrationAssignmentsClient,
	}
}
