package backups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupRestoreFiles struct {
	DestinationVolumeId string   `json:"destinationVolumeId"`
	FileList            []string `json:"fileList"`
	RestoreFilePath     *string  `json:"restoreFilePath,omitempty"`
}
