package tables

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResultStatistics struct {
	IngestedRecords *int64   `json:"ingestedRecords,omitempty"`
	Progress        *float64 `json:"progress,omitempty"`
	ScannedGb       *float64 `json:"scannedGb,omitempty"`
}
