package sqlvirtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SqlVirtualMachineProperties struct {
	AssessmentSettings                     *AssessmentSettings                     `json:"assessmentSettings"`
	AutoBackupSettings                     *AutoBackupSettings                     `json:"autoBackupSettings"`
	AutoPatchingSettings                   *AutoPatchingSettings                   `json:"autoPatchingSettings"`
	KeyVaultCredentialSettings             *KeyVaultCredentialSettings             `json:"keyVaultCredentialSettings"`
	ProvisioningState                      *string                                 `json:"provisioningState,omitempty"`
	ServerConfigurationsManagementSettings *ServerConfigurationsManagementSettings `json:"serverConfigurationsManagementSettings"`
	SqlImageOffer                          *string                                 `json:"sqlImageOffer,omitempty"`
	SqlImageSku                            *SqlImageSku                            `json:"sqlImageSku,omitempty"`
	SqlManagement                          *SqlManagementMode                      `json:"sqlManagement,omitempty"`
	SqlServerLicenseType                   *SqlServerLicenseType                   `json:"sqlServerLicenseType,omitempty"`
	SqlVirtualMachineGroupResourceId       *string                                 `json:"sqlVirtualMachineGroupResourceId,omitempty"`
	StorageConfigurationSettings           *StorageConfigurationSettings           `json:"storageConfigurationSettings"`
	VirtualMachineResourceId               *string                                 `json:"virtualMachineResourceId,omitempty"`
	WsfcDomainCredentials                  *WsfcDomainCredentials                  `json:"wsfcDomainCredentials"`
	WsfcStaticIP                           *string                                 `json:"wsfcStaticIp,omitempty"`
}
