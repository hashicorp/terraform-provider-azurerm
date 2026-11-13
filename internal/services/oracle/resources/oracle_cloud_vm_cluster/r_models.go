// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_cloud_vm_cluster

type CloudVmClusterResourceModel struct {
	// Azure
	Location          string            `tfschema:"location"`
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Tags              map[string]string `tfschema:"tags"`

	// Required
	CloudExadataInfrastructureId string   `tfschema:"cloud_exadata_infrastructure_id"`
	CpuCoreCount                 int64    `tfschema:"cpu_core_count"`
	DataStorageSizeInTbs         float64  `tfschema:"data_storage_size_in_tbs"`
	DbNodeStorageSizeInGbs       int64    `tfschema:"db_node_storage_size_in_gbs"`
	DbServers                    []string `tfschema:"db_servers"`
	DisplayName                  string   `tfschema:"display_name"`
	GiVersion                    string   `tfschema:"gi_version"`
	Hostname                     string   `tfschema:"hostname"`
	HostnameActual               string   `tfschema:"hostname_actual"`
	LicenseModel                 string   `tfschema:"license_model"`
	MemorySizeInGbs              int64    `tfschema:"memory_size_in_gbs"`
	SshPublicKeys                []string `tfschema:"ssh_public_keys"`
	SubnetId                     string   `tfschema:"subnet_id"`
	VnetId                       string   `tfschema:"virtual_network_id"`

	// Optional
	BackupSubnetCidr         string                         `tfschema:"backup_subnet_cidr"`
	ClusterName              string                         `tfschema:"cluster_name"`
	DataCollectionOptions    []DataCollectionOptionsModel   `tfschema:"data_collection_options"`
	DataStoragePercentage    int64                          `tfschema:"data_storage_percentage"`
	Domain                   string                         `tfschema:"domain"`
	IsLocalBackupEnabled     bool                           `tfschema:"local_backup_enabled"`
	IsSparseDiskgroupEnabled bool                           `tfschema:"sparse_diskgroup_enabled"`
	Ocid                     string                         `tfschema:"ocid"`
	ScanListenerPortTcp      int64                          `tfschema:"scan_listener_port_tcp"`
	ScanListenerPortTcpSsl   int64                          `tfschema:"scan_listener_port_tcp_ssl"`
	SystemVersion            string                         `tfschema:"system_version"`
	TimeZone                 string                         `tfschema:"time_zone"`
	ZoneId                   string                         `tfschema:"zone_id"`
	FileSystemConfiguration  []FileSystemConfigurationModel `tfschema:"file_system_configuration"`
}

func (Resource) ModelObject() any {
	return &Resource{}
}
