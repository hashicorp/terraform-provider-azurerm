package replications

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReplicationProperties struct {
	ProvisioningState     *ProvisioningState `json:"provisioningState,omitempty"`
	RegionEndpointEnabled *bool              `json:"regionEndpointEnabled,omitempty"`
	Status                *Status            `json:"status,omitempty"`
	ZoneRedundancy        *ZoneRedundancy    `json:"zoneRedundancy,omitempty"`
}
