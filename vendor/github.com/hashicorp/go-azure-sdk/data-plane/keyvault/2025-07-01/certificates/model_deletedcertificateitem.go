package certificates

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeletedCertificateItem struct {
	Attributes         *CertificateAttributes `json:"attributes,omitempty"`
	DeletedDate        *int64                 `json:"deletedDate,omitempty"`
	Id                 *string                `json:"id,omitempty"`
	RecoveryId         *string                `json:"recoveryId,omitempty"`
	ScheduledPurgeDate *int64                 `json:"scheduledPurgeDate,omitempty"`
	Tags               *map[string]string     `json:"tags,omitempty"`
	X5t                *string                `json:"x5t,omitempty"`
}
