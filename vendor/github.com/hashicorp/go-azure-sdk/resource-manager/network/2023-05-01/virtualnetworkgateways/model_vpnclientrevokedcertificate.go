package virtualnetworkgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VpnClientRevokedCertificate struct {
	Etag       *string                                      `json:"etag,omitempty"`
	Id         *string                                      `json:"id,omitempty"`
	Name       *string                                      `json:"name,omitempty"`
	Properties *VpnClientRevokedCertificatePropertiesFormat `json:"properties,omitempty"`
}
