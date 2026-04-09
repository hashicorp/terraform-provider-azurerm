package firewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExplicitProxy struct {
	EnableExplicitProxy *bool   `json:"enableExplicitProxy,omitempty"`
	EnablePacFile       *bool   `json:"enablePacFile,omitempty"`
	HTTPPort            *int64  `json:"httpPort,omitempty"`
	HTTPSPort           *int64  `json:"httpsPort,omitempty"`
	PacFile             *string `json:"pacFile,omitempty"`
	PacFilePort         *int64  `json:"pacFilePort,omitempty"`
}
