package servicelinker

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FirewallRules struct {
	AzureServices  *AllowType `json:"azureServices,omitempty"`
	CallerClientIP *AllowType `json:"callerClientIP,omitempty"`
	IPRanges       *[]string  `json:"ipRanges,omitempty"`
}
