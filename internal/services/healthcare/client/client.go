package client

import (
	healthcare "github.com/Azure/azure-sdk-for-go/services/healthcareapis/mgmt/2020-03-30/healthcareapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	healthcareWorkspaceDicom "github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/sdk/2021-06-01-preview/dicomservices"
	healthcareWorkspaceFhirService "github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/sdk/2021-06-01-preview/fhirservices"
	healthcareWorkspaceIotConnector "github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/sdk/2021-06-01-preview/iotconnectors"
	healthcareWorkspace "github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/sdk/2021-06-01-preview/workspaces"
)

type Client struct {
	HealthcareServiceClient               *healthcare.ServicesClient
	HealthcareWorkspaceClient             *healthcareWorkspace.WorkspacesClient
	HealthcareWorkspaceFhirServiceClient  *healthcareWorkspaceFhirService.FhirServicesClient
	HealthcareWorkspaceDicomServiceClient *healthcareWorkspaceDicom.DicomServicesClient
	HealthcareWorkspaceIotConnectorClient *healthcareWorkspaceIotConnector.IotConnectorsClient
}

func NewClient(o *common.ClientOptions) *Client {
	HealthcareServiceClient := healthcare.NewServicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&HealthcareServiceClient.Client, o.ResourceManagerAuthorizer)

	HealthcareWorkspaceClient := healthcareWorkspace.NewWorkspacesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&HealthcareWorkspaceClient.Client, o.ResourceManagerAuthorizer)

	HealthcareWorkspaceFhirServiceClient := healthcareWorkspaceFhirService.NewFhirServicesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&HealthcareWorkspaceFhirServiceClient.Client, o.ResourceManagerAuthorizer)

	HealthcareWorkspaceIotConnectorClient := healthcareWorkspaceIotConnector.NewIotConnectorsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&HealthcareWorkspaceFhirServiceClient.Client, o.ResourceManagerAuthorizer)

	HealthcareWorkspaceDicomServiceClient := healthcareWorkspaceDicom.NewDicomServicesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&HealthcareWorkspaceFhirServiceClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		HealthcareServiceClient:               &HealthcareServiceClient,
		HealthcareWorkspaceClient:             &HealthcareWorkspaceClient,
		HealthcareWorkspaceFhirServiceClient:  &HealthcareWorkspaceFhirServiceClient,
		HealthcareWorkspaceIotConnectorClient: &HealthcareWorkspaceIotConnectorClient,
		HealthcareWorkspaceDicomServiceClient: &HealthcareWorkspaceDicomServiceClient,
	}
}
