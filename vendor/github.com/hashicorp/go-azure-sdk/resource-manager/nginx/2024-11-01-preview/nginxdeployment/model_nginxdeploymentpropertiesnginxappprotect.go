package nginxdeployment

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NginxDeploymentPropertiesNginxAppProtect struct {
	WebApplicationFirewallSettings WebApplicationFirewallSettings `json:"webApplicationFirewallSettings"`
	WebApplicationFirewallStatus   *WebApplicationFirewallStatus  `json:"webApplicationFirewallStatus,omitempty"`
}
