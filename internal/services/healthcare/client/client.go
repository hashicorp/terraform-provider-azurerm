package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/dicomservices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/fhirservices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/iotconnectors"
	service "github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/resource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	HealthcareServiceClient                *service.ResourceClient
	HealthcareWorkspaceClient              *workspaces.WorkspacesClient
	HealthcareWorkspaceDicomServiceClient  *dicomservices.DicomServicesClient
	HealthcareWorkspaceFhirServiceClient   *fhirservices.FhirServicesClient
	HealthcareWorkspaceIotConnectorsClient *iotconnectors.IotConnectorsClient
}

func NewClient(o *common.ClientOptions) *Client {
	HealthcareServiceClient := service.NewResourceClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&HealthcareServiceClient.Client, o.ResourceManagerAuthorizer)

	HealthcareWorkspaceClient := workspaces.NewWorkspacesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&HealthcareWorkspaceClient.Client, o.ResourceManagerAuthorizer)

	HealthcareWorkspaceDicomServiceClient := dicomservices.NewDicomServicesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&HealthcareWorkspaceDicomServiceClient.Client, o.ResourceManagerAuthorizer)

	HealthcareWorkspaceFhirServiceClient := fhirservices.NewFhirServicesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&HealthcareWorkspaceFhirServiceClient.Client, o.ResourceManagerAuthorizer)

	HealthcareWorkspaceIotConnectorsClient := iotconnectors.NewIotConnectorsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&HealthcareWorkspaceIotConnectorsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		HealthcareServiceClient:                &HealthcareServiceClient,
		HealthcareWorkspaceClient:              &HealthcareWorkspaceClient,
		HealthcareWorkspaceDicomServiceClient:  &HealthcareWorkspaceDicomServiceClient,
		HealthcareWorkspaceFhirServiceClient:   &HealthcareWorkspaceFhirServiceClient,
		HealthcareWorkspaceIotConnectorsClient: &HealthcareWorkspaceIotConnectorsClient,
	}
}
