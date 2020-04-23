package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/blueprint/mgmt/2018-11-01-preview/blueprint"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	AssignmentsClient         *blueprint.AssignmentsClient
	BlueprintsClient          *blueprint.BlueprintsClient
	PublishedBlueprintsClient *blueprint.PublishedBlueprintsClient
}

func NewClient(o *common.ClientOptions) *Client {
	assignmentsClient := blueprint.NewAssignmentsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&assignmentsClient.Client, o.ResourceManagerAuthorizer)

	blueprintsClient := blueprint.NewBlueprintsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&blueprintsClient.Client, o.ResourceManagerAuthorizer)

	publishedBlueprintsClient := blueprint.NewPublishedBlueprintsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&publishedBlueprintsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AssignmentsClient:         &assignmentsClient,
		BlueprintsClient:          &blueprintsClient,
		PublishedBlueprintsClient: &publishedBlueprintsClient,
	}
}
