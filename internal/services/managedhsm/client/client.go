// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2023-02-01/managedhsms"
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

	// TODO: rename these once this is ported over

	ManagedHsmClient          *managedhsms.ManagedHsmsClient
	ManagementClient          *dataplane.BaseClient
	MHSMSDClient              *dataplane.HSMSecurityDomainClient
	MHSMRoleClient            *dataplane.RoleDefinitionsClient
	MHSMRoleAssignmentsClient *dataplane.RoleAssignmentsClient
}

func NewClient(o *common.ClientOptions) *Client {
	managedHsmClient := managedhsms.NewManagedHsmsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&managedHsmClient.Client, o.ResourceManagerAuthorizer)

	managementClient := dataplane.New()
	o.ConfigureClient(&managementClient.Client, o.KeyVaultAuthorizer)

	sdClient := dataplane.NewHSMSecurityDomainClient()
	o.ConfigureClient(&sdClient.Client, o.ManagedHSMAuthorizer)

	mhsmRoleDefineClient := dataplane.NewRoleDefinitionsClient()
	o.ConfigureClient(&mhsmRoleDefineClient.Client, o.ManagedHSMAuthorizer)

	mhsmRoleAssignClient := dataplane.NewRoleAssignmentsClient()
	o.ConfigureClient(&mhsmRoleAssignClient.Client, o.ManagedHSMAuthorizer)

	return &Client{
		ManagedHsmClient:          &managedHsmClient,
		ManagementClient:          &managementClient,
		MHSMSDClient:              &sdClient,
		MHSMRoleClient:            &mhsmRoleDefineClient,
		MHSMRoleAssignmentsClient: &mhsmRoleAssignClient,
	}
}
