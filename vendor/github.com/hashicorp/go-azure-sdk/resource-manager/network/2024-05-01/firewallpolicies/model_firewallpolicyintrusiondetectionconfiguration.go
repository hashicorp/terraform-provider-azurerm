package firewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FirewallPolicyIntrusionDetectionConfiguration struct {
	BypassTrafficSettings *[]FirewallPolicyIntrusionDetectionBypassTrafficSpecifications `json:"bypassTrafficSettings,omitempty"`
	PrivateRanges         *[]string                                                      `json:"privateRanges,omitempty"`
	SignatureOverrides    *[]FirewallPolicyIntrusionDetectionSignatureSpecification      `json:"signatureOverrides,omitempty"`
}
