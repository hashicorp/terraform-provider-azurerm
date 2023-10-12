package virtualwans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetVpnSitesConfigurationRequest struct {
	OutputBlobSasUrl string    `json:"outputBlobSasUrl"`
	VpnSites         *[]string `json:"vpnSites,omitempty"`
}
