package client

import (
	"github.com/Azure/azure-sdk-for-go/services/storageimportexport/mgmt/2016-11-01/storageimportexport"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	JobClient *storageimportexport.JobsClient
}

func NewClient(o *common.ClientOptions) *Client {
	jobClient := storageimportexport.NewJobsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, "")
	o.ConfigureClient(&jobClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		JobClient: &jobClient,
	}
}
