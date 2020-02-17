package client

import (
	"github.com/Azure/azure-sdk-for-go/services/databox/mgmt/2019-09-01/databox"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	JobClient *databox.JobsClient
}

func NewClient(o *common.ClientOptions) *Client {
	jobClient := databox.NewJobsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&jobClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		JobClient: &jobClient,
	}
}
