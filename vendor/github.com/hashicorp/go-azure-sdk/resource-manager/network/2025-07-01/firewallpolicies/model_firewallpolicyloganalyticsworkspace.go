package firewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FirewallPolicyLogAnalyticsWorkspace struct {
	Region      *string      `json:"region,omitempty"`
	WorkspaceId *SubResource `json:"workspaceId,omitempty"`
}
