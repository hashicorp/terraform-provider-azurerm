package storage

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeletedSasDefinitionBundle struct {
	Attributes         *SasDefinitionAttributes `json:"attributes,omitempty"`
	DeletedDate        *int64                   `json:"deletedDate,omitempty"`
	Id                 *string                  `json:"id,omitempty"`
	RecoveryId         *string                  `json:"recoveryId,omitempty"`
	SasType            *SasTokenType            `json:"sasType,omitempty"`
	ScheduledPurgeDate *int64                   `json:"scheduledPurgeDate,omitempty"`
	Sid                *string                  `json:"sid,omitempty"`
	Tags               *map[string]string       `json:"tags,omitempty"`
	TemplateUri        *string                  `json:"templateUri,omitempty"`
	ValidityPeriod     *string                  `json:"validityPeriod,omitempty"`
}
