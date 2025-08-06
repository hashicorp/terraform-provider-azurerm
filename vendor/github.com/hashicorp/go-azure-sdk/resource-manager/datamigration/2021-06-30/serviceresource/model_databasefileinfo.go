package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatabaseFileInfo struct {
	DatabaseName     *string           `json:"databaseName,omitempty"`
	FileType         *DatabaseFileType `json:"fileType,omitempty"`
	Id               *string           `json:"id,omitempty"`
	LogicalName      *string           `json:"logicalName,omitempty"`
	PhysicalFullName *string           `json:"physicalFullName,omitempty"`
	RestoreFullName  *string           `json:"restoreFullName,omitempty"`
	SizeMB           *float64          `json:"sizeMB,omitempty"`
}
