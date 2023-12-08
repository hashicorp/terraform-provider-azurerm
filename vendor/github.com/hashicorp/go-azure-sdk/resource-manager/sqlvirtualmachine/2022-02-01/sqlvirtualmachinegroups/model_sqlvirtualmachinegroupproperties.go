package sqlvirtualmachinegroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SqlVirtualMachineGroupProperties struct {
	ClusterConfiguration *ClusterConfiguration `json:"clusterConfiguration,omitempty"`
	ClusterManagerType   *ClusterManagerType   `json:"clusterManagerType,omitempty"`
	ProvisioningState    *string               `json:"provisioningState,omitempty"`
	ScaleType            *ScaleType            `json:"scaleType,omitempty"`
	SqlImageOffer        *string               `json:"sqlImageOffer,omitempty"`
	SqlImageSku          *SqlVMGroupImageSku   `json:"sqlImageSku,omitempty"`
	WsfcDomainProfile    *WsfcDomainProfile    `json:"wsfcDomainProfile,omitempty"`
}
