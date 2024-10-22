package nginxdeployment

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NginxDeploymentUpdateProperties struct {
	AutoUpgradeProfile       *AutoUpgradeProfile                             `json:"autoUpgradeProfile,omitempty"`
	EnableDiagnosticsSupport *bool                                           `json:"enableDiagnosticsSupport,omitempty"`
	Logging                  *NginxLogging                                   `json:"logging,omitempty"`
	NginxAppProtect          *NginxDeploymentUpdatePropertiesNginxAppProtect `json:"nginxAppProtect,omitempty"`
	ScalingProperties        *NginxDeploymentScalingProperties               `json:"scalingProperties,omitempty"`
	UserProfile              *NginxDeploymentUserProfile                     `json:"userProfile,omitempty"`
}
