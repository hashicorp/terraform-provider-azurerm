package managedclustersnapshots

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkProfileForSnapshot struct {
	LoadBalancerSku   *LoadBalancerSku   `json:"loadBalancerSku,omitempty"`
	NetworkMode       *NetworkMode       `json:"networkMode,omitempty"`
	NetworkPlugin     *NetworkPlugin     `json:"networkPlugin,omitempty"`
	NetworkPluginMode *NetworkPluginMode `json:"networkPluginMode,omitempty"`
	NetworkPolicy     *NetworkPolicy     `json:"networkPolicy,omitempty"`
}
