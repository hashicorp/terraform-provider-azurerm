package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterPoolNetworkProfile struct {
	ApiServerAuthorizedIPRanges *[]string     `json:"apiServerAuthorizedIpRanges,omitempty"`
	EnablePrivateApiServer      *bool         `json:"enablePrivateApiServer,omitempty"`
	OutboundType                *OutboundType `json:"outboundType,omitempty"`
	SubnetId                    string        `json:"subnetId"`
}
