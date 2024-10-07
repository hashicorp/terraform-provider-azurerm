// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2023-02-01/vaults"
	vaults20230701 "github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2023-07-01/vaults"
	resources20151101 "github.com/hashicorp/go-azure-sdk/resource-manager/resources/2015-11-01/resources"
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

	ManagementClient *dataplane.BaseClient // TODO: we should rename this DataPlaneClient in time

	// NOTE: @tombuildsstuff: this client is intentionally internal-only so that it's not used directly
	resources20151101Client *resources20151101.ResourcesClient

	// @tombuildsstuff: I'm intentionally vendoring this API version separately to take advantage
	// of the updated List API behaviour/new base layer (since the API now always returns a nextLink)
	// which `Azure/go-autorest` doesn't handle cleanly. Before migrating the Key Vault resources over
	// to using this new API Version, there's a few Enum/casing related items to resolve, and as such
	// that's intentionally split-out into a separate task. For now, please continue to use VaultsClient
	// for regular operations, and we can remove this internal client one the newer API version is used
	// across the Provider.
	vaults20230701Client *vaults20230701.VaultsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	// These clients use `hashicorp/go-azure-sdk`
	updatedVaultsClient, err := vaults20230701.NewVaultsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Vaults client: %+v", err)
	}
	o.Configure(updatedVaultsClient.Client, o.Authorizers.ResourceManager)

	resources20151101Client, err := resources20151101.NewResourcesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building legacy Resources client: %+v", err)
	}
	o.Configure(resources20151101Client.Client, o.Authorizers.ResourceManager)

	// These clients use `Azure/azure-sdk-for-go` and/or `Azure/go-autorest`
	vaultsClient := vaults.NewVaultsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&vaultsClient.Client, o.ResourceManagerAuthorizer)
	managementClient := dataplane.New()
	o.ConfigureClient(&managementClient.Client, o.KeyVaultAuthorizer)

	return &Client{
		ManagementClient: &managementClient,
		VaultsClient:     &vaultsClient,

		// intentionally internal to this package for now, see above.
		resources20151101Client: resources20151101Client,
		vaults20230701Client:    updatedVaultsClient,
	}, nil
}
