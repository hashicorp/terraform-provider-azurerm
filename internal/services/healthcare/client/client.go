package client

import (
	"github.com/Azure/azure-sdk-for-go/services/healthcareapis/mgmt/2021-11-01/healthcareapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	HealthcareServiceClient   *healthcareapis.ServicesClient
	HealthcareWorkspaceClient *healthcareapis.WorkspacesClient
}

func NewClient(o *common.ClientOptions) *Client {
	HealthcareServiceClient := healthcareapis.NewServicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&HealthcareServiceClient.Client, o.ResourceManagerAuthorizer)

	HealthcareWorkspaceClient := healthcareapis.NewWorkspacesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&HealthcareWorkspaceClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		HealthcareServiceClient:   &HealthcareServiceClient,
		HealthcareWorkspaceClient: &HealthcareWorkspaceClient,
	}
}
