// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	dataplane7_4 "github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2026-02-01/deletedvaults"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2026-02-01/vaults"
	resources20151101 "github.com/hashicorp/go-azure-sdk/resource-manager/resources/2015-11-01/resources"
	dataplaneClient "github.com/hashicorp/go-azure-sdk/sdk/client/dataplane"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	dataplane "github.com/jackofallops/kermit/sdk/keyvault/7.4/keyvault"
)

type Client struct {
	// NOTE: Key Vault and Managed HSMs are /intentionally/ split into two different service packages
	// whilst the service shares a similar interface - the behaviours and functionalities of the service
	// including the casing that is required to be used for the constants - differs between the two
	// services.
	//
	// As such this separation on our side is intentional to avoid code reuse given these differences.
	VaultsClient        *vaults.VaultsClient
	DeletedVaultsClient *deletedvaults.DeletedVaultsClient

	ManagementClient        *dataplane.BaseClient // TODO: we should rename this DataPlaneClient in time
	DataPlaneKeyVaultClient *dataplane7_4.Client

	// NOTE: @tombuildsstuff: this client is intentionally internal-only so that it's not used directly
	resources20151101Client *resources20151101.ResourcesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	resources20151101Client, err := resources20151101.NewResourcesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building legacy Resources client: %+v", err)
	}
	o.Configure(resources20151101Client.Client, o.Authorizers.ResourceManager)

	vaultsClient, err := vaults.NewVaultsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Vaults client: %+v", err)
	}
	o.Configure(vaultsClient.Client, o.Authorizers.ResourceManager)

	deletedVaultsClient, err := deletedvaults.NewDeletedVaultsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building DeletedVaults client: %+v", err)
	}
	o.Configure(deletedVaultsClient.Client, o.Authorizers.ResourceManager)

	managementClient := dataplane.New()
	o.ConfigureClient(&managementClient.Client, o.KeyVaultAuthorizer)

	dataplaneKeyvaultClient, err := dataplane7_4.NewClient(func(c *dataplaneClient.Client) {
		o.Configure(c.Client, o.Authorizers.KeyVault)
	})
	if err != nil {
		return nil, fmt.Errorf("building data-plane KeyVault client: %+v", err)
	}

	return &Client{
		ManagementClient:    &managementClient,
		VaultsClient:        vaultsClient,
		DeletedVaultsClient: deletedVaultsClient,

		DataPlaneKeyVaultClient: dataplaneKeyvaultClient,

		// intentionally internal to this package for now, see above.
		resources20151101Client: resources20151101Client,
	}, nil
}
