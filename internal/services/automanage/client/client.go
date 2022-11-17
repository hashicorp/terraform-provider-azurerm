package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/tombuildsstuff/kermit/sdk/automanage/2022-05-04/automanage"
)

type Client struct {
	ConfigurationProfileClient              *automanage.ConfigurationProfilesClient
	ConfigurationProfileAssignmentClient    *automanage.ConfigurationProfileAssignmentsClient
	ConfigurationProfileHCIAssignmentClient *automanage.ConfigurationProfileHCIAssignmentsClient
}

func NewClient(o *common.ClientOptions) *Client {
	configurationProfileClient := automanage.NewConfigurationProfilesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&configurationProfileClient.Client, o.ResourceManagerAuthorizer)

	configurationProfileAssignmentClient := automanage.NewConfigurationProfileAssignmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&configurationProfileAssignmentClient.Client, o.ResourceManagerAuthorizer)

	configurationProfileHCIAssignmentClient := automanage.NewConfigurationProfileHCIAssignmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&configurationProfileHCIAssignmentClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ConfigurationProfileClient:              &configurationProfileClient,
		ConfigurationProfileAssignmentClient:    &configurationProfileAssignmentClient,
		ConfigurationProfileHCIAssignmentClient: &configurationProfileHCIAssignmentClient,
	}
}
