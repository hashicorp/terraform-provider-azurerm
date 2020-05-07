package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/operationalinsights/mgmt/2015-11-01-preview/operationalinsights"
	"github.com/Azure/azure-sdk-for-go/services/preview/operationsmanagement/mgmt/2015-11-01-preview/operationsmanagement"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	DataSourcesClient    *operationalinsights.DataSourcesClient
	LinkedServicesClient *operationalinsights.LinkedServicesClient
	SolutionsClient      *operationsmanagement.SolutionsClient
	WorkspacesClient     *operationalinsights.WorkspacesClient
}

func NewClient(o *common.ClientOptions) *Client {
	DataSourcesClient := operationalinsights.NewDataSourcesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DataSourcesClient.Client, o.ResourceManagerAuthorizer)

	WorkspacesClient := operationalinsights.NewWorkspacesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&WorkspacesClient.Client, o.ResourceManagerAuthorizer)

	SolutionsClient := operationsmanagement.NewSolutionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, "Microsoft.OperationsManagement", "solutions", "testing")
	o.ConfigureClient(&SolutionsClient.Client, o.ResourceManagerAuthorizer)

	LinkedServicesClient := operationalinsights.NewLinkedServicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&LinkedServicesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		DataSourcesClient:    &DataSourcesClient,
		LinkedServicesClient: &LinkedServicesClient,
		SolutionsClient:      &SolutionsClient,
		WorkspacesClient:     &WorkspacesClient,
	}
}
