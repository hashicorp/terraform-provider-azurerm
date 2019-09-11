package healthcare

import (
	healthcare "github.com/Azure/azure-sdk-for-go/services/preview/healthcareapis/mgmt/2018-08-20-preview/healthcareapis"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	HealthcareServiceClient *healthcare.ServicesClient
}

func BuildClient(o *common.ClientOptions) *Client {

	HealthcareServiceClient := healthcare.NewServicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&HealthcareServiceClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		HealthcareServiceClient: &HealthcareServiceClient,
	}
}
