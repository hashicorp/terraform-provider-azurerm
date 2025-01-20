package cloudvmclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CloudVMClusterUpdateProperties struct {
	ComputeNodes           *[]string              `json:"computeNodes,omitempty"`
	CpuCoreCount           *int64                 `json:"cpuCoreCount,omitempty"`
	DataCollectionOptions  *DataCollectionOptions `json:"dataCollectionOptions,omitempty"`
	DataStorageSizeInTbs   *float64               `json:"dataStorageSizeInTbs,omitempty"`
	DbNodeStorageSizeInGbs *int64                 `json:"dbNodeStorageSizeInGbs,omitempty"`
	DisplayName            *string                `json:"displayName,omitempty"`
	LicenseModel           *LicenseModel          `json:"licenseModel,omitempty"`
	MemorySizeInGbs        *int64                 `json:"memorySizeInGbs,omitempty"`
	OcpuCount              *float64               `json:"ocpuCount,omitempty"`
	SshPublicKeys          *[]string              `json:"sshPublicKeys,omitempty"`
	StorageSizeInGbs       *int64                 `json:"storageSizeInGbs,omitempty"`
}
