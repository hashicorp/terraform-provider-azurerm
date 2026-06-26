package networkwatchers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureReachabilityReportItem struct {
	AzureLocation *string                               `json:"azureLocation,omitempty"`
	Latencies     *[]AzureReachabilityReportLatencyInfo `json:"latencies,omitempty"`
	Provider      *string                               `json:"provider,omitempty"`
}
