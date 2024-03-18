// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofilehciassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/tombuildsstuff/kermit/sdk/automanage/2022-05-04/automanage"
)

type Client struct {
	ConfigurationProfilesClient              *configurationprofiles.ConfigurationProfilesClient
	ConfigurationProfileHCIAssignmentsClient *configurationprofilehciassignments.ConfigurationProfileHCIAssignmentsClient

	// NOTE: these clients use `tombuildsstuff/kermit` (a variant of `Azure/azure-sdk-for-go`) and shouldn't be used going forwards.
	ConfigurationClient *automanage.ConfigurationProfilesClient
	HCIAssignmentClient *automanage.ConfigurationProfileHCIAssignmentsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	configurationProfilesClient, err := configurationprofiles.NewConfigurationProfilesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ConfigurationProfiles client: %+v", err)
	}
	o.Configure(configurationProfilesClient.Client, o.Authorizers.ResourceManager)

	configurationProfileHCIAssignmentsClient, err := configurationprofilehciassignments.NewConfigurationProfileHCIAssignmentsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ConfigurationProfilesHCIAssignments client: %+v", err)
	}
	o.Configure(configurationProfileHCIAssignmentsClient.Client, o.Authorizers.ResourceManager)

	// NOTE: these clients use `tombuildsstuff/kermit` (a variant of `Azure/azure-sdk-for-go`) and shouldn't be used going forwards.
	legacyConfigurationProfilesClient := automanage.NewConfigurationProfilesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&legacyConfigurationProfilesClient.Client, o.ResourceManagerAuthorizer)
	legacyHCIAssignmentsClient := automanage.NewConfigurationProfileHCIAssignmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&legacyHCIAssignmentsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ConfigurationProfilesClient:              configurationProfilesClient,
		ConfigurationProfileHCIAssignmentsClient: configurationProfileHCIAssignmentsClient,

		ConfigurationClient: &legacyConfigurationProfilesClient,
		HCIAssignmentClient: &legacyHCIAssignmentsClient,
	}, nil
}
