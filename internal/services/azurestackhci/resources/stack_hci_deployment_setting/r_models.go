// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package stack_hci_deployment_setting

type StackHCIDeploymentSettingModel struct {
	StackHCIClusterId string           `tfschema:"stack_hci_cluster_id"`
	ArcResourceIds    []string         `tfschema:"arc_resource_ids"`
	Version           string           `tfschema:"version"`
	ScaleUnit         []ScaleUnitModel `tfschema:"scale_unit"`
}

type ScaleUnitModel struct {
	ActiveDirectoryOrganizationalUnitPath string                       `tfschema:"active_directory_organizational_unit_path"`
	Cluster                               []ClusterModel               `tfschema:"cluster"`
	DomainFqdn                            string                       `tfschema:"domain_fqdn"`
	HostNetwork                           []HostNetworkModel           `tfschema:"host_network"`
	InfrastructureNetwork                 []InfrastructureNetworkModel `tfschema:"infrastructure_network"`
	NamePrefix                            string                       `tfschema:"name_prefix"`
	OptionalService                       []OptionalServiceModel       `tfschema:"optional_service"`
	PhysicalNode                          []PhysicalNodeModel          `tfschema:"physical_node"`
	SecretsLocation                       string                       `tfschema:"secrets_location"`
	Storage                               []StorageModel               `tfschema:"storage"`

	// flatten 'observability' block, for API always return them
	StreamingDataClientEnabled bool `tfschema:"streaming_data_client_enabled"`
	EuLocationEnabled          bool `tfschema:"eu_location_enabled"`
	EpisodicDataUploadEnabled  bool `tfschema:"episodic_data_upload_enabled"`

	// flatten 'securitySetting' block, for API always return them
	BitlockerBootVolumeEnabled   bool `tfschema:"bitlocker_boot_volume_enabled"`
	BitlockerDataVolumeEnabled   bool `tfschema:"bitlocker_data_volume_enabled"`
	CredentialGuardEnabled       bool `tfschema:"credential_guard_enabled"`
	DriftControlEnabled          bool `tfschema:"drift_control_enabled"`
	DrtmProtectionEnabled        bool `tfschema:"drtm_protection_enabled"`
	HvciProtectionEnabled        bool `tfschema:"hvci_protection_enabled"`
	SideChannelMitigationEnabled bool `tfschema:"side_channel_mitigation_enabled"`
	SmbSigningEnabled            bool `tfschema:"smb_signing_enabled"`
	SmbClusterEncryptionEnabled  bool `tfschema:"smb_cluster_encryption_enabled"`
	WdacEnabled                  bool `tfschema:"wdac_enabled"`
}

type ClusterModel struct {
	AzureServiceEndpoint string `tfschema:"azure_service_endpoint"`
	CloudAccountName     string `tfschema:"cloud_account_name"`
	Name                 string `tfschema:"name"`
	WitnessType          string `tfschema:"witness_type"`
	WitnessPath          string `tfschema:"witness_path"`
}

type HostNetworkModel struct {
	Intent                               []HostNetworkIntentModel         `tfschema:"intent"`
	StorageAutoIpEnabled                 bool                             `tfschema:"storage_auto_ip_enabled"`
	StorageConnectivitySwitchlessEnabled bool                             `tfschema:"storage_connectivity_switchless_enabled"`
	StorageNetwork                       []HostNetworkStorageNetworkModel `tfschema:"storage_network"`
}

type HostNetworkIntentModel struct {
	Adapter                                   []string                                                   `tfschema:"adapter"`
	AdapterPropertyOverride                   []HostNetworkIntentAdapterPropertyOverrideModel            `tfschema:"adapter_property_override"`
	AdapterPropertyOverrideEnabled            bool                                                       `tfschema:"adapter_property_override_enabled"`
	Name                                      string                                                     `tfschema:"name"`
	QosPolicyOverride                         []HostNetworkIntentQosPolicyOverrideModel                  `tfschema:"qos_policy_override"`
	QosPolicyOverrideEnabled                  bool                                                       `tfschema:"qos_policy_override_enabled"`
	TrafficType                               []string                                                   `tfschema:"traffic_type"`
	VirtualSwitchConfigurationOverride        []HostNetworkIntentVirtualSwitchConfigurationOverrideModel `tfschema:"virtual_switch_configuration_override"`
	VirtualSwitchConfigurationOverrideEnabled bool                                                       `tfschema:"virtual_switch_configuration_override_enabled"`
}

type HostNetworkIntentAdapterPropertyOverrideModel struct {
	JumboPacket             string `tfschema:"jumbo_packet"`
	NetworkDirect           string `tfschema:"network_direct"`
	NetworkDirectTechnology string `tfschema:"network_direct_technology"`
}

type HostNetworkIntentQosPolicyOverrideModel struct {
	BandWidthPercentageSMB         string `tfschema:"bandwidth_percentage_smb"`
	PriorityValue8021ActionCluster string `tfschema:"priority_value8021_action_cluster"`
	PriorityValue8021ActionSMB     string `tfschema:"priority_value8021_action_smb"`
}

type HostNetworkIntentVirtualSwitchConfigurationOverrideModel struct {
	EnableIov              string `tfschema:"enable_iov"`
	LoadBalancingAlgorithm string `tfschema:"load_balancing_algorithm"`
}

type HostNetworkStorageNetworkModel struct {
	Name               string `tfschema:"name"`
	NetworkAdapterName string `tfschema:"network_adapter_name"`
	VlanId             string `tfschema:"vlan_id"`
}

type InfrastructureNetworkModel struct {
	DhcpEnabled bool          `tfschema:"dhcp_enabled"`
	SubnetMask  string        `tfschema:"subnet_mask"`
	Gateway     string        `tfschema:"gateway"`
	IpPool      []IpPoolModel `tfschema:"ip_pool"`
	DnsServer   []string      `tfschema:"dns_server"`
}

type IpPoolModel struct {
	StartingAddress string `tfschema:"starting_address"`
	EndingAddress   string `tfschema:"ending_address"`
}

type OptionalServiceModel struct {
	CustomLocation string `tfschema:"custom_location"`
}

type PhysicalNodeModel struct {
	Name        string `tfschema:"name"`
	Ipv4Address string `tfschema:"ipv4_address"`
}

type StorageModel struct {
	ConfigurationMode string `tfschema:"configuration_mode"`
}

func (Resource) ModelObject() any {
	return &StackHCIDeploymentSettingModel{}
}
