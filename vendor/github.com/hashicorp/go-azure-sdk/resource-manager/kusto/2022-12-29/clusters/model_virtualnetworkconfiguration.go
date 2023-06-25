package clusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualNetworkConfiguration struct {
	DataManagementPublicIPId string `json:"dataManagementPublicIpId"`
	EnginePublicIPId         string `json:"enginePublicIpId"`
	SubnetId                 string `json:"subnetId"`
}
