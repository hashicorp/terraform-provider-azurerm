package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedservices/2019-06-01/registrationassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedservices/2019-06-01/registrationdefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AssignmentsClient *registrationassignments.RegistrationAssignmentsClient
	DefinitionsClient *registrationdefinitions.RegistrationDefinitionsClient
}

func NewClient(o *common.ClientOptions) *Client {
	assignmentsClient := registrationassignments.NewRegistrationAssignmentsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&assignmentsClient.Client, o.ResourceManagerAuthorizer)

	definitionsClient := registrationdefinitions.NewRegistrationDefinitionsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&definitionsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		DefinitionsClient: &definitionsClient,
		AssignmentsClient: &assignmentsClient,
	}
}
