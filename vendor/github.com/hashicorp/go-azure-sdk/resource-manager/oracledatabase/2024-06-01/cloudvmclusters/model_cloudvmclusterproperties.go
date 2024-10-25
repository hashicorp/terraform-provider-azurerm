package cloudvmclusters

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CloudVMClusterProperties struct {
	BackupSubnetCidr             *string                         `json:"backupSubnetCidr,omitempty"`
	CloudExadataInfrastructureId string                          `json:"cloudExadataInfrastructureId"`
	ClusterName                  *string                         `json:"clusterName,omitempty"`
	CompartmentId                *string                         `json:"compartmentId,omitempty"`
	ComputeNodes                 *[]string                       `json:"computeNodes,omitempty"`
	CpuCoreCount                 int64                           `json:"cpuCoreCount"`
	DataCollectionOptions        *DataCollectionOptions          `json:"dataCollectionOptions,omitempty"`
	DataStoragePercentage        *int64                          `json:"dataStoragePercentage,omitempty"`
	DataStorageSizeInTbs         *float64                        `json:"dataStorageSizeInTbs,omitempty"`
	DbNodeStorageSizeInGbs       *int64                          `json:"dbNodeStorageSizeInGbs,omitempty"`
	DbServers                    *[]string                       `json:"dbServers,omitempty"`
	DiskRedundancy               *DiskRedundancy                 `json:"diskRedundancy,omitempty"`
	DisplayName                  string                          `json:"displayName"`
	Domain                       *string                         `json:"domain,omitempty"`
	GiVersion                    string                          `json:"giVersion"`
	Hostname                     string                          `json:"hostname"`
	IormConfigCache              *ExadataIormConfig              `json:"iormConfigCache,omitempty"`
	IsLocalBackupEnabled         *bool                           `json:"isLocalBackupEnabled,omitempty"`
	IsSparseDiskgroupEnabled     *bool                           `json:"isSparseDiskgroupEnabled,omitempty"`
	LastUpdateHistoryEntryId     *string                         `json:"lastUpdateHistoryEntryId,omitempty"`
	LicenseModel                 *LicenseModel                   `json:"licenseModel,omitempty"`
	LifecycleDetails             *string                         `json:"lifecycleDetails,omitempty"`
	LifecycleState               *CloudVMClusterLifecycleState   `json:"lifecycleState,omitempty"`
	ListenerPort                 *int64                          `json:"listenerPort,omitempty"`
	MemorySizeInGbs              *int64                          `json:"memorySizeInGbs,omitempty"`
	NodeCount                    *int64                          `json:"nodeCount,omitempty"`
	NsgCidrs                     *[]NsgCidr                      `json:"nsgCidrs,omitempty"`
	NsgURL                       *string                         `json:"nsgUrl,omitempty"`
	OciURL                       *string                         `json:"ociUrl,omitempty"`
	Ocid                         *string                         `json:"ocid,omitempty"`
	OcpuCount                    *float64                        `json:"ocpuCount,omitempty"`
	ProvisioningState            *AzureResourceProvisioningState `json:"provisioningState,omitempty"`
	ScanDnsName                  *string                         `json:"scanDnsName,omitempty"`
	ScanDnsRecordId              *string                         `json:"scanDnsRecordId,omitempty"`
	ScanIPIds                    *[]string                       `json:"scanIpIds,omitempty"`
	ScanListenerPortTcp          *int64                          `json:"scanListenerPortTcp,omitempty"`
	ScanListenerPortTcpSsl       *int64                          `json:"scanListenerPortTcpSsl,omitempty"`
	Shape                        *string                         `json:"shape,omitempty"`
	SshPublicKeys                []string                        `json:"sshPublicKeys"`
	StorageSizeInGbs             *int64                          `json:"storageSizeInGbs,omitempty"`
	SubnetId                     string                          `json:"subnetId"`
	SubnetOcid                   *string                         `json:"subnetOcid,omitempty"`
	SystemVersion                *string                         `json:"systemVersion,omitempty"`
	TimeCreated                  *string                         `json:"timeCreated,omitempty"`
	TimeZone                     *string                         `json:"timeZone,omitempty"`
	VipIds                       *[]string                       `json:"vipIds,omitempty"`
	VnetId                       string                          `json:"vnetId"`
	ZoneId                       *string                         `json:"zoneId,omitempty"`
}

func (o *CloudVMClusterProperties) GetTimeCreatedAsTime() (*time.Time, error) {
	if o.TimeCreated == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TimeCreated, "2006-01-02T15:04:05Z07:00")
}

func (o *CloudVMClusterProperties) SetTimeCreatedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TimeCreated = &formatted
}
