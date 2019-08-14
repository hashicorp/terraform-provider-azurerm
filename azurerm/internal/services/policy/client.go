package policy

import (
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/policy"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	AssignmentsClient    *policy.AssignmentsClient
	DefinitionsClient    *policy.DefinitionsClient
	SetDefinitionsClient *policy.SetDefinitionsClient
}

func BuildClient(o *common.ClientOptions) *Client {

	AssignmentsClient := policy.NewAssignmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&AssignmentsClient.Client, o.ResourceManagerAuthorizer)

	DefinitionsClient := policy.NewDefinitionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DefinitionsClient.Client, o.ResourceManagerAuthorizer)

	SetDefinitionsClient := policy.NewSetDefinitionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&SetDefinitionsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AssignmentsClient:    &AssignmentsClient,
		DefinitionsClient:    &DefinitionsClient,
		SetDefinitionsClient: &SetDefinitionsClient,
	}
}
