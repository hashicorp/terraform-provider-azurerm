package client

import (
	"github.com/Azure/azure-sdk-for-go/services/maintenance/mgmt/2021-05-01/maintenance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ConfigurationsClient           *maintenance.ConfigurationsClient
	ConfigurationAssignmentsClient *maintenance.ConfigurationAssignmentsClient
	PublicConfigurationsClient     *maintenance.PublicMaintenanceConfigurationsClient
}

func NewClient(o *common.ClientOptions) *Client {
	configurationsClient := maintenance.NewConfigurationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&configurationsClient.Client, o.ResourceManagerAuthorizer)

	configurationAssignmentsClient := maintenance.NewConfigurationAssignmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&configurationAssignmentsClient.Client, o.ResourceManagerAuthorizer)

	publicConfigurationsClient := maintenance.NewPublicMaintenanceConfigurationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&publicConfigurationsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ConfigurationsClient:           &configurationsClient,
		ConfigurationAssignmentsClient: &configurationAssignmentsClient,
		PublicConfigurationsClient:     &publicConfigurationsClient,
	}
}
