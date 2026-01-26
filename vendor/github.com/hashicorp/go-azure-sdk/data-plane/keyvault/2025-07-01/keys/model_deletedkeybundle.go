package keys

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeletedKeyBundle struct {
	Attributes         *KeyAttributes     `json:"attributes,omitempty"`
	DeletedDate        *int64             `json:"deletedDate,omitempty"`
	Key                *JsonWebKey        `json:"key,omitempty"`
	Managed            *bool              `json:"managed,omitempty"`
	RecoveryId         *string            `json:"recoveryId,omitempty"`
	ReleasePolicy      *KeyReleasePolicy  `json:"release_policy,omitempty"`
	ScheduledPurgeDate *int64             `json:"scheduledPurgeDate,omitempty"`
	Tags               *map[string]string `json:"tags,omitempty"`
}
