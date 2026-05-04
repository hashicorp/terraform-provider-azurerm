package endpoints

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeepCreatedOriginGroupProperties struct {
	HealthProbeSettings                                   *HealthProbeParameters                       `json:"healthProbeSettings,omitempty"`
	Origins                                               []ResourceReference                          `json:"origins"`
	ResponseBasedOriginErrorDetectionSettings             *ResponseBasedOriginErrorDetectionParameters `json:"responseBasedOriginErrorDetectionSettings,omitempty"`
	TrafficRestorationTimeToHealedOrNewEndpointsInMinutes *int64                                       `json:"trafficRestorationTimeToHealedOrNewEndpointsInMinutes,omitempty"`
}
