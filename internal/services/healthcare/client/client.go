package client

import (
	healthcare "github.com/Azure/azure-sdk-for-go/services/healthcareapis/mgmt/2020-03-30/healthcareapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	healthcareWorkspaceDicom "github.com/Azure/azure-sdk-for-go/services/healthcareapis/mgmt/2021-11-01/healthcareapis"

)

type Client struct {
	HealthcareServiceClient *healthcare.ServicesClient
	HealthcareWorkspaceDicomServiceClient *healthcareWorkspaceDicom.DicomServicesClient
}

func NewClient(o *common.ClientOptions) *Client {
	HealthcareServiceClient := healthcare.NewServicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&HealthcareServiceClient.Client, o.ResourceManagerAuthorizer)

	HealthcareWorkspaceDicomServiceClient := healthcareWorkspaceDicom.NewDicomServicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&HealthcareWorkspaceDicomServiceClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		HealthcareServiceClient: &HealthcareServiceClient,
		HealthcareWorkspaceDicomServiceClient: &HealthcareWorkspaceDicomServiceClient,
	}
}
