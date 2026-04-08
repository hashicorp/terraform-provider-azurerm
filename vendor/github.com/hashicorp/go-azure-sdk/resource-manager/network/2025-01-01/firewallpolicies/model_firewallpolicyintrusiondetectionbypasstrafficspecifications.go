package firewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FirewallPolicyIntrusionDetectionBypassTrafficSpecifications struct {
	Description          *string                                   `json:"description,omitempty"`
	DestinationAddresses *[]string                                 `json:"destinationAddresses,omitempty"`
	DestinationIPGroups  *[]string                                 `json:"destinationIpGroups,omitempty"`
	DestinationPorts     *[]string                                 `json:"destinationPorts,omitempty"`
	Name                 *string                                   `json:"name,omitempty"`
	Protocol             *FirewallPolicyIntrusionDetectionProtocol `json:"protocol,omitempty"`
	SourceAddresses      *[]string                                 `json:"sourceAddresses,omitempty"`
	SourceIPGroups       *[]string                                 `json:"sourceIpGroups,omitempty"`
}
