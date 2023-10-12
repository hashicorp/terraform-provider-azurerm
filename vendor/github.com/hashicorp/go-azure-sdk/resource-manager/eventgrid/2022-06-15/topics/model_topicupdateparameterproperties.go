package topics

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TopicUpdateParameterProperties struct {
	DataResidencyBoundary *DataResidencyBoundary `json:"dataResidencyBoundary,omitempty"`
	DisableLocalAuth      *bool                  `json:"disableLocalAuth,omitempty"`
	InboundIPRules        *[]InboundIPRule       `json:"inboundIpRules,omitempty"`
	PublicNetworkAccess   *PublicNetworkAccess   `json:"publicNetworkAccess,omitempty"`
}
