package hcpopenshiftclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterAutoscalingProfile struct {
	MaxNodeProvisionTimeSeconds *int64 `json:"maxNodeProvisionTimeSeconds,omitempty"`
	MaxNodesTotal               *int64 `json:"maxNodesTotal,omitempty"`
	MaxPodGracePeriodSeconds    *int64 `json:"maxPodGracePeriodSeconds,omitempty"`
	PodPriorityThreshold        *int64 `json:"podPriorityThreshold,omitempty"`
}
