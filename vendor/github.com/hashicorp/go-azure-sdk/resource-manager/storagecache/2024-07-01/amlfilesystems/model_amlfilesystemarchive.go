package amlfilesystems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AmlFilesystemArchive struct {
	FilesystemPath *string                     `json:"filesystemPath,omitempty"`
	Status         *AmlFilesystemArchiveStatus `json:"status,omitempty"`
}
