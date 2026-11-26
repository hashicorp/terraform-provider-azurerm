package serverdevopsaudit

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerDevOpsAuditSettingsProperties struct {
	IsAzureMonitorTargetEnabled  *bool                   `json:"isAzureMonitorTargetEnabled,omitempty"`
	IsManagedIdentityInUse       *bool                   `json:"isManagedIdentityInUse,omitempty"`
	State                        BlobAuditingPolicyState `json:"state"`
	StorageAccountAccessKey      *string                 `json:"storageAccountAccessKey,omitempty"`
	StorageAccountSubscriptionId *string                 `json:"storageAccountSubscriptionId,omitempty"`
	StorageEndpoint              *string                 `json:"storageEndpoint,omitempty"`
}
