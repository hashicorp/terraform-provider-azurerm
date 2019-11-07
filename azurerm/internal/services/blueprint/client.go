package blueprint

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/blueprint/mgmt/2018-11-01-preview/blueprint"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	BlueprintsClient          *blueprint.BlueprintsClient
	AssignmentsClient         *blueprint.AssignmentsClient
	ArtifactsClient           *blueprint.ArtifactsClient
	PublishedArtifactsClient  *blueprint.PublishedArtifactsClient
	PublishedBlueprintsClient *blueprint.PublishedBlueprintsClient
}

func BuildClient(o *common.ClientOptions) *Client {
	BlueprintsClient := blueprint.NewBlueprintsClient()
	o.ConfigureClient(&BlueprintsClient.Client, o.ResourceManagerAuthorizer)

	AssignmentsClient := blueprint.NewAssignmentsClient()
	o.ConfigureClient(&AssignmentsClient.Client, o.ResourceManagerAuthorizer)

	ArtifactsClient := blueprint.NewArtifactsClient()
	o.ConfigureClient(&ArtifactsClient.Client, o.ResourceManagerAuthorizer)

	PublishedArtifactsClient := blueprint.NewPublishedArtifactsClient()
	o.ConfigureClient(&PublishedArtifactsClient.Client, o.ResourceManagerAuthorizer)

	PublishedBlueprintsClient := blueprint.NewPublishedBlueprintsClient()
	o.ConfigureClient(&PublishedBlueprintsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		BlueprintsClient:          &BlueprintsClient,
		AssignmentsClient:         &AssignmentsClient,
		ArtifactsClient:           &ArtifactsClient,
		PublishedArtifactsClient:  &PublishedArtifactsClient,
		PublishedBlueprintsClient: &PublishedBlueprintsClient,
	}
}
