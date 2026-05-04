package nginxdeployment

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebApplicationFirewallStatus struct {
	AttackSignaturesPackage *WebApplicationFirewallPackage           `json:"attackSignaturesPackage,omitempty"`
	BotSignaturesPackage    *WebApplicationFirewallPackage           `json:"botSignaturesPackage,omitempty"`
	ComponentVersions       *WebApplicationFirewallComponentVersions `json:"componentVersions,omitempty"`
	ThreatCampaignsPackage  *WebApplicationFirewallPackage           `json:"threatCampaignsPackage,omitempty"`
}
