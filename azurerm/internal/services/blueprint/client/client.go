package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/blueprint/mgmt/2018-11-01-preview/blueprint"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	BlueprintClient  *blueprint.BlueprintsClient
	AssignmentClient *blueprint.AssignmentsClient
	PublishClient    *blueprint.PublishedBlueprintsClient
}

func NewClient(o *common.ClientOptions) *Client {
	blueprintsClient := blueprint.NewBlueprintsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&blueprintsClient.Client, o.ResourceManagerAuthorizer)

	assignmentsClient := blueprint.NewAssignmentsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&assignmentsClient.Client, o.ResourceManagerAuthorizer)

	publishClient := blueprint.NewPublishedBlueprintsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&publishClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		BlueprintClient:  &blueprintsClient,
		AssignmentClient: &assignmentsClient,
		PublishClient:    &publishClient,
	}
}
