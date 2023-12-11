package blobservice

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BlobServicePropertiesProperties struct {
	AutomaticSnapshotPolicyEnabled *bool                         `json:"automaticSnapshotPolicyEnabled,omitempty"`
	ChangeFeed                     *ChangeFeed                   `json:"changeFeed,omitempty"`
	ContainerDeleteRetentionPolicy *DeleteRetentionPolicy        `json:"containerDeleteRetentionPolicy,omitempty"`
	Cors                           *CorsRules                    `json:"cors,omitempty"`
	DefaultServiceVersion          *string                       `json:"defaultServiceVersion,omitempty"`
	DeleteRetentionPolicy          *DeleteRetentionPolicy        `json:"deleteRetentionPolicy,omitempty"`
	IsVersioningEnabled            *bool                         `json:"isVersioningEnabled,omitempty"`
	LastAccessTimeTrackingPolicy   *LastAccessTimeTrackingPolicy `json:"lastAccessTimeTrackingPolicy,omitempty"`
	RestorePolicy                  *RestorePolicyProperties      `json:"restorePolicy,omitempty"`
}
