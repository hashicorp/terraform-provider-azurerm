package policy

import (
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/policy"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	AssignmentsClient    policy.AssignmentsClient
	DefinitionsClient    policy.DefinitionsClient
	SetDefinitionsClient policy.SetDefinitionsClient
}

func BuildClient(o *common.ClientOptions) *Client {
	c := Client{}

	c.AssignmentsClient = policy.NewAssignmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.AssignmentsClient.Client, o.ResourceManagerAuthorizer)

	c.DefinitionsClient = policy.NewDefinitionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.DefinitionsClient.Client, o.ResourceManagerAuthorizer)

	c.SetDefinitionsClient = policy.NewSetDefinitionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.SetDefinitionsClient.Client, o.ResourceManagerAuthorizer)

	return &c
}
