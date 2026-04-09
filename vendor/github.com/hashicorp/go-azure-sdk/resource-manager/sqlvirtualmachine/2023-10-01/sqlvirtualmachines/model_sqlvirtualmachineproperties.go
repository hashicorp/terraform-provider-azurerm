package sqlvirtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SqlVirtualMachineProperties struct {
	AdditionalVMPatch                      *AdditionalOsPatch                      `json:"additionalVmPatch,omitempty"`
	AssessmentSettings                     *AssessmentSettings                     `json:"assessmentSettings,omitempty"`
	AutoBackupSettings                     *AutoBackupSettings                     `json:"autoBackupSettings,omitempty"`
	AutoPatchingSettings                   *AutoPatchingSettings                   `json:"autoPatchingSettings,omitempty"`
	EnableAutomaticUpgrade                 *bool                                   `json:"enableAutomaticUpgrade,omitempty"`
	KeyVaultCredentialSettings             *KeyVaultCredentialSettings             `json:"keyVaultCredentialSettings,omitempty"`
	LeastPrivilegeMode                     *LeastPrivilegeMode                     `json:"leastPrivilegeMode,omitempty"`
	OsType                                 *OsType                                 `json:"osType,omitempty"`
	ProvisioningState                      *string                                 `json:"provisioningState,omitempty"`
	ServerConfigurationsManagementSettings *ServerConfigurationsManagementSettings `json:"serverConfigurationsManagementSettings,omitempty"`
	SqlImageOffer                          *string                                 `json:"sqlImageOffer,omitempty"`
	SqlImageSku                            *SqlImageSku                            `json:"sqlImageSku,omitempty"`
	SqlManagement                          *SqlManagementMode                      `json:"sqlManagement,omitempty"`
	SqlServerLicenseType                   *SqlServerLicenseType                   `json:"sqlServerLicenseType,omitempty"`
	SqlVirtualMachineGroupResourceId       *string                                 `json:"sqlVirtualMachineGroupResourceId,omitempty"`
	StorageConfigurationSettings           *StorageConfigurationSettings           `json:"storageConfigurationSettings,omitempty"`
	TroubleshootingStatus                  *TroubleshootingStatus                  `json:"troubleshootingStatus,omitempty"`
	VirtualMachineIdentitySettings         *VirtualMachineIdentity                 `json:"virtualMachineIdentitySettings,omitempty"`
	VirtualMachineResourceId               *string                                 `json:"virtualMachineResourceId,omitempty"`
	WsfcDomainCredentials                  *WsfcDomainCredentials                  `json:"wsfcDomainCredentials,omitempty"`
	WsfcStaticIP                           *string                                 `json:"wsfcStaticIp,omitempty"`
}
