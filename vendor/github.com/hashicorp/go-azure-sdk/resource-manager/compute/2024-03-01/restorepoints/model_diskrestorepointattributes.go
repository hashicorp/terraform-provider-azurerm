package restorepoints

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiskRestorePointAttributes struct {
	Encryption             *RestorePointEncryption `json:"encryption,omitempty"`
	Id                     *string                 `json:"id,omitempty"`
	SourceDiskRestorePoint *ApiEntityReference     `json:"sourceDiskRestorePoint,omitempty"`
}
