package exadbvmclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExadbVMClusterProperties struct {
	BackupSubnetCidr          *string                         `json:"backupSubnetCidr,omitempty"`
	BackupSubnetOcid          *string                         `json:"backupSubnetOcid,omitempty"`
	ClusterName               *string                         `json:"clusterName,omitempty"`
	DataCollectionOptions     *DataCollectionOptions          `json:"dataCollectionOptions,omitempty"`
	DisplayName               string                          `json:"displayName"`
	Domain                    *string                         `json:"domain,omitempty"`
	EnabledEcpuCount          int64                           `json:"enabledEcpuCount"`
	ExascaleDbStorageVaultId  string                          `json:"exascaleDbStorageVaultId"`
	GiVersion                 *string                         `json:"giVersion,omitempty"`
	GridImageOcid             *string                         `json:"gridImageOcid,omitempty"`
	GridImageType             *GridImageType                  `json:"gridImageType,omitempty"`
	Hostname                  string                          `json:"hostname"`
	IormConfigCache           *ExadataIormConfig              `json:"iormConfigCache,omitempty"`
	LicenseModel              *LicenseModel                   `json:"licenseModel,omitempty"`
	LifecycleDetails          *string                         `json:"lifecycleDetails,omitempty"`
	LifecycleState            *ExadbVMClusterLifecycleState   `json:"lifecycleState,omitempty"`
	ListenerPort              *int64                          `json:"listenerPort,omitempty"`
	MemorySizeInGbs           *int64                          `json:"memorySizeInGbs,omitempty"`
	NodeCount                 int64                           `json:"nodeCount"`
	NsgCidrs                  *[]NsgCidr                      `json:"nsgCidrs,omitempty"`
	NsgURL                    *string                         `json:"nsgUrl,omitempty"`
	OciURL                    *string                         `json:"ociUrl,omitempty"`
	Ocid                      *string                         `json:"ocid,omitempty"`
	PrivateZoneOcid           *string                         `json:"privateZoneOcid,omitempty"`
	ProvisioningState         *AzureResourceProvisioningState `json:"provisioningState,omitempty"`
	ScanDnsName               *string                         `json:"scanDnsName,omitempty"`
	ScanDnsRecordId           *string                         `json:"scanDnsRecordId,omitempty"`
	ScanIPIds                 *[]string                       `json:"scanIpIds,omitempty"`
	ScanListenerPortTcp       *int64                          `json:"scanListenerPortTcp,omitempty"`
	ScanListenerPortTcpSsl    *int64                          `json:"scanListenerPortTcpSsl,omitempty"`
	Shape                     string                          `json:"shape"`
	SnapshotFileSystemStorage *ExadbVMClusterStorageDetails   `json:"snapshotFileSystemStorage,omitempty"`
	SshPublicKeys             []string                        `json:"sshPublicKeys"`
	SubnetId                  string                          `json:"subnetId"`
	SubnetOcid                *string                         `json:"subnetOcid,omitempty"`
	SystemVersion             *string                         `json:"systemVersion,omitempty"`
	TimeZone                  *string                         `json:"timeZone,omitempty"`
	TotalEcpuCount            int64                           `json:"totalEcpuCount"`
	TotalFileSystemStorage    *ExadbVMClusterStorageDetails   `json:"totalFileSystemStorage,omitempty"`
	VMFileSystemStorage       ExadbVMClusterStorageDetails    `json:"vmFileSystemStorage"`
	VipIds                    *[]string                       `json:"vipIds,omitempty"`
	VnetId                    string                          `json:"vnetId"`
	ZoneOcid                  *string                         `json:"zoneOcid,omitempty"`
}
