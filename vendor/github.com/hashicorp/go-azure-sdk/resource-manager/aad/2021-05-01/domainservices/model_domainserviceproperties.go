package domainservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DomainServiceProperties struct {
	ConfigDiagnostics       *ConfigDiagnostics      `json:"configDiagnostics,omitempty"`
	DeploymentId            *string                 `json:"deploymentId,omitempty"`
	DomainConfigurationType *string                 `json:"domainConfigurationType,omitempty"`
	DomainName              *string                 `json:"domainName,omitempty"`
	DomainSecuritySettings  *DomainSecuritySettings `json:"domainSecuritySettings,omitempty"`
	FilteredSync            *FilteredSync           `json:"filteredSync,omitempty"`
	LdapsSettings           *LdapsSettings          `json:"ldapsSettings,omitempty"`
	MigrationProperties     *MigrationProperties    `json:"migrationProperties,omitempty"`
	NotificationSettings    *NotificationSettings   `json:"notificationSettings,omitempty"`
	ProvisioningState       *string                 `json:"provisioningState,omitempty"`
	ReplicaSets             *[]ReplicaSet           `json:"replicaSets,omitempty"`
	ResourceForestSettings  *ResourceForestSettings `json:"resourceForestSettings,omitempty"`
	Sku                     *string                 `json:"sku,omitempty"`
	SyncOwner               *string                 `json:"syncOwner,omitempty"`
	TenantId                *string                 `json:"tenantId,omitempty"`
	Version                 *int64                  `json:"version,omitempty"`
}
