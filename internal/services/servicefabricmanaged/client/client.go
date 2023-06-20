package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicefabricmanagedcluster/2021-05-01/managedcluster"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicefabricmanagedcluster/2021-05-01/nodetype"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ManagedClusterClient *managedcluster.ManagedClusterClient
	NodeTypeClient       *nodetype.NodeTypeClient
}

func NewClient(o *common.ClientOptions) *Client {
	managedCluster := managedcluster.NewManagedClusterClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&managedCluster.Client, o.ResourceManagerAuthorizer)

	nodeType := nodetype.NewNodeTypeClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&nodeType.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ManagedClusterClient: &managedCluster,
		NodeTypeClient:       &nodeType,
	}
}
