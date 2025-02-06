package sqlvirtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SQLTempDbSettings struct {
	DataFileCount     *int64   `json:"dataFileCount,omitempty"`
	DataFileSize      *int64   `json:"dataFileSize,omitempty"`
	DataGrowth        *int64   `json:"dataGrowth,omitempty"`
	DefaultFilePath   *string  `json:"defaultFilePath,omitempty"`
	LogFileSize       *int64   `json:"logFileSize,omitempty"`
	LogGrowth         *int64   `json:"logGrowth,omitempty"`
	Luns              *[]int64 `json:"luns,omitempty"`
	PersistFolder     *bool    `json:"persistFolder,omitempty"`
	PersistFolderPath *string  `json:"persistFolderPath,omitempty"`
	UseStoragePool    *bool    `json:"useStoragePool,omitempty"`
}
