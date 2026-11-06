// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package compute_fleet

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type ComputeFleetResourceModel struct {
	Name                       string                                     `tfschema:"name"`
	ResourceGroupName          string                                     `tfschema:"resource_group_name"`
	Identity                   []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Location                   string                                     `tfschema:"location"`
	Plan                       []PlanModel                                `tfschema:"plan"`
	AdditionalCapabilities     []AdditionalCapabilitiesModel              `tfschema:"additional_capabilities"`
	VirtualMachineProfile      []VirtualMachineProfileModel               `tfschema:"virtual_machine_profile"`
	ComputeApiVersion          string                                     `tfschema:"compute_api_version"`
	PlatformFaultDomainCount   int64                                      `tfschema:"platform_fault_domain_count"`
	OnDemandCapacity           []OnDemandCapacityModel                    `tfschema:"on_demand_capacity"`
	SpotCapacity               []SpotCapacityModel                        `tfschema:"spot_capacity"`
	UniqueId                   string                                     `tfschema:"unique_id"`
	VirtualMachineSizesProfile []VirtualMachineSizeProfileModel           `tfschema:"virtual_machine_sizes_profile"`
	Tags                       map[string]string                          `tfschema:"tags"`
	Zones                      []string                                   `tfschema:"zones"`
}

type AdditionalCapabilitiesModel struct {
	UltraSsdEnabled    bool `tfschema:"ultra_ssd_enabled"`
	HibernationEnabled bool `tfschema:"hibernation_enabled"`
}
type VirtualMachineProfileModel struct {
	GalleryApplicationProfile                []GalleryApplicationModel   `tfschema:"gallery_application"`
	CapacityReservationGroupId               string                      `tfschema:"capacity_reservation_group_id"`
	BootDiagnosticEnabled                    bool                        `tfschema:"boot_diagnostic_enabled"`
	BootDiagnosticStorageAccountEndpoint     string                      `tfschema:"boot_diagnostic_storage_account_endpoint"`
	Extension                                []ExtensionModel            `tfschema:"extension"`
	ExtensionsTimeBudgetDuration             string                      `tfschema:"extensions_time_budget_duration"`
	ExtensionOperationsEnabled               bool                        `tfschema:"extension_operations_enabled"`
	LicenseType                              string                      `tfschema:"license_type"`
	NetworkInterface                         []NetworkInterfaceModel     `tfschema:"network_interface"`
	OsProfile                                []OSProfileModel            `tfschema:"os_profile"`
	ScheduledEventTerminationTimeoutDuration string                      `tfschema:"scheduled_event_termination_timeout_duration"`
	ScheduledEventOsImageTimeoutDuration     string                      `tfschema:"scheduled_event_os_image_timeout_duration"`
	EncryptionAtHostEnabled                  bool                        `tfschema:"encryption_at_host_enabled"`
	SecureBootEnabled                        bool                        `tfschema:"secure_boot_enabled"`
	VTpmEnabled                              bool                        `tfschema:"vtpm_enabled"`
	DataDisks                                []DataDiskModel             `tfschema:"data_disk"`
	OsDisk                                   []OSDiskModel               `tfschema:"os_disk"`
	SourceImageReference                     []SourceImageReferenceModel `tfschema:"source_image_reference"`
	SourceImageId                            string                      `tfschema:"source_image_id"`
	UserDataBase64                           string                      `tfschema:"user_data_base64"`
	NetworkApiVersion                        string                      `tfschema:"network_api_version"`
}

type GalleryApplicationModel struct {
	ConfigurationBlobUri                   string `tfschema:"configuration_blob_uri"`
	AutomaticUpgradeEnabled                bool   `tfschema:"automatic_upgrade_enabled"`
	Order                                  int64  `tfschema:"order"`
	VersionId                              string `tfschema:"version_id"`
	Tag                                    string `tfschema:"tag"`
	TreatFailureAsDeploymentFailureEnabled bool   `tfschema:"treat_failure_as_deployment_failure_enabled"`
}

type ExtensionModel struct {
	Name                                             string                               `tfschema:"name"`
	Publisher                                        string                               `tfschema:"publisher"`
	Type                                             string                               `tfschema:"type"`
	TypeHandlerVersion                               string                               `tfschema:"type_handler_version"`
	AutoUpgradeMinorVersionEnabled                   bool                                 `tfschema:"auto_upgrade_minor_version_enabled"`
	AutomaticUpgradeEnabled                          bool                                 `tfschema:"automatic_upgrade_enabled"`
	ForceExtensionExecutionOnChange                  string                               `tfschema:"force_extension_execution_on_change"`
	ProtectedSettingsJson                            string                               `tfschema:"protected_settings_json"`
	ProtectedSettingsFromKeyVault                    []ProtectedSettingsFromKeyVaultModel `tfschema:"protected_settings_from_key_vault"`
	ExtensionsToProvisionAfterVirtualMachineCreation []string                             `tfschema:"extensions_to_provision_after_virtual_machine_creation"`
	FailureSuppressionEnabled                        bool                                 `tfschema:"failure_suppression_enabled"`
	SettingsJson                                     string                               `tfschema:"settings_json"`
}

type ProtectedSettingsFromKeyVaultModel struct {
	SecretUrl     string `tfschema:"secret_url"`
	SourceVaultId string `tfschema:"source_vault_id"`
}

type NetworkInterfaceModel struct {
	Name                           string                 `tfschema:"name"`
	AuxiliaryMode                  string                 `tfschema:"auxiliary_mode"`
	AuxiliarySku                   string                 `tfschema:"auxiliary_sku"`
	DeleteOption                   string                 `tfschema:"delete_option"`
	DnsServers                     []string               `tfschema:"dns_servers"`
	AcceleratedNetworkingEnabled   bool                   `tfschema:"accelerated_networking_enabled"`
	IPForwardingEnabled            bool                   `tfschema:"ip_forwarding_enabled"`
	IPConfiguration                []IPConfigurationModel `tfschema:"ip_configuration"`
	NetworkSecurityGroupId         string                 `tfschema:"network_security_group_id"`
	PrimaryNetworkInterfaceEnabled bool                   `tfschema:"primary_network_interface_enabled"`
}

type IPConfigurationModel struct {
	Name                                    string                 `tfschema:"name"`
	ApplicationGatewayBackendAddressPoolIds []string               `tfschema:"application_gateway_backend_address_pool_ids"`
	ApplicationSecurityGroupIds             []string               `tfschema:"application_security_group_ids"`
	LoadBalancerBackendAddressPoolIds       []string               `tfschema:"load_balancer_backend_address_pool_ids"`
	PrimaryIpConfigurationEnabled           bool                   `tfschema:"primary_ip_configuration_enabled"`
	Version                                 string                 `tfschema:"version"`
	PublicIPAddress                         []PublicIPAddressModel `tfschema:"public_ip_address"`
	SubnetId                                string                 `tfschema:"subnet_id"`
}

type PublicIPAddressModel struct {
	Name                 string `tfschema:"name"`
	DeleteOption         string `tfschema:"delete_option"`
	DomainNameLabel      string `tfschema:"domain_name_label"`
	DomainNameLabelScope string `tfschema:"domain_name_label_scope"`
	IdleTimeoutInMinutes int64  `tfschema:"idle_timeout_in_minutes"`
	Version              string `tfschema:"version"`
	PublicIPPrefix       string `tfschema:"public_ip_prefix_id"`
	SkuName              string `tfschema:"sku_name"`
}

type OSProfileModel struct {
	CustomDataBase64     string                      `tfschema:"custom_data_base64"`
	LinuxConfiguration   []LinuxConfigurationModel   `tfschema:"linux_configuration"`
	WindowsConfiguration []WindowsConfigurationModel `tfschema:"windows_configuration"`
}

type LinuxConfigurationModel struct {
	AdminPassword                             string             `tfschema:"admin_password"`
	AdminUsername                             string             `tfschema:"admin_username"`
	ComputerNamePrefix                        string             `tfschema:"computer_name_prefix"`
	Secret                                    []LinuxSecretModel `tfschema:"secret"`
	PasswordAuthenticationEnabled             bool               `tfschema:"password_authentication_enabled"`
	VirtualMachineAgentPlatformUpdatesEnabled bool               `tfschema:"virtual_machine_agent_platform_updates_enabled"`
	PatchMode                                 string             `tfschema:"patch_mode"`
	BypassPlatformSafetyChecksEnabled         bool               `tfschema:"bypass_platform_safety_checks_enabled"`
	PatchRebooting                            string             `tfschema:"patch_rebooting"`
	ProvisionVirtualMachineAgentEnabled       bool               `tfschema:"provision_virtual_machine_agent_enabled"`
	AdminSSHKeys                              []string           `tfschema:"admin_ssh_keys"`
}

type LinuxSecretModel struct {
	KeyVaultId  string                  `tfschema:"key_vault_id"`
	Certificate []LinuxCertificateModel `tfschema:"certificate"`
}

type WindowsSecretModel struct {
	KeyVaultId  string                    `tfschema:"key_vault_id"`
	Certificate []WindowsCertificateModel `tfschema:"certificate"`
}

type LinuxCertificateModel struct {
	Url string `tfschema:"url"`
}

type WindowsCertificateModel struct {
	Store string `tfschema:"store"`
	Url   string `tfschema:"url"`
}

type WindowsConfigurationModel struct {
	AdminPassword                             string                           `tfschema:"admin_password"`
	AdminUsername                             string                           `tfschema:"admin_username"`
	ComputerNamePrefix                        string                           `tfschema:"computer_name_prefix"`
	Secret                                    []WindowsSecretModel             `tfschema:"secret"`
	AdditionalUnattendContent                 []AdditionalUnattendContentModel `tfschema:"additional_unattend_content"`
	AutomaticUpdatesEnabled                   bool                             `tfschema:"automatic_updates_enabled"`
	VirtualMachineAgentPlatformUpdatesEnabled bool                             `tfschema:"virtual_machine_agent_platform_updates_enabled"`
	PatchMode                                 string                           `tfschema:"patch_mode"`
	BypassPlatformSafetyChecksEnabled         bool                             `tfschema:"bypass_platform_safety_checks_enabled"`
	PatchRebooting                            string                           `tfschema:"patch_rebooting"`
	HotPatchingEnabled                        bool                             `tfschema:"hot_patching_enabled"`
	ProvisionVirtualMachineAgentEnabled       bool                             `tfschema:"provision_virtual_machine_agent_enabled"`
	TimeZone                                  string                           `tfschema:"time_zone"`
	WinRM                                     []WinRMModel                     `tfschema:"winrm_listener"`
}

type AdditionalUnattendContentModel struct {
	Xml     string `tfschema:"xml"`
	Setting string `tfschema:"setting"`
}

type WinRMModel struct {
	CertificateUrl string `tfschema:"certificate_url"`
	Protocol       string `tfschema:"protocol"`
}

type DataDiskModel struct {
	Caching                 string `tfschema:"caching"`
	CreateOption            string `tfschema:"create_option"`
	DeleteOption            string `tfschema:"delete_option"`
	DiskSizeInGiB           int64  `tfschema:"disk_size_in_gib"`
	DiskEncryptionSetId     string `tfschema:"disk_encryption_set_id"`
	StorageAccountType      string `tfschema:"storage_account_type"`
	Lun                     int64  `tfschema:"lun"`
	WriteAcceleratorEnabled bool   `tfschema:"write_accelerator_enabled"`
}

type SourceImageReferenceModel struct {
	Offer     string `tfschema:"offer"`
	Publisher string `tfschema:"publisher"`
	Sku       string `tfschema:"sku"`
	Version   string `tfschema:"version"`
}

type OSDiskModel struct {
	Caching                 string `tfschema:"caching"`
	DeleteOption            string `tfschema:"delete_option"`
	DiffDiskOption          string `tfschema:"diff_disk_option"`
	DiffDiskPlacement       string `tfschema:"diff_disk_placement"`
	DiskSizeInGiB           int64  `tfschema:"disk_size_in_gib"`
	DiskEncryptionSetId     string `tfschema:"disk_encryption_set_id"`
	SecurityEncryptionType  string `tfschema:"security_encryption_type"`
	StorageAccountType      string `tfschema:"storage_account_type"`
	WriteAcceleratorEnabled bool   `tfschema:"write_accelerator_enabled"`
}

type PlanModel struct {
	Name          string `tfschema:"name"`
	Product       string `tfschema:"product"`
	PromotionCode string `tfschema:"promotion_code"`
	Publisher     string `tfschema:"publisher"`
}

type OnDemandCapacityModel struct {
	AllocationStrategy      string `tfschema:"allocation_strategy"`
	TargetCapacity          int64  `tfschema:"target_capacity"`
	MinimumStartingCapacity int64  `tfschema:"minimum_starting_capacity"`
}

type SpotCapacityModel struct {
	AllocationStrategy                     string  `tfschema:"allocation_strategy"`
	TargetCapacity                         int64   `tfschema:"target_capacity"`
	EvictionPolicy                         string  `tfschema:"eviction_policy"`
	MaintainCapacityEnabled                bool    `tfschema:"maintain_capacity_enabled"`
	MaximumHourlyPricePerVirtualMachineUsd float64 `tfschema:"maximum_hourly_price_per_virtual_machine_usd"`
	MinimumCapacity                        int64   `tfschema:"minimum_capacity"`
}

type VirtualMachineSizeProfileModel struct {
	Name string `tfschema:"name"`
	Rank int64  `tfschema:"rank"`
}

func (r Resource) ModelObject() any {
	return &ComputeFleetResourceModel{}
}
