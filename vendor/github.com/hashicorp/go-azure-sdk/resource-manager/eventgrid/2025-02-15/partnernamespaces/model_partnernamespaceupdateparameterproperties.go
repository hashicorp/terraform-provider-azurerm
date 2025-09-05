package partnernamespaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PartnerNamespaceUpdateParameterProperties struct {
	DisableLocalAuth         *bool                `json:"disableLocalAuth,omitempty"`
	InboundIPRules           *[]InboundIPRule     `json:"inboundIpRules,omitempty"`
	MinimumTlsVersionAllowed *TlsVersion          `json:"minimumTlsVersionAllowed,omitempty"`
	PublicNetworkAccess      *PublicNetworkAccess `json:"publicNetworkAccess,omitempty"`
}
