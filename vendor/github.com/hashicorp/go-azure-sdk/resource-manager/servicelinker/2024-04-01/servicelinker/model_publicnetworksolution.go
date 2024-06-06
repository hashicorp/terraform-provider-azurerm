package servicelinker

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PublicNetworkSolution struct {
	Action                 *ActionType             `json:"action,omitempty"`
	DeleteOrUpdateBehavior *DeleteOrUpdateBehavior `json:"deleteOrUpdateBehavior,omitempty"`
	FirewallRules          *FirewallRules          `json:"firewallRules,omitempty"`
}
