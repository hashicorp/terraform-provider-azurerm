package backupvaultresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExportJobsResult struct {
	BlobSasKey          *string `json:"blobSasKey,omitempty"`
	BlobURL             *string `json:"blobUrl,omitempty"`
	ExcelFileBlobSasKey *string `json:"excelFileBlobSasKey,omitempty"`
	ExcelFileBlobURL    *string `json:"excelFileBlobUrl,omitempty"`
}
