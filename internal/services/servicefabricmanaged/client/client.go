package client

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicefabricmanaged/sdk/2021-05-01/managedcluster"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicefabricmanaged/sdk/2021-05-01/nodetype"
)

type Client struct {
	ManagedClusterClient *managedcluster.ManagedClusterClient
	NodeTypeClient       *nodetype.NodeTypeClient
	tokenFunc            func(endpoint string) (autorest.Authorizer, error)
	configureClientFunc  func(c *autorest.Client, authorizer autorest.Authorizer)
}

func NewClient(o *common.ClientOptions) *Client {
	managedCluster := managedcluster.NewManagedClusterClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&managedCluster.Client, o.ResourceManagerAuthorizer)

	nodeType := nodetype.NewNodeTypeClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&nodeType.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ManagedClusterClient: &managedCluster,
		NodeTypeClient:       &nodeType,
		tokenFunc:            o.TokenFunc,
		configureClientFunc:  o.ConfigureClient,
	}
}
