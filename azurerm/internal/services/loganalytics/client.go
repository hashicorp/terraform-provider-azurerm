package loganalytics

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/operationalinsights/mgmt/2015-11-01-preview/operationalinsights"
	"github.com/Azure/azure-sdk-for-go/services/preview/operationsmanagement/mgmt/2015-11-01-preview/operationsmanagement"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	LinkedServicesClient operationalinsights.LinkedServicesClient
	SolutionsClient      operationsmanagement.SolutionsClient
	WorkspacesClient     operationalinsights.WorkspacesClient
}

func BuildClient(o *common.ClientOptions) *Client {
	c := Client{}

	c.WorkspacesClient = operationalinsights.NewWorkspacesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.WorkspacesClient.Client, o.ResourceManagerAuthorizer)

	c.SolutionsClient = operationsmanagement.NewSolutionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, "Microsoft.OperationsManagement", "solutions", "testing")
	o.ConfigureClient(&c.SolutionsClient.Client, o.ResourceManagerAuthorizer)

	c.LinkedServicesClient = operationalinsights.NewLinkedServicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.LinkedServicesClient.Client, o.ResourceManagerAuthorizer)

	return &c
}
