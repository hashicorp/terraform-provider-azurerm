package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupFileInfo struct {
	FamilySequenceNumber *int64            `json:"familySequenceNumber,omitempty"`
	FileLocation         *string           `json:"fileLocation,omitempty"`
	Status               *BackupFileStatus `json:"status,omitempty"`
}
