package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedClusterNATGatewayProfile struct {
	EffectiveOutboundIPs     *[]ResourceReference                    `json:"effectiveOutboundIPs,omitempty"`
	IdleTimeoutInMinutes     *int64                                  `json:"idleTimeoutInMinutes,omitempty"`
	ManagedOutboundIPProfile *ManagedClusterManagedOutboundIPProfile `json:"managedOutboundIPProfile,omitempty"`
}
