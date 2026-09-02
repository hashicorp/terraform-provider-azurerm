package vaults

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SoftDeleteSettings struct {
	EnhancedSecurityState           *EnhancedSecurityState `json:"enhancedSecurityState,omitempty"`
	SoftDeleteRetentionPeriodInDays *int64                 `json:"softDeleteRetentionPeriodInDays,omitempty"`
	SoftDeleteState                 *SoftDeleteState       `json:"softDeleteState,omitempty"`
}
