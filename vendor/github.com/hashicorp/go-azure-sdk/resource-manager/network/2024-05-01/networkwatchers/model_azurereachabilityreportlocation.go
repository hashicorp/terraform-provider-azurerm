package networkwatchers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureReachabilityReportLocation struct {
	City    *string `json:"city,omitempty"`
	Country string  `json:"country"`
	State   *string `json:"state,omitempty"`
}
