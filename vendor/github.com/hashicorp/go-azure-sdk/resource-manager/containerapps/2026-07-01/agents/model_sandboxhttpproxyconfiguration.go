package agents

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SandboxHTTPProxyConfiguration struct {
	HTTPSURL *string   `json:"httpsUrl,omitempty"`
	HTTPURL  *string   `json:"httpUrl,omitempty"`
	NoProxy  *[]string `json:"noProxy,omitempty"`
}
