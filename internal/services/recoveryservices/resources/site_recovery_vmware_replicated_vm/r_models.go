// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_vmware_replicated_vm

type IncludedDiskModel struct {
	DiskId                    string `tfschema:"disk_id"`
	LogStorageAccountId       string `tfschema:"log_storage_account_id"`
	TargetDiskEncryptionSetId string `tfschema:"target_disk_encryption_set_id"`
	TargetDiskType            string `tfschema:"target_disk_type"`
}

type NetworkInterfaceModel struct {
	SourceMacAddress string `tfschema:"source_mac_address"`
	TargetStaticIp   string `tfschema:"target_static_ip"`
	TargetSubnetName string `tfschema:"target_subnet_name"`
	TestSubnetName   string `tfschema:"test_subnet_name"`
	IsPrimary        bool   `tfschema:"is_primary"`
}

type SiteRecoveryReplicatedVmVMwareModel struct {
	Name                                  string                  `tfschema:"name"`
	RecoveryVaultId                       string                  `tfschema:"recovery_vault_id"`
	SourceVmName                          string                  `tfschema:"source_vm_name"`
	ApplianceName                         string                  `tfschema:"appliance_name"`
	RecoveryReplicationPolicyId           string                  `tfschema:"recovery_replication_policy_id"`
	PhysicalServerCredentialName          string                  `tfschema:"physical_server_credential_name"`
	LicenseType                           string                  `tfschema:"license_type"`
	TargetResourceGroupId                 string                  `tfschema:"target_resource_group_id"`
	TargetVmName                          string                  `tfschema:"target_vm_name"`
	MultiVmGroupName                      string                  `tfschema:"multi_vm_group_name"`
	TargetProximityPlacementGroupId       string                  `tfschema:"target_proximity_placement_group_id"`
	TargetVmSize                          string                  `tfschema:"target_vm_size"`
	TargetAvailabilitySetId               string                  `tfschema:"target_availability_set_id"`
	TargetZone                            string                  `tfschema:"target_zone"`
	TargetNetworkId                       string                  `tfschema:"target_network_id"`
	TestNetworkId                         string                  `tfschema:"test_network_id"`
	TargetBootDiagnosticsStorageAccountId string                  `tfschema:"target_boot_diagnostics_storage_account_id"`
	DiskToInclude                         []IncludedDiskModel     `tfschema:"managed_disk"`
	NetworkInterface                      []NetworkInterfaceModel `tfschema:"network_interface"`
	DefaultLogStorageAccountId            string                  `tfschema:"default_log_storage_account_id"`
	DefaultRecoveryDiskType               string                  `tfschema:"default_recovery_disk_type"`
	DefaultTargetDiskEncryptionSetId      string                  `tfschema:"default_target_disk_encryption_set_id"`
}

func (r Resource) ModelObject() any {
	return &SiteRecoveryReplicatedVmVMwareModel{}
}
