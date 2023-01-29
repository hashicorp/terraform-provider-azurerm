package lab

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LabNetworkProfile struct {
	LoadBalancerId *string `json:"loadBalancerId,omitempty"`
	PublicIPId     *string `json:"publicIpId,omitempty"`
	SubnetId       *string `json:"subnetId,omitempty"`
}
