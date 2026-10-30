package afdorigingroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AFDOriginGroupProperties struct {
	Authentication                                        *OriginAuthenticationProperties  `json:"authentication,omitempty"`
	DeploymentStatus                                      *DeploymentStatus                `json:"deploymentStatus,omitempty"`
	HealthProbeSettings                                   *HealthProbeParameters           `json:"healthProbeSettings,omitempty"`
	LoadBalancingSettings                                 *LoadBalancingSettingsParameters `json:"loadBalancingSettings,omitempty"`
	ProfileName                                           *string                          `json:"profileName,omitempty"`
	ProvisioningState                                     *AfdProvisioningState            `json:"provisioningState,omitempty"`
	SessionAffinityState                                  *EnabledState                    `json:"sessionAffinityState,omitempty"`
	TrafficRestorationTimeToHealedOrNewEndpointsInMinutes *int64                           `json:"trafficRestorationTimeToHealedOrNewEndpointsInMinutes,omitempty"`
}
