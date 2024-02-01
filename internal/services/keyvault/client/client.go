// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2023-02-01/vaults"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	dataplane "github.com/tombuildsstuff/kermit/sdk/keyvault/7.4/keyvault"
)

type Client struct {
	// NOTE: Key Vault and Managed HSMs are /intentionally/ split into two different service packages
	// whilst the service shares a similar interface - the behaviours and functionalities of the service
	// including the casing that is required to be used for the constants - differs between the two
	// services.
	//
	// As such this separation on our side is intentional to avoid code reuse given these differences.

	VaultsClient *vaults.VaultsClient

	ManagementClient *dataplane.BaseClient
}

func NewClient(o *common.ClientOptions) *Client {
	managementClient := dataplane.New()
	o.ConfigureClient(&managementClient.Client, o.KeyVaultAuthorizer)

	vaultsClient := vaults.NewVaultsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&vaultsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ManagementClient: &managementClient,
		VaultsClient:     &vaultsClient,
	}
}
