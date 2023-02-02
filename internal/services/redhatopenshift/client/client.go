package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/redhatopenshift/2022-09-04/openshiftclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	OpenShiftClustersClient *openshiftclusters.OpenShiftClustersClient
}

func NewClient(o *common.ClientOptions) *Client {
	openshiftClustersClient := openshiftclusters.NewOpenShiftClustersClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&openshiftClustersClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		OpenShiftClustersClient: &openshiftClustersClient,
	}
}
