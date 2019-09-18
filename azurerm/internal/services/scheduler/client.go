package scheduler

import (
	"github.com/Azure/azure-sdk-for-go/services/scheduler/mgmt/2016-03-01/scheduler"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

// TODO: remove in 2.0
type Client struct {
	JobCollectionsClient *scheduler.JobCollectionsClient //nolint: megacheck
	JobsClient           *scheduler.JobsClient           //nolint: megacheck
}

func BuildClient(o *common.ClientOptions) *Client {

	JobCollectionsClient := scheduler.NewJobCollectionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId) //nolint: megacheck
	o.ConfigureClient(&JobCollectionsClient.Client, o.ResourceManagerAuthorizer)

	JobsClient := scheduler.NewJobsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId) //nolint: megacheck
	o.ConfigureClient(&JobsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		JobCollectionsClient: &JobCollectionsClient,
		JobsClient:           &JobsClient,
	}
}
