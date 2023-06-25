package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/datamigration/2018-04-19/projectresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datamigration/2018-04-19/serviceresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ServicesClient *serviceresource.ServiceResourceClient
	ProjectsClient *projectresource.ProjectResourceClient
}

func NewClient(o *common.ClientOptions) *Client {
	servicesClient := serviceresource.NewServiceResourceClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&servicesClient.Client, o.ResourceManagerAuthorizer)

	projectsClient := projectresource.NewProjectResourceClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&projectsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ServicesClient: &servicesClient,
		ProjectsClient: &projectsClient,
	}
}
