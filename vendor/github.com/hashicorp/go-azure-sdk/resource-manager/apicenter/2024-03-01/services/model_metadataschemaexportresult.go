package services

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MetadataSchemaExportResult struct {
	Format *MetadataSchemaExportFormat `json:"format,omitempty"`
	Value  *string                     `json:"value,omitempty"`
}
