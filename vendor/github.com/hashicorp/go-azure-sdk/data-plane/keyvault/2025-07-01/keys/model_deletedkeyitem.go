package keys

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeletedKeyItem struct {
	Attributes         *KeyAttributes     `json:"attributes,omitempty"`
	DeletedDate        *int64             `json:"deletedDate,omitempty"`
	Kid                *string            `json:"kid,omitempty"`
	Managed            *bool              `json:"managed,omitempty"`
	RecoveryId         *string            `json:"recoveryId,omitempty"`
	ScheduledPurgeDate *int64             `json:"scheduledPurgeDate,omitempty"`
	Tags               *map[string]string `json:"tags,omitempty"`
}
