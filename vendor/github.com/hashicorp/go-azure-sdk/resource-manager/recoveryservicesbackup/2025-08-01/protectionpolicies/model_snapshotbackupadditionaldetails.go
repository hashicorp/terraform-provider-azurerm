package protectionpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SnapshotBackupAdditionalDetails struct {
	InstantRPDetails                   *string                             `json:"instantRPDetails,omitempty"`
	InstantRpRetentionRangeInDays      *int64                              `json:"instantRpRetentionRangeInDays,omitempty"`
	UserAssignedManagedIdentityDetails *UserAssignedManagedIdentityDetails `json:"userAssignedManagedIdentityDetails,omitempty"`
}
