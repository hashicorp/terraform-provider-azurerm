package vpnsites

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IPConfigurationBgpPeeringAddress struct {
	CustomBgpIPAddresses  *[]string `json:"customBgpIpAddresses,omitempty"`
	DefaultBgpIPAddresses *[]string `json:"defaultBgpIpAddresses,omitempty"`
	IPconfigurationId     *string   `json:"ipconfigurationId,omitempty"`
	TunnelIPAddresses     *[]string `json:"tunnelIpAddresses,omitempty"`
}
