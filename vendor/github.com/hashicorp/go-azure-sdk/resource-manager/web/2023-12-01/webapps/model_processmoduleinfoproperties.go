package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProcessModuleInfoProperties struct {
	BaseAddress      *string `json:"base_address,omitempty"`
	FileDescription  *string `json:"file_description,omitempty"`
	FileName         *string `json:"file_name,omitempty"`
	FilePath         *string `json:"file_path,omitempty"`
	FileVersion      *string `json:"file_version,omitempty"`
	Href             *string `json:"href,omitempty"`
	IsDebug          *bool   `json:"is_debug,omitempty"`
	Language         *string `json:"language,omitempty"`
	ModuleMemorySize *int64  `json:"module_memory_size,omitempty"`
	Product          *string `json:"product,omitempty"`
	ProductVersion   *string `json:"product_version,omitempty"`
}
