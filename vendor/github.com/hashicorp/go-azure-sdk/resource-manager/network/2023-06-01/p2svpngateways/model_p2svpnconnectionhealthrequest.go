package p2svpngateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type P2SVpnConnectionHealthRequest struct {
	OutputBlobSasUrl   *string   `json:"outputBlobSasUrl,omitempty"`
	VpnUserNamesFilter *[]string `json:"vpnUserNamesFilter,omitempty"`
}
