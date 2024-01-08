package firewalls

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkProfile struct {
	EgressNatIP       *[]IPAddress       `json:"egressNatIp,omitempty"`
	EnableEgressNat   EgressNat          `json:"enableEgressNat"`
	NetworkType       NetworkType        `json:"networkType"`
	PublicIPs         []IPAddress        `json:"publicIps"`
	TrustedRanges     *[]string          `json:"trustedRanges,omitempty"`
	VnetConfiguration *VnetConfiguration `json:"vnetConfiguration,omitempty"`
	VwanConfiguration *VwanConfiguration `json:"vwanConfiguration,omitempty"`
}
