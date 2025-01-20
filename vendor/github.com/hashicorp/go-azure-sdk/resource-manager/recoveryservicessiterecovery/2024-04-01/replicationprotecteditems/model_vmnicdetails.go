package replicationprotecteditems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VMNicDetails struct {
	EnableAcceleratedNetworkingOnRecovery *bool              `json:"enableAcceleratedNetworkingOnRecovery,omitempty"`
	EnableAcceleratedNetworkingOnTfo      *bool              `json:"enableAcceleratedNetworkingOnTfo,omitempty"`
	IPConfigs                             *[]IPConfigDetails `json:"ipConfigs,omitempty"`
	NicId                                 *string            `json:"nicId,omitempty"`
	RecoveryNetworkSecurityGroupId        *string            `json:"recoveryNetworkSecurityGroupId,omitempty"`
	RecoveryNicName                       *string            `json:"recoveryNicName,omitempty"`
	RecoveryNicResourceGroupName          *string            `json:"recoveryNicResourceGroupName,omitempty"`
	RecoveryVMNetworkId                   *string            `json:"recoveryVMNetworkId,omitempty"`
	ReplicaNicId                          *string            `json:"replicaNicId,omitempty"`
	ReuseExistingNic                      *bool              `json:"reuseExistingNic,omitempty"`
	SelectionType                         *string            `json:"selectionType,omitempty"`
	SourceNicArmId                        *string            `json:"sourceNicArmId,omitempty"`
	TargetNicName                         *string            `json:"targetNicName,omitempty"`
	TfoNetworkSecurityGroupId             *string            `json:"tfoNetworkSecurityGroupId,omitempty"`
	TfoRecoveryNicName                    *string            `json:"tfoRecoveryNicName,omitempty"`
	TfoRecoveryNicResourceGroupName       *string            `json:"tfoRecoveryNicResourceGroupName,omitempty"`
	TfoReuseExistingNic                   *bool              `json:"tfoReuseExistingNic,omitempty"`
	TfoVMNetworkId                        *string            `json:"tfoVMNetworkId,omitempty"`
	VMNetworkName                         *string            `json:"vMNetworkName,omitempty"`
}
