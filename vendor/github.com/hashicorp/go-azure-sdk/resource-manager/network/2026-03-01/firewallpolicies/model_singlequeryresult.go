package firewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SingleQueryResult struct {
	Description               *string                               `json:"description,omitempty"`
	DestinationPorts          *[]string                             `json:"destinationPorts,omitempty"`
	Direction                 *FirewallPolicyIDPSSignatureDirection `json:"direction,omitempty"`
	Group                     *string                               `json:"group,omitempty"`
	InheritedFromParentPolicy *bool                                 `json:"inheritedFromParentPolicy,omitempty"`
	LastUpdated               *string                               `json:"lastUpdated,omitempty"`
	Mode                      *FirewallPolicyIDPSSignatureMode      `json:"mode,omitempty"`
	Protocol                  *string                               `json:"protocol,omitempty"`
	Severity                  *FirewallPolicyIDPSSignatureSeverity  `json:"severity,omitempty"`
	SignatureId               *int64                                `json:"signatureId,omitempty"`
	SourcePorts               *[]string                             `json:"sourcePorts,omitempty"`
}
