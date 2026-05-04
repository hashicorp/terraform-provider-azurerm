package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeletedAppRestoreRequestProperties struct {
	DeletedSiteId        *string `json:"deletedSiteId,omitempty"`
	RecoverConfiguration *bool   `json:"recoverConfiguration,omitempty"`
	SnapshotTime         *string `json:"snapshotTime,omitempty"`
	UseDRSecondary       *bool   `json:"useDRSecondary,omitempty"`
}
