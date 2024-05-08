package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TrinoTelemetryConfig struct {
	HivecatalogName          *string `json:"hivecatalogName,omitempty"`
	HivecatalogSchema        *string `json:"hivecatalogSchema,omitempty"`
	PartitionRetentionInDays *int64  `json:"partitionRetentionInDays,omitempty"`
	Path                     *string `json:"path,omitempty"`
}
