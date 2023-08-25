// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2023-02-01/managedhsms"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2023-02-01/vaults"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	dataplane "github.com/tombuildsstuff/kermit/sdk/keyvault/7.4/keyvault"
)

type Client struct {
	ManagedHsmClient *managedhsms.ManagedHsmsClient
	ManagementClient *dataplane.BaseClient
	VaultsClient     *vaults.VaultsClient

	MHSMSDClient   *dataplane.HSMSecurityDomainClient
	MHSMRoleClient *dataplane.RoleDefinitionsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	managedHsmClient, err := managedhsms.NewManagedHsmsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ManagedHsms client: %+v", err)
	}
	o.Configure(managedHsmClient.Client, o.Authorizers.ResourceManager)

	managementClient := dataplane.New()
	o.ConfigureClient(&managementClient.Client, o.KeyVaultAuthorizer)

	vaultsClient, err := vaults.NewVaultsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Vaults client: %+v", err)
	}
	o.Configure(vaultsClient.Client, o.Authorizers.ResourceManager)

	sdClient := dataplane.NewHSMSecurityDomainClient()
	o.ConfigureClient(&sdClient.Client, o.ManagedHSMAuthorizer)

	mhsmRoleDefineClient := dataplane.NewRoleDefinitionsClient()
	o.ConfigureClient(&mhsmRoleDefineClient.Client, o.ManagedHSMAuthorizer)

	return &Client{
		ManagedHsmClient: managedHsmClient,
		VaultsClient:     vaultsClient,

		ManagementClient: &managementClient,
		MHSMSDClient:     &sdClient,
		MHSMRoleClient:   &mhsmRoleDefineClient,
	}, nil
}
