package nginxdeployment

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NginxDeploymentProperties struct {
	AutoUpgradeProfile       *AutoUpgradeProfile                       `json:"autoUpgradeProfile,omitempty"`
	EnableDiagnosticsSupport *bool                                     `json:"enableDiagnosticsSupport,omitempty"`
	IPAddress                *string                                   `json:"ipAddress,omitempty"`
	Logging                  *NginxLogging                             `json:"logging,omitempty"`
	ManagedResourceGroup     *string                                   `json:"managedResourceGroup,omitempty"`
	NetworkProfile           *NginxNetworkProfile                      `json:"networkProfile,omitempty"`
	NginxAppProtect          *NginxDeploymentPropertiesNginxAppProtect `json:"nginxAppProtect,omitempty"`
	NginxVersion             *string                                   `json:"nginxVersion,omitempty"`
	ProvisioningState        *ProvisioningState                        `json:"provisioningState,omitempty"`
	ScalingProperties        *NginxDeploymentScalingProperties         `json:"scalingProperties,omitempty"`
	UserProfile              *NginxDeploymentUserProfile               `json:"userProfile,omitempty"`
}
