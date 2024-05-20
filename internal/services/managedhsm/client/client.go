// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2023-07-01/managedhsms"
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

	// Resource Manager
	ManagedHsmClient *managedhsms.ManagedHsmsClient

	// Data Plane
	DataPlaneKeysClient            *dataplane.BaseClient
	DataPlaneRoleAssignmentsClient *dataplane.RoleAssignmentsClient
	DataPlaneRoleDefinitionsClient *dataplane.RoleDefinitionsClient
	DataPlaneSecurityDomainsClient *dataplane.HSMSecurityDomainClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	managedHsmClient, err := managedhsms.NewManagedHsmsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ManagedHsms client: %+v", err)
	}
	o.Configure(managedHsmClient.Client, o.Authorizers.ResourceManager)

	managementKeysClient := dataplane.New()
	o.ConfigureClient(&managementKeysClient.Client, o.ManagedHSMAuthorizer)

	securityDomainClient := dataplane.NewHSMSecurityDomainClient()
	o.ConfigureClient(&securityDomainClient.Client, o.ManagedHSMAuthorizer)

	roleDefinitionsClient := dataplane.NewRoleDefinitionsClient()
	o.ConfigureClient(&roleDefinitionsClient.Client, o.ManagedHSMAuthorizer)

	roleAssignmentsClient := dataplane.NewRoleAssignmentsClient()
	o.ConfigureClient(&roleAssignmentsClient.Client, o.ManagedHSMAuthorizer)

	return &Client{
		// Resource Manger
		ManagedHsmClient: managedHsmClient,

		// Data Plane
		DataPlaneKeysClient:            &managementKeysClient,
		DataPlaneSecurityDomainsClient: &securityDomainClient,
		DataPlaneRoleDefinitionsClient: &roleDefinitionsClient,
		DataPlaneRoleAssignmentsClient: &roleAssignmentsClient,
	}, nil
}
