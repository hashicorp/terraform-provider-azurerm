package nginxconfiguration

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NginxConfigurationProperties struct {
	Files             *[]NginxConfigurationFile  `json:"files,omitempty"`
	Package           *NginxConfigurationPackage `json:"package,omitempty"`
	ProtectedFiles    *[]NginxConfigurationFile  `json:"protectedFiles,omitempty"`
	ProvisioningState *ProvisioningState         `json:"provisioningState,omitempty"`
	RootFile          *string                    `json:"rootFile,omitempty"`
}
