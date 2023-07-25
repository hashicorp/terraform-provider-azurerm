// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2022-07-01-preview/configurationassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2022-07-01-preview/maintenanceconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2022-07-01-preview/publicmaintenanceconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ConfigurationsClient           *maintenanceconfigurations.MaintenanceConfigurationsClient
	ConfigurationAssignmentsClient *configurationassignments.ConfigurationAssignmentsClient
	PublicConfigurationsClient     *publicmaintenanceconfigurations.PublicMaintenanceConfigurationsClient
}

func NewClient(o *common.ClientOptions) *Client {
	configurationsClient := maintenanceconfigurations.NewMaintenanceConfigurationsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&configurationsClient.Client, o.ResourceManagerAuthorizer)

	configurationAssignmentsClient := configurationassignments.NewConfigurationAssignmentsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&configurationAssignmentsClient.Client, o.ResourceManagerAuthorizer)

	publicConfigurationsClient := publicmaintenanceconfigurations.NewPublicMaintenanceConfigurationsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&publicConfigurationsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ConfigurationsClient:           &configurationsClient,
		ConfigurationAssignmentsClient: &configurationAssignmentsClient,
		PublicConfigurationsClient:     &publicConfigurationsClient,
	}
}
