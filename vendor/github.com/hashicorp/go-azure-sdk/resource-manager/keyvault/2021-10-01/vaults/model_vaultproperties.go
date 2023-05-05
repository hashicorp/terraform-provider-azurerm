package vaults

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VaultProperties struct {
	AccessPolicies               *[]AccessPolicyEntry             `json:"accessPolicies,omitempty"`
	CreateMode                   *CreateMode                      `json:"createMode,omitempty"`
	EnablePurgeProtection        *bool                            `json:"enablePurgeProtection,omitempty"`
	EnableRbacAuthorization      *bool                            `json:"enableRbacAuthorization,omitempty"`
	EnableSoftDelete             *bool                            `json:"enableSoftDelete,omitempty"`
	EnabledForDeployment         *bool                            `json:"enabledForDeployment,omitempty"`
	EnabledForDiskEncryption     *bool                            `json:"enabledForDiskEncryption,omitempty"`
	EnabledForTemplateDeployment *bool                            `json:"enabledForTemplateDeployment,omitempty"`
	HsmPoolResourceId            *string                          `json:"hsmPoolResourceId,omitempty"`
	NetworkAcls                  *NetworkRuleSet                  `json:"networkAcls,omitempty"`
	PrivateEndpointConnections   *[]PrivateEndpointConnectionItem `json:"privateEndpointConnections,omitempty"`
	ProvisioningState            *VaultProvisioningState          `json:"provisioningState,omitempty"`
	PublicNetworkAccess          *string                          `json:"publicNetworkAccess,omitempty"`
	Sku                          Sku                              `json:"sku"`
	SoftDeleteRetentionInDays    *int64                           `json:"softDeleteRetentionInDays,omitempty"`
	TenantId                     string                           `json:"tenantId"`
	VaultUri                     *string                          `json:"vaultUri,omitempty"`
}
