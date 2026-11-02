package defenderforstorage

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BlobsScanSummary struct {
	FailedBlobsCount    *int64   `json:"failedBlobsCount,omitempty"`
	MaliciousBlobsCount *int64   `json:"maliciousBlobsCount,omitempty"`
	ScannedBlobsInGB    *float64 `json:"scannedBlobsInGB,omitempty"`
	SkippedBlobsCount   *int64   `json:"skippedBlobsCount,omitempty"`
	TotalBlobsScanned   *int64   `json:"totalBlobsScanned,omitempty"`
}
