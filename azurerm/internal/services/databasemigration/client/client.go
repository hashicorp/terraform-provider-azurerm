package client

import (
	"github.com/Azure/azure-sdk-for-go/services/datamigration/mgmt/2018-04-19/datamigration"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ServicesClient *datamigration.ServicesClient
	ProjectsClient *datamigration.ProjectsClient
}

func NewClient(o *common.ClientOptions) *Client {
	servicesClient := datamigration.NewServicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&servicesClient.Client, o.ResourceManagerAuthorizer)

	projectsClient := datamigration.NewProjectsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&projectsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ServicesClient: &servicesClient,
		ProjectsClient: &projectsClient,
	}
}
