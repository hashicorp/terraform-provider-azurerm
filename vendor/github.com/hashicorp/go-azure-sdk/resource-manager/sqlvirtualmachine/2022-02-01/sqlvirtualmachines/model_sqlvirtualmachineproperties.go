package sqlvirtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SqlVirtualMachineProperties struct {
	AssessmentSettings                     *AssessmentSettings                     `json:"assessmentSettings,omitempty"`
	AutoBackupSettings                     *AutoBackupSettings                     `json:"autoBackupSettings,omitempty"`
	AutoPatchingSettings                   *AutoPatchingSettings                   `json:"autoPatchingSettings,omitempty"`
	KeyVaultCredentialSettings             *KeyVaultCredentialSettings             `json:"keyVaultCredentialSettings,omitempty"`
	ProvisioningState                      *string                                 `json:"provisioningState,omitempty"`
	ServerConfigurationsManagementSettings *ServerConfigurationsManagementSettings `json:"serverConfigurationsManagementSettings,omitempty"`
	SqlImageOffer                          *string                                 `json:"sqlImageOffer,omitempty"`
	SqlImageSku                            *SqlImageSku                            `json:"sqlImageSku,omitempty"`
	SqlManagement                          *SqlManagementMode                      `json:"sqlManagement,omitempty"`
	SqlServerLicenseType                   *SqlServerLicenseType                   `json:"sqlServerLicenseType,omitempty"`
	SqlVirtualMachineGroupResourceId       *string                                 `json:"sqlVirtualMachineGroupResourceId,omitempty"`
	StorageConfigurationSettings           *StorageConfigurationSettings           `json:"storageConfigurationSettings,omitempty"`
	VirtualMachineResourceId               *string                                 `json:"virtualMachineResourceId,omitempty"`
	WsfcDomainCredentials                  *WsfcDomainCredentials                  `json:"wsfcDomainCredentials,omitempty"`
	WsfcStaticIP                           *string                                 `json:"wsfcStaticIp,omitempty"`
}
