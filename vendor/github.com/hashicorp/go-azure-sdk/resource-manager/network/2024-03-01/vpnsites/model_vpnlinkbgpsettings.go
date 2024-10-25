package vpnsites

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VpnLinkBgpSettings struct {
	Asn               *int64  `json:"asn,omitempty"`
	BgpPeeringAddress *string `json:"bgpPeeringAddress,omitempty"`
}
