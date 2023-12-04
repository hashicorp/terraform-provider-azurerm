package nginxdeployment

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NginxPrivateIPAddress struct {
	PrivateIPAddress          *string                         `json:"privateIPAddress,omitempty"`
	PrivateIPAllocationMethod *NginxPrivateIPAllocationMethod `json:"privateIPAllocationMethod,omitempty"`
	SubnetId                  *string                         `json:"subnetId,omitempty"`
}
