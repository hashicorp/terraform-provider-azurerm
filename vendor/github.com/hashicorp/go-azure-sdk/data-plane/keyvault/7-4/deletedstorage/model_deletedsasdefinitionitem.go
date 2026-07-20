package deletedstorage

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeletedSasDefinitionItem struct {
	Attributes         *SasDefinitionAttributes `json:"attributes,omitempty"`
	DeletedDate        *int64                   `json:"deletedDate,omitempty"`
	Id                 *string                  `json:"id,omitempty"`
	RecoveryId         *string                  `json:"recoveryId,omitempty"`
	ScheduledPurgeDate *int64                   `json:"scheduledPurgeDate,omitempty"`
	Sid                *string                  `json:"sid,omitempty"`
	Tags               *map[string]string       `json:"tags,omitempty"`
}
