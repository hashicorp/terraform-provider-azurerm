package deployments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentProperties struct {
	CallRateLimit            *CallRateLimit                       `json:"callRateLimit,omitempty"`
	Capabilities             *map[string]string                   `json:"capabilities,omitempty"`
	CapacitySettings         *DeploymentCapacitySettings          `json:"capacitySettings,omitempty"`
	CurrentCapacity          *int64                               `json:"currentCapacity,omitempty"`
	DynamicThrottlingEnabled *bool                                `json:"dynamicThrottlingEnabled,omitempty"`
	Model                    *DeploymentModel                     `json:"model,omitempty"`
	ParentDeploymentName     *string                              `json:"parentDeploymentName,omitempty"`
	ProvisioningState        *DeploymentProvisioningState         `json:"provisioningState,omitempty"`
	RaiPolicyName            *string                              `json:"raiPolicyName,omitempty"`
	RateLimits               *[]ThrottlingRule                    `json:"rateLimits,omitempty"`
	ScaleSettings            *DeploymentScaleSettings             `json:"scaleSettings,omitempty"`
	VersionUpgradeOption     *DeploymentModelVersionUpgradeOption `json:"versionUpgradeOption,omitempty"`
}
