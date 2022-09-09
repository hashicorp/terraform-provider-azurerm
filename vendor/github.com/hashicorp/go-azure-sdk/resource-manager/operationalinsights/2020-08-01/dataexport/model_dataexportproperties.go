package dataexport

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataExportProperties struct {
	CreatedDate      *string      `json:"createdDate,omitempty"`
	DataExportId     *string      `json:"dataExportId,omitempty"`
	Destination      *Destination `json:"destination,omitempty"`
	Enable           *bool        `json:"enable,omitempty"`
	LastModifiedDate *string      `json:"lastModifiedDate,omitempty"`
	TableNames       []string     `json:"tableNames"`
}
