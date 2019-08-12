package resource

import (
	providers "github.com/Azure/azure-sdk-for-go/profiles/2017-03-09/resources/mgmt/resources"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2016-09-01/locks"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	GroupsClient      resources.GroupsClient
	DeploymentsClient resources.DeploymentsClient
	LocksClient       locks.ManagementLocksClient
	ProvidersClient   providers.ProvidersClient
}

func BuildClient(o *common.ClientOptions) *Client {
	c := Client{}

	c.LocksClient = locks.NewManagementLocksClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.LocksClient.Client, o.ResourceManagerAuthorizer)

	c.DeploymentsClient = resources.NewDeploymentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.DeploymentsClient.Client, o.ResourceManagerAuthorizer)

	c.GroupsClient = resources.NewGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.GroupsClient.Client, o.ResourceManagerAuthorizer)

	// this has to come from the Profile since this is shared with Stack
	c.ProvidersClient = providers.NewProvidersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ProvidersClient.Client, o.ResourceManagerAuthorizer)

	return &c
}
