package virtualnetworkgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VpnClientIPsecParameters struct {
	DhGroup             DhGroup         `json:"dhGroup"`
	IPsecEncryption     IPsecEncryption `json:"ipsecEncryption"`
	IPsecIntegrity      IPsecIntegrity  `json:"ipsecIntegrity"`
	IkeEncryption       IkeEncryption   `json:"ikeEncryption"`
	IkeIntegrity        IkeIntegrity    `json:"ikeIntegrity"`
	PfsGroup            PfsGroup        `json:"pfsGroup"`
	SaDataSizeKilobytes int64           `json:"saDataSizeKilobytes"`
	SaLifeTimeSeconds   int64           `json:"saLifeTimeSeconds"`
}
