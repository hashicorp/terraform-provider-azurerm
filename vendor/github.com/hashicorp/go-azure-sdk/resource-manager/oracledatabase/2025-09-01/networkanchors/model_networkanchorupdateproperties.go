package networkanchors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkAnchorUpdateProperties struct {
	IsOracleDnsForwardingEndpointEnabled *bool   `json:"isOracleDnsForwardingEndpointEnabled,omitempty"`
	IsOracleDnsListeningEndpointEnabled  *bool   `json:"isOracleDnsListeningEndpointEnabled,omitempty"`
	IsOracleToAzureDnsZoneSyncEnabled    *bool   `json:"isOracleToAzureDnsZoneSyncEnabled,omitempty"`
	OciBackupCidrBlock                   *string `json:"ociBackupCidrBlock,omitempty"`
}
