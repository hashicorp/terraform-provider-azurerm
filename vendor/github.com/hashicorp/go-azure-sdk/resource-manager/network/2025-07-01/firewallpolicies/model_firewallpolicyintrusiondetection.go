package firewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FirewallPolicyIntrusionDetection struct {
	Configuration *FirewallPolicyIntrusionDetectionConfiguration `json:"configuration,omitempty"`
	Mode          *FirewallPolicyIntrusionDetectionStateType     `json:"mode,omitempty"`
	Profile       *FirewallPolicyIntrusionDetectionProfileType   `json:"profile,omitempty"`
}
