package agents

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SandboxEgressConfiguration struct {
	AllowHTTPMcpServerNetworkAccess *bool                          `json:"allowHttpMcpServerNetworkAccess,omitempty"`
	AllowedCodeRepositories         *[]string                      `json:"allowedCodeRepositories,omitempty"`
	AllowedHosts                    *[]string                      `json:"allowedHosts,omitempty"`
	AllowedRegistries               *[]string                      `json:"allowedRegistries,omitempty"`
	HTTPProxy                       *SandboxHTTPProxyConfiguration `json:"httpProxy,omitempty"`
	Mode                            *SandboxEgressMode             `json:"mode,omitempty"`
}
