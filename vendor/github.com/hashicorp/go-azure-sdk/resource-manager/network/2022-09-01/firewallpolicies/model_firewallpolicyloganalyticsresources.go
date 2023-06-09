package firewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FirewallPolicyLogAnalyticsResources struct {
	DefaultWorkspaceId *SubResource                           `json:"defaultWorkspaceId,omitempty"`
	Workspaces         *[]FirewallPolicyLogAnalyticsWorkspace `json:"workspaces,omitempty"`
}
