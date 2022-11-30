package exports

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExportDataset struct {
	Configuration *ExportDatasetConfiguration `json:"configuration,omitempty"`
	Granularity   *GranularityType            `json:"granularity,omitempty"`
}
