package provisionedclusterinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AgentPoolProperties struct {
	Count             *int64                             `json:"count,omitempty"`
	EnableAutoScaling *bool                              `json:"enableAutoScaling,omitempty"`
	KubernetesVersion *string                            `json:"kubernetesVersion,omitempty"`
	MaxCount          *int64                             `json:"maxCount,omitempty"`
	MaxPods           *int64                             `json:"maxPods,omitempty"`
	MinCount          *int64                             `json:"minCount,omitempty"`
	NodeLabels        *map[string]string                 `json:"nodeLabels,omitempty"`
	NodeTaints        *[]string                          `json:"nodeTaints,omitempty"`
	OsSKU             *OSSKU                             `json:"osSKU,omitempty"`
	OsType            *OsType                            `json:"osType,omitempty"`
	ProvisioningState *ResourceProvisioningState         `json:"provisioningState,omitempty"`
	Status            *AgentPoolProvisioningStatusStatus `json:"status,omitempty"`
	VMSize            *string                            `json:"vmSize,omitempty"`
}
