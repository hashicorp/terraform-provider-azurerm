package ddoscustompolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DdosDetectionRulePropertiesFormat struct {
	DetectionMode        *DdosDetectionMode    `json:"detectionMode,omitempty"`
	ProvisioningState    *ProvisioningState    `json:"provisioningState,omitempty"`
	TrafficDetectionRule *TrafficDetectionRule `json:"trafficDetectionRule,omitempty"`
}
