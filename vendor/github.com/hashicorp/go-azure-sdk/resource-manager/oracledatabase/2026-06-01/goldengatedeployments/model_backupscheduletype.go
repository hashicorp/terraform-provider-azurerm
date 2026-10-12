package goldengatedeployments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupScheduleType struct {
	BucketName               *string        `json:"bucketName,omitempty"`
	CompartmentId            *string        `json:"compartmentId,omitempty"`
	FrequencyBackupScheduled *FrequencyType `json:"frequencyBackupScheduled,omitempty"`
	IsMetadataOnly           *bool          `json:"isMetadataOnly,omitempty"`
	NamespaceName            *string        `json:"namespaceName,omitempty"`
	TimeBackupScheduled      *string        `json:"timeBackupScheduled,omitempty"`
}
