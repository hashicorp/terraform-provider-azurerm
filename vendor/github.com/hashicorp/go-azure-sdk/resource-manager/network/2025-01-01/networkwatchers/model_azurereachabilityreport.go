package networkwatchers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureReachabilityReport struct {
	AggregationLevel   string                          `json:"aggregationLevel"`
	ProviderLocation   AzureReachabilityReportLocation `json:"providerLocation"`
	ReachabilityReport []AzureReachabilityReportItem   `json:"reachabilityReport"`
}
