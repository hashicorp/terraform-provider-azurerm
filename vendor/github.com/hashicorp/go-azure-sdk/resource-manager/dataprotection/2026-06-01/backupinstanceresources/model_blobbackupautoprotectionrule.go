package backupinstanceresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BlobBackupAutoProtectionRule struct {
	Mode       BlobBackupRuleMode    `json:"mode"`
	ObjectType string                `json:"objectType"`
	Pattern    string                `json:"pattern"`
	Type       BlobBackupPatternType `json:"type"`
}
