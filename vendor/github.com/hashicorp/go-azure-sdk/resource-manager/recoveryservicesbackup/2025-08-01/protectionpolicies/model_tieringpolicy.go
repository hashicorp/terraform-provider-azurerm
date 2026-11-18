package protectionpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TieringPolicy struct {
	Duration     *int64                 `json:"duration,omitempty"`
	DurationType *RetentionDurationType `json:"durationType,omitempty"`
	TieringMode  *TieringMode           `json:"tieringMode,omitempty"`
}
