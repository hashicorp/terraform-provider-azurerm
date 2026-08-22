package machines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AgentConfiguration struct {
	ConfigMode                *AgentConfigurationMode   `json:"configMode,omitempty"`
	ExtensionsAllowList       *[]ConfigurationExtension `json:"extensionsAllowList,omitempty"`
	ExtensionsBlockList       *[]ConfigurationExtension `json:"extensionsBlockList,omitempty"`
	ExtensionsEnabled         *string                   `json:"extensionsEnabled,omitempty"`
	GuestConfigurationEnabled *string                   `json:"guestConfigurationEnabled,omitempty"`
	IncomingConnectionsPorts  *[]string                 `json:"incomingConnectionsPorts,omitempty"`
	ProxyBypass               *[]string                 `json:"proxyBypass,omitempty"`
	ProxyURL                  *string                   `json:"proxyUrl,omitempty"`
}
