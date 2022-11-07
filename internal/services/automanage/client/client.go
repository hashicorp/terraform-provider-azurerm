package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/tombuildsstuff/kermit/sdk/automanage/2022-05-04/automanage"
)

type Client struct {
	ConfigurationProfileClient               *automanage.ConfigurationProfilesClient
	ConfigurationProfilesVersionClient       *automanage.ConfigurationProfilesVersionsClient
	ConfigurationProfileAssignmentClient     *automanage.ConfigurationProfileAssignmentsClient
	ConfigurationProfileHCRPAssignmentClient *automanage.ConfigurationProfileHCRPAssignmentsClient
	ConfigurationProfileHCIAssignmentClient  *automanage.ConfigurationProfileHCIAssignmentsClient
}

func NewClient(o *common.ClientOptions) *Client {
	configurationProfileClient := automanage.NewConfigurationProfilesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&configurationProfileClient.Client, o.ResourceManagerAuthorizer)

	configurationProfilesVersionClient := automanage.NewConfigurationProfilesVersionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&configurationProfilesVersionClient.Client, o.ResourceManagerAuthorizer)

	configurationProfileAssignmentClient := automanage.NewConfigurationProfileAssignmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&configurationProfileAssignmentClient.Client, o.ResourceManagerAuthorizer)

	configurationProfileHCRPAssignmentClient := automanage.NewConfigurationProfileHCRPAssignmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&configurationProfileHCRPAssignmentClient.Client, o.ResourceManagerAuthorizer)

	configurationProfileHCIAssignmentClient := automanage.NewConfigurationProfileHCIAssignmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&configurationProfileHCIAssignmentClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ConfigurationProfileClient:               &configurationProfileClient,
		ConfigurationProfilesVersionClient:       &configurationProfilesVersionClient,
		ConfigurationProfileAssignmentClient:     &configurationProfileAssignmentClient,
		ConfigurationProfileHCRPAssignmentClient: &configurationProfileHCRPAssignmentClient,
		ConfigurationProfileHCIAssignmentClient:  &configurationProfileHCIAssignmentClient,
	}
}
