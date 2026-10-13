package firewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FirewallPolicyIntrusionDetectionSignatureSpecification struct {
	Id   *string                                    `json:"id,omitempty"`
	Mode *FirewallPolicyIntrusionDetectionStateType `json:"mode,omitempty"`
}
