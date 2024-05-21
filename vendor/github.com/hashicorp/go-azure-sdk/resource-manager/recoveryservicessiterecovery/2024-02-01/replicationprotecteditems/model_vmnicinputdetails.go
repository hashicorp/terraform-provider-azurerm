package replicationprotecteditems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VMNicInputDetails struct {
	EnableAcceleratedNetworkingOnRecovery *bool                   `json:"enableAcceleratedNetworkingOnRecovery,omitempty"`
	EnableAcceleratedNetworkingOnTfo      *bool                   `json:"enableAcceleratedNetworkingOnTfo,omitempty"`
	IPConfigs                             *[]IPConfigInputDetails `json:"ipConfigs,omitempty"`
	NicId                                 *string                 `json:"nicId,omitempty"`
	RecoveryNetworkSecurityGroupId        *string                 `json:"recoveryNetworkSecurityGroupId,omitempty"`
	RecoveryNicName                       *string                 `json:"recoveryNicName,omitempty"`
	RecoveryNicResourceGroupName          *string                 `json:"recoveryNicResourceGroupName,omitempty"`
	ReuseExistingNic                      *bool                   `json:"reuseExistingNic,omitempty"`
	SelectionType                         *string                 `json:"selectionType,omitempty"`
	TargetNicName                         *string                 `json:"targetNicName,omitempty"`
	TfoNetworkSecurityGroupId             *string                 `json:"tfoNetworkSecurityGroupId,omitempty"`
	TfoNicName                            *string                 `json:"tfoNicName,omitempty"`
	TfoNicResourceGroupName               *string                 `json:"tfoNicResourceGroupName,omitempty"`
	TfoReuseExistingNic                   *bool                   `json:"tfoReuseExistingNic,omitempty"`
}
