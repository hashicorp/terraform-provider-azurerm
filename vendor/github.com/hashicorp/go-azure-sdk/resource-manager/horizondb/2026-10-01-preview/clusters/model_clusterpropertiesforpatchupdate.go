package clusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterPropertiesForPatchUpdate struct {
	AiModelManagement *AiModelManagementStatus                      `json:"aiModelManagement,omitempty"`
	Backup            *BackupPropertiesForPatchUpdate               `json:"backup,omitempty"`
	DataEncryption    *HorizonDbClusterDataEncryptionForPatchUpdate `json:"dataEncryption,omitempty"`
	Network           *NetworkForPatchUpdate                        `json:"network,omitempty"`
}
