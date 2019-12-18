package client

import (
	healthcare "github.com/Azure/azure-sdk-for-go/services/healthcareapis/mgmt/2019-09-16/healthcareapis"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	HealthcareServiceClient *healthcare.ServicesClient
}

func NewClient(o *common.ClientOptions) *Client {
	HealthcareServiceClient := healthcare.NewServicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&HealthcareServiceClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		HealthcareServiceClient: &HealthcareServiceClient,
	}
}
