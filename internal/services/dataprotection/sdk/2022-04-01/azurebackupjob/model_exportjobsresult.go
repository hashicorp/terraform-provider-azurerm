package azurebackupjob

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExportJobsResult struct {
	BlobSasKey          *string `json:"blobSasKey,omitempty"`
	BlobUrl             *string `json:"blobUrl,omitempty"`
	ExcelFileBlobSasKey *string `json:"excelFileBlobSasKey,omitempty"`
	ExcelFileBlobUrl    *string `json:"excelFileBlobUrl,omitempty"`
}
