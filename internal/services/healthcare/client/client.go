package client

import (
	"github.com/Azure/azure-sdk-for-go/services/healthcareapis/mgmt/2021-11-01/healthcareapis" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	HealthcareServiceClient                                *healthcareapis.ServicesClient
	HealthcareWorkspaceClient                              *healthcareapis.WorkspacesClient
	HealthcareWorkspaceDicomServiceClient                  *healthcareapis.DicomServicesClient
	HealthcareWorkspaceFhirServiceClient                   *healthcareapis.FhirServicesClient
	HealthcareWorkspaceMedTechServiceClient                *healthcareapis.IotConnectorsClient
	HealthcareWorkspaceMedTechServiceFhirDestinationClient *healthcareapis.IotConnectorFhirDestinationClient
}

func NewClient(o *common.ClientOptions) *Client {
	HealthcareServiceClient := healthcareapis.NewServicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&HealthcareServiceClient.Client, o.ResourceManagerAuthorizer)

	HealthcareWorkspaceClient := healthcareapis.NewWorkspacesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&HealthcareWorkspaceClient.Client, o.ResourceManagerAuthorizer)

	HealthcareWorkspaceDicomServiceClient := healthcareapis.NewDicomServicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&HealthcareWorkspaceDicomServiceClient.Client, o.ResourceManagerAuthorizer)

	HealthcareWorkspaceFhirServiceClient := healthcareapis.NewFhirServicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&HealthcareWorkspaceFhirServiceClient.Client, o.ResourceManagerAuthorizer)

	HealthcareWorkspaceMedTechServiceClient := healthcareapis.NewIotConnectorsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&HealthcareWorkspaceMedTechServiceClient.Client, o.ResourceManagerAuthorizer)

	HealthcareWorkspaceMedTechServiceFhirDestinationClient := healthcareapis.NewIotConnectorFhirDestinationClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&HealthcareWorkspaceMedTechServiceFhirDestinationClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		HealthcareServiceClient:                                &HealthcareServiceClient,
		HealthcareWorkspaceClient:                              &HealthcareWorkspaceClient,
		HealthcareWorkspaceDicomServiceClient:                  &HealthcareWorkspaceDicomServiceClient,
		HealthcareWorkspaceFhirServiceClient:                   &HealthcareWorkspaceFhirServiceClient,
		HealthcareWorkspaceMedTechServiceClient:                &HealthcareWorkspaceMedTechServiceClient,
		HealthcareWorkspaceMedTechServiceFhirDestinationClient: &HealthcareWorkspaceMedTechServiceFhirDestinationClient,
	}
}
