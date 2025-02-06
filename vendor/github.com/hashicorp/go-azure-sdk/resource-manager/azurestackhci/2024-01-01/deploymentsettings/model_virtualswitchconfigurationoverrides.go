package deploymentsettings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualSwitchConfigurationOverrides struct {
	EnableIov              *string `json:"enableIov,omitempty"`
	LoadBalancingAlgorithm *string `json:"loadBalancingAlgorithm,omitempty"`
}
