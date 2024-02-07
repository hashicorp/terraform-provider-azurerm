// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/tombuildsstuff/kermit/sdk/automanage/2022-05-04/automanage"
)

type Client struct {
	ConfigurationClient  *automanage.ConfigurationProfilesClient
	HCIAssignmentClient  *automanage.ConfigurationProfileHCIAssignmentsClient
	HCRPAssignmentClient *automanage.ConfigurationProfileHCRPAssignmentsClient
	VMAssignmentClient   *automanage.ConfigurationProfileAssignmentsClient
}

func NewClient(o *common.ClientOptions) *Client {
	configurationProfileClient := automanage.NewConfigurationProfilesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&configurationProfileClient.Client, o.ResourceManagerAuthorizer)

	hciAssignmentClient := automanage.NewConfigurationProfileHCIAssignmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&hciAssignmentClient.Client, o.ResourceManagerAuthorizer)

	hcrpAssignmentClient := automanage.NewConfigurationProfileHCRPAssignmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&hcrpAssignmentClient.Client, o.ResourceManagerAuthorizer)

	vmAssignmentClient := automanage.NewConfigurationProfileAssignmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&vmAssignmentClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ConfigurationClient:  &configurationProfileClient,
		HCIAssignmentClient:  &hciAssignmentClient,
		HCRPAssignmentClient: &hcrpAssignmentClient,
		VMAssignmentClient:   &vmAssignmentClient,
	}
}
