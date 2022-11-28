package blobservice

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BlobServicePropertiesProperties struct {
	AutomaticSnapshotPolicyEnabled *bool                         `json:"automaticSnapshotPolicyEnabled,omitempty"`
	ChangeFeed                     *ChangeFeed                   `json:"changeFeed"`
	ContainerDeleteRetentionPolicy *DeleteRetentionPolicy        `json:"containerDeleteRetentionPolicy"`
	Cors                           *CorsRules                    `json:"cors"`
	DefaultServiceVersion          *string                       `json:"defaultServiceVersion,omitempty"`
	DeleteRetentionPolicy          *DeleteRetentionPolicy        `json:"deleteRetentionPolicy"`
	IsVersioningEnabled            *bool                         `json:"isVersioningEnabled,omitempty"`
	LastAccessTimeTrackingPolicy   *LastAccessTimeTrackingPolicy `json:"lastAccessTimeTrackingPolicy"`
	RestorePolicy                  *RestorePolicyProperties      `json:"restorePolicy"`
}
