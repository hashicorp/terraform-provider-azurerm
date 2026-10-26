package blobauditing

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExtendedServerBlobAuditingPolicyProperties struct {
	AuditActionsAndGroups        *[]string               `json:"auditActionsAndGroups,omitempty"`
	IsAzureMonitorTargetEnabled  *bool                   `json:"isAzureMonitorTargetEnabled,omitempty"`
	IsDevopsAuditEnabled         *bool                   `json:"isDevopsAuditEnabled,omitempty"`
	IsManagedIdentityInUse       *bool                   `json:"isManagedIdentityInUse,omitempty"`
	IsStorageSecondaryKeyInUse   *bool                   `json:"isStorageSecondaryKeyInUse,omitempty"`
	PredicateExpression          *string                 `json:"predicateExpression,omitempty"`
	QueueDelayMs                 *int64                  `json:"queueDelayMs,omitempty"`
	RetentionDays                *int64                  `json:"retentionDays,omitempty"`
	State                        BlobAuditingPolicyState `json:"state"`
	StorageAccountAccessKey      *string                 `json:"storageAccountAccessKey,omitempty"`
	StorageAccountSubscriptionId *string                 `json:"storageAccountSubscriptionId,omitempty"`
	StorageEndpoint              *string                 `json:"storageEndpoint,omitempty"`
}
