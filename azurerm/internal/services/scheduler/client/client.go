package client

import (
	"github.com/Azure/azure-sdk-for-go/services/scheduler/mgmt/2016-03-01/scheduler"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

// TODO: remove in 2.0
type Client struct {
	JobCollectionsClient *scheduler.JobCollectionsClient //nolint: megacheck
	JobsClient           *scheduler.JobsClient           //nolint: megacheck
}

func NewClient(o *common.ClientOptions) *Client {
	jobCollectionsClient := scheduler.NewJobCollectionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId) //nolint: megacheck
	o.ConfigureClient(&jobCollectionsClient.Client, o.ResourceManagerAuthorizer)

	jobsClient := scheduler.NewJobsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId) //nolint: megacheck
	o.ConfigureClient(&jobsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		JobCollectionsClient: &jobCollectionsClient,
		JobsClient:           &jobsClient,
	}
}
