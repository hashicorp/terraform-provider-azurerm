package client

import (
	"github.com/Azure/azure-sdk-for-go/services/redhatopenshift/mgmt/2020-04-30/redhatopenshift"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	OpenShiftClustersClient *redhatopenshift.OpenShiftClustersClient
}

func NewClient(o *common.ClientOptions) *Client {
	openshiftClustersClient := redhatopenshift.NewOpenShiftClustersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&openshiftClustersClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		OpenShiftClustersClient: &openshiftClustersClient,
	}
}
