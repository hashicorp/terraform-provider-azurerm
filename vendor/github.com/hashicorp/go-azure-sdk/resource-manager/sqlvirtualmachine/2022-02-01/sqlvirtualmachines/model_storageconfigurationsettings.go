package sqlvirtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageConfigurationSettings struct {
	DiskConfigurationType *DiskConfigurationType `json:"diskConfigurationType,omitempty"`
	SqlDataSettings       *SQLStorageSettings    `json:"sqlDataSettings,omitempty"`
	SqlLogSettings        *SQLStorageSettings    `json:"sqlLogSettings,omitempty"`
	SqlSystemDbOnDataDisk *bool                  `json:"sqlSystemDbOnDataDisk,omitempty"`
	SqlTempDbSettings     *SQLTempDbSettings     `json:"sqlTempDbSettings,omitempty"`
	StorageWorkloadType   *StorageWorkloadType   `json:"storageWorkloadType,omitempty"`
}
