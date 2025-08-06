// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/fhirservices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/iotconnectors"
	service "github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/resource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2024-03-31/dicomservices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2024-03-31/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	HealthcareServiceClient                *service.ResourceClient
	HealthcareWorkspaceClient              *workspaces.WorkspacesClient
	HealthcareWorkspaceDicomServiceClient  *dicomservices.DicomServicesClient
	HealthcareWorkspaceFhirServiceClient   *fhirservices.FhirServicesClient
	HealthcareWorkspaceIotConnectorsClient *iotconnectors.IotConnectorsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	healthcareServiceClient, err := service.NewResourceClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building HealthcareService Client: %+v", err)
	}
	o.Configure(healthcareServiceClient.Client, o.Authorizers.ResourceManager)

	healthcareWorkspaceClient, err := workspaces.NewWorkspacesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building HealthcareWorkspace Client: %+v", err)
	}
	o.Configure(healthcareWorkspaceClient.Client, o.Authorizers.ResourceManager)

	healthcareWorkspaceDicomServiceClient, err := dicomservices.NewDicomServicesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building HealthcareWorkspaceDicomService Client: %+v", err)
	}
	o.Configure(healthcareWorkspaceDicomServiceClient.Client, o.Authorizers.ResourceManager)

	healthcareWorkspaceFhirServiceClient, err := fhirservices.NewFhirServicesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building HealthcareWorkspaceFhirService Client: %+v", err)
	}
	o.Configure(healthcareWorkspaceFhirServiceClient.Client, o.Authorizers.ResourceManager)

	healthcareWorkspaceIotConnectorsClient, err := iotconnectors.NewIotConnectorsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building HealthcareWorkspaceIotConnectors Client: %+v", err)
	}
	o.Configure(healthcareWorkspaceIotConnectorsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		HealthcareServiceClient:                healthcareServiceClient,
		HealthcareWorkspaceClient:              healthcareWorkspaceClient,
		HealthcareWorkspaceDicomServiceClient:  healthcareWorkspaceDicomServiceClient,
		HealthcareWorkspaceFhirServiceClient:   healthcareWorkspaceFhirServiceClient,
		HealthcareWorkspaceIotConnectorsClient: healthcareWorkspaceIotConnectorsClient,
	}, nil
}
