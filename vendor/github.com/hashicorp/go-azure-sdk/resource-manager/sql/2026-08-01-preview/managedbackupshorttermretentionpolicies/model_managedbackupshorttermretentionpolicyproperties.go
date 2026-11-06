package managedbackupshorttermretentionpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedBackupShortTermRetentionPolicyProperties struct {
	ImmutabilityStatus *ImmutabilityStatus `json:"immutabilityStatus,omitempty"`
	LockImmutability   *bool               `json:"lockImmutability,omitempty"`
	RetentionDays      *int64              `json:"retentionDays,omitempty"`
}
