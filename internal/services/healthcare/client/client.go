package client

import (
	healthcare "github.com/Azure/azure-sdk-for-go/services/healthcareapis/mgmt/2020-03-30/healthcareapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	healthcareWorkspaceIotConnector "github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/sdk/2021-06-01-preview/iotconnectors"
)

type Client struct {
	HealthcareServiceClient               *healthcare.ServicesClient
	HealthcareWorkspaceIotConnectorClient *healthcareWorkspaceIotConnector.IotConnectorsClient
}

func NewClient(o *common.ClientOptions) *Client {
	HealthcareServiceClient := healthcare.NewServicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&HealthcareServiceClient.Client, o.ResourceManagerAuthorizer)

	HealthcareWorkspaceIotConnectorClient := healthcareWorkspaceIotConnector.NewIotConnectorsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&HealthcareWorkspaceIotConnectorClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		HealthcareServiceClient:               &HealthcareServiceClient,
		HealthcareWorkspaceIotConnectorClient: &HealthcareWorkspaceIotConnectorClient,
	}
}
