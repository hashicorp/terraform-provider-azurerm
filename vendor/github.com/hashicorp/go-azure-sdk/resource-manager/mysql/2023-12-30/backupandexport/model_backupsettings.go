package backupandexport

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupSettings struct {
	BackupFormat *BackupFormat `json:"backupFormat,omitempty"`
	BackupName   string        `json:"backupName"`
}
