package networkinterfaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProbePropertiesFormat struct {
	IntervalInSeconds  *int64             `json:"intervalInSeconds,omitempty"`
	LoadBalancingRules *[]SubResource     `json:"loadBalancingRules,omitempty"`
	NumberOfProbes     *int64             `json:"numberOfProbes,omitempty"`
	Port               int64              `json:"port"`
	ProbeThreshold     *int64             `json:"probeThreshold,omitempty"`
	Protocol           ProbeProtocol      `json:"protocol"`
	ProvisioningState  *ProvisioningState `json:"provisioningState,omitempty"`
	RequestPath        *string            `json:"requestPath,omitempty"`
}
