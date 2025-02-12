package computefleet

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurefleet/2024-11-01/fleets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ComputeFleetResourceModel struct {
	Name                                     string                                     `tfschema:"name"`
	ResourceGroupName                        string                                     `tfschema:"resource_group_name"`
	Identity                                 []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Location                                 string                                     `tfschema:"location"`
	Plan                                     []PlanModel                                `tfschema:"plan"`
	AdditionalLocationProfile                []AdditionalLocationProfileModel           `tfschema:"additional_location_profile"`
	AdditionalCapabilitiesUltraSSDEnabled    bool                                       `tfschema:"additional_capabilities_ultra_ssd_enabled"`
	AdditionalCapabilitiesHibernationEnabled bool                                       `tfschema:"additional_capabilities_hibernation_enabled"`
	VirtualMachineProfile                    []VirtualMachineProfileModel               `tfschema:"virtual_machine_profile"`
	ComputeApiVersion                        string                                     `tfschema:"compute_api_version"`
	PlatformFaultDomainCount                 int64                                      `tfschema:"platform_fault_domain_count"`
	RegularPriorityProfile                   []RegularPriorityProfileModel              `tfschema:"regular_priority_profile"`
	SpotPriorityProfile                      []SpotPriorityProfileModel                 `tfschema:"spot_priority_profile"`
	UniqueId                                 string                                     `tfschema:"unique_id"`
	VMAttributes                             []VMAttributesModel                        `tfschema:"vm_attributes"`
	VMSizesProfile                           []VMSizeProfileModel                       `tfschema:"vm_sizes_profile"`
	Tags                                     map[string]string                          `tfschema:"tags"`
	Zones                                    []string                                   `tfschema:"zones"`
}

type AdditionalLocationProfileModel struct {
	Location                      string                       `tfschema:"location"`
	VirtualMachineProfileOverride []VirtualMachineProfileModel `tfschema:"virtual_machine_profile_override"`
}

type VirtualMachineProfileModel struct {
	GalleryApplicationProfile            []GalleryApplicationModel   `tfschema:"gallery_application"`
	CapacityReservationGroupId           string                      `tfschema:"capacity_reservation_group_id"`
	BootDiagnosticEnabled                bool                        `tfschema:"boot_diagnostic_enabled"`
	BootDiagnosticStorageAccountEndpoint string                      `tfschema:"boot_diagnostic_storage_account_endpoint"`
	Extension                            []ExtensionModel            `tfschema:"extension"`
	ExtensionsTimeBudget                 string                      `tfschema:"extensions_time_budget"`
	ExtensionOperationsEnabled           bool                        `tfschema:"extension_operations_enabled"`
	LicenseType                          string                      `tfschema:"license_type"`
	NetworkInterface                     []NetworkInterfaceModel     `tfschema:"network_interface"`
	OsProfile                            []OSProfileModel            `tfschema:"os_profile"`
	ScheduledEventTerminationTimeout     string                      `tfschema:"scheduled_event_termination_timeout"`
	ScheduledEventOsImageTimeout         string                      `tfschema:"scheduled_event_os_image_timeout"`
	EncryptionAtHostEnabled              bool                        `tfschema:"encryption_at_host_enabled"`
	SecureBootEnabled                    bool                        `tfschema:"secure_boot_enabled"`
	VTpmEnabled                          bool                        `tfschema:"vtpm_enabled"`
	DataDisks                            []DataDiskModel             `tfschema:"data_disk"`
	OsDisk                               []OSDiskModel               `tfschema:"os_disk"`
	SourceImageReference                 []SourceImageReferenceModel `tfschema:"source_image_reference"`
	SourceImageId                        string                      `tfschema:"source_image_id"`
	UserDataBase64                       string                      `tfschema:"user_data_base64"`
	NetworkApiVersion                    string                      `tfschema:"network_api_version"`
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
	Name                                 string                               `tfschema:"name"`
	Publisher                            string                               `tfschema:"publisher"`
	Type                                 string                               `tfschema:"type"`
	TypeHandlerVersion                   string                               `tfschema:"type_handler_version"`
	AutoUpgradeMinorVersionEnabled       bool                                 `tfschema:"auto_upgrade_minor_version_enabled"`
	AutomaticUpgradeEnabled              bool                                 `tfschema:"automatic_upgrade_enabled"`
	ForceExtensionExecutionOnChange      string                               `tfschema:"force_extension_execution_on_change"`
	ProtectedSettingsJson                string                               `tfschema:"protected_settings_json"`
	ProtectedSettingsFromKeyVault        []ProtectedSettingsFromKeyVaultModel `tfschema:"protected_settings_from_key_vault"`
	ExtensionsToProvisionAfterVmCreation []string                             `tfschema:"extensions_to_provision_after_vm_creation"`
	FailureSuppressionEnabled            bool                                 `tfschema:"failure_suppression_enabled"`
	SettingsJson                         string                               `tfschema:"settings_json"`
}

type ProtectedSettingsFromKeyVaultModel struct {
	SecretUrl     string `tfschema:"secret_url"`
	SourceVaultId string `tfschema:"source_vault_id"`
}

type NetworkInterfaceModel struct {
	Name                         string                 `tfschema:"name"`
	AuxiliaryMode                string                 `tfschema:"auxiliary_mode"`
	AuxiliarySku                 string                 `tfschema:"auxiliary_sku"`
	DeleteOption                 string                 `tfschema:"delete_option"`
	DnsServers                   []string               `tfschema:"dns_servers"`
	AcceleratedNetworkingEnabled bool                   `tfschema:"accelerated_networking_enabled"`
	IPForwardingEnabled          bool                   `tfschema:"ip_forwarding_enabled"`
	IPConfiguration              []IPConfigurationModel `tfschema:"ip_configuration"`
	NetworkSecurityGroupId       string                 `tfschema:"network_security_group_id"`
	Primary                      bool                   `tfschema:"primary"`
}

type IPConfigurationModel struct {
	Name                                    string                 `tfschema:"name"`
	ApplicationGatewayBackendAddressPoolIds []string               `tfschema:"application_gateway_backend_address_pool_ids"`
	ApplicationSecurityGroupIds             []string               `tfschema:"application_security_group_ids"`
	LoadBalancerBackendAddressPoolIds       []string               `tfschema:"load_balancer_backend_address_pool_ids"`
	Primary                                 bool                   `tfschema:"primary"`
	Version                                 string                 `tfschema:"version"`
	PublicIPAddress                         []PublicIPAddressModel `tfschema:"public_ip_address"`
	SubnetId                                string                 `tfschema:"subnet_id"`
}

type PublicIPAddressModel struct {
	Name                 string       `tfschema:"name"`
	DeleteOption         string       `tfschema:"delete_option"`
	DomainNameLabel      string       `tfschema:"domain_name_label"`
	DomainNameLabelScope string       `tfschema:"domain_name_label_scope"`
	IdleTimeoutInMinutes int64        `tfschema:"idle_timeout_in_minutes"`
	IPTag                []IPTagModel `tfschema:"ip_tag"`
	Version              string       `tfschema:"version"`
	PublicIPPrefix       string       `tfschema:"public_ip_prefix_id"`
	SkuName              string       `tfschema:"sku_name"`
}

type IPTagModel struct {
	Type string `tfschema:"type"`
	Tag  string `tfschema:"tag"`
}

type OSProfileModel struct {
	CustomDataBase64     string                      `tfschema:"custom_data_base64"`
	LinuxConfiguration   []LinuxConfigurationModel   `tfschema:"linux_configuration"`
	WindowsConfiguration []WindowsConfigurationModel `tfschema:"windows_configuration"`
}

type LinuxConfigurationModel struct {
	AdminPassword                     string             `tfschema:"admin_password"`
	AdminUsername                     string             `tfschema:"admin_username"`
	ComputerNamePrefix                string             `tfschema:"computer_name_prefix"`
	Secret                            []LinuxSecretModel `tfschema:"secret"`
	PasswordAuthenticationEnabled     bool               `tfschema:"password_authentication_enabled"`
	VMAgentPlatformUpdatesEnabled     bool               `tfschema:"vm_agent_platform_updates_enabled"`
	PatchAssessmentMode               string             `tfschema:"patch_assessment_mode"`
	PatchMode                         string             `tfschema:"patch_mode"`
	BypassPlatformSafetyChecksEnabled bool               `tfschema:"bypass_platform_safety_checks_enabled"`
	RebootSetting                     string             `tfschema:"reboot_setting"`
	ProvisionVMAgentEnabled           bool               `tfschema:"provision_vm_agent_enabled"`
	AdminSSHKeys                      []string           `tfschema:"admin_ssh_keys"`
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
	AdminPassword                     string                           `tfschema:"admin_password"`
	AdminUsername                     string                           `tfschema:"admin_username"`
	ComputerNamePrefix                string                           `tfschema:"computer_name_prefix"`
	Secret                            []WindowsSecretModel             `tfschema:"secret"`
	AdditionalUnattendContent         []AdditionalUnattendContentModel `tfschema:"additional_unattend_content"`
	AutomaticUpdatesEnabled           bool                             `tfschema:"automatic_updates_enabled"`
	VMAgentPlatformUpdatesEnabled     bool                             `tfschema:"vm_agent_platform_updates_enabled"`
	PatchAssessmentMode               string                           `tfschema:"patch_assessment_mode"`
	PatchMode                         string                           `tfschema:"patch_mode"`
	BypassPlatformSafetyChecksEnabled bool                             `tfschema:"bypass_platform_safety_checks_enabled"`
	RebootSetting                     string                           `tfschema:"reboot_setting"`
	HotPatchingEnabled                bool                             `tfschema:"hot_patching_enabled"`
	ProvisionVMAgentEnabled           bool                             `tfschema:"provision_vm_agent_enabled"`
	TimeZone                          string                           `tfschema:"time_zone"`
	WinRM                             []WinRMModel                     `tfschema:"winrm_listener"`
}

type AdditionalUnattendContentModel struct {
	Content string `tfschema:"content"`
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
	DiskSizeInGB            int64  `tfschema:"disk_size_in_gb"`
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
	DiskSizeInGB            int64  `tfschema:"disk_size_in_gb"`
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

type RegularPriorityProfileModel struct {
	AllocationStrategy string `tfschema:"allocation_strategy"`
	Capacity           int64  `tfschema:"capacity"`
	MinCapacity        int64  `tfschema:"min_capacity"`
}

type SpotPriorityProfileModel struct {
	AllocationStrategy  string  `tfschema:"allocation_strategy"`
	Capacity            int64   `tfschema:"capacity"`
	EvictionPolicy      string  `tfschema:"eviction_policy"`
	MaintainEnabled     bool    `tfschema:"maintain_enabled"`
	MaxHourlyPricePerVM float64 `tfschema:"max_hourly_price_per_vm"`
	MinCapacity         int64   `tfschema:"min_capacity"`
}

type VMAttributesModel struct {
	AcceleratorCount          []VMAttributeMinMaxIntegerModel `tfschema:"accelerator_count"`
	AcceleratorManufacturers  []string                        `tfschema:"accelerator_manufacturers"`
	AcceleratorSupport        string                          `tfschema:"accelerator_support"`
	AcceleratorTypes          []string                        `tfschema:"accelerator_types"`
	ArchitectureTypes         []string                        `tfschema:"architecture_types"`
	BurstableSupport          string                          `tfschema:"burstable_support"`
	CpuManufacturers          []string                        `tfschema:"cpu_manufacturers"`
	DataDiskCount             []VMAttributeMinMaxIntegerModel `tfschema:"data_disk_count"`
	ExcludedVMSizes           []string                        `tfschema:"excluded_vm_sizes"`
	LocalStorageDiskTypes     []string                        `tfschema:"local_storage_disk_types"`
	LocalStorageInGib         []VMAttributeMinMaxDoubleModel  `tfschema:"local_storage_in_gib"`
	LocalStorageSupport       string                          `tfschema:"local_storage_support"`
	MemoryInGib               []VMAttributeMinMaxDoubleModel  `tfschema:"memory_in_gib"`
	MemoryInGibPerVCPU        []VMAttributeMinMaxDoubleModel  `tfschema:"memory_in_gib_per_vcpu"`
	NetworkBandwidthInMbps    []VMAttributeMinMaxDoubleModel  `tfschema:"network_bandwidth_in_mbps"`
	NetworkInterfaceCount     []VMAttributeMinMaxIntegerModel `tfschema:"network_interface_count"`
	RdmaNetworkInterfaceCount []VMAttributeMinMaxIntegerModel `tfschema:"rdma_network_interface_count"`
	RdmaSupport               string                          `tfschema:"rdma_support"`
	VCPUCount                 []VMAttributeMinMaxIntegerModel `tfschema:"vcpu_count"`
	VMCategories              []string                        `tfschema:"vm_categories"`
}

type VMAttributeMinMaxIntegerModel struct {
	Max int64 `tfschema:"max"`
	Min int64 `tfschema:"min"`
}

type VMAttributeMinMaxDoubleModel struct {
	Max float64 `tfschema:"max"`
	Min float64 `tfschema:"min"`
}

type VMSizeProfileModel struct {
	Name string `tfschema:"name"`
	Rank int64  `tfschema:"rank"`
}

type ComputeFleetResource struct{}

var _ sdk.ResourceWithUpdate = ComputeFleetResource{}

var _ sdk.ResourceWithCustomizeDiff = ComputeFleetResource{}

func (r ComputeFleetResource) ResourceType() string {
	return "azurerm_compute_fleet"
}

func (r ComputeFleetResource) ModelObject() interface{} {
	return &ComputeFleetResourceModel{}
}

func (r ComputeFleetResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return fleets.ValidateFleetID
}

func (r ComputeFleetResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9.-]{0,62}[a-zA-Z0-9]$"),
				"The fleet name can only start or end with a number or a letter, and can contain only letters, numbers, periods (.), hyphens (-), up to 64 characters",
			),
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"virtual_machine_profile": virtualMachineProfileSchema(),

		"additional_capabilities_hibernation_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},

		"additional_capabilities_ultra_ssd_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},

		"additional_location_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"location": commonschema.LocationWithoutForceNew(),

					"virtual_machine_profile_override": virtualMachineProfileSchema(),
				},
			},
		},

		"plan": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						ForceNew: true,
						Required: true,
					},

					"product": {
						Type:     pluginsdk.TypeString,
						ForceNew: true,
						Required: true,
					},

					"publisher": {
						Type:     pluginsdk.TypeString,
						ForceNew: true,
						Required: true,
					},

					"promotion_code": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"platform_fault_domain_count": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Default:  1,
			ForceNew: true,
		},

		"vm_attributes": vmAttributesSchema(),

		"zones": commonschema.ZonesMultipleOptionalForceNew(),

		"compute_api_version": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"regular_priority_profile": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MaxItems:     1,
			AtLeastOneOf: []string{"regular_priority_profile", "spot_priority_profile"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"allocation_strategy": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						Default:  string(fleets.RegularPriorityAllocationStrategyLowestPrice),
						ValidateFunc: validation.StringInSlice([]string{
							string(fleets.RegularPriorityAllocationStrategyLowestPrice),
							string(fleets.RegularPriorityAllocationStrategyPrioritized),
						}, false),
					},

					"min_capacity": {
						Type:     pluginsdk.TypeInt,
						ForceNew: true,
						Optional: true,
					},

					"capacity": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 10000),
					},
				},
			},
		},

		"spot_priority_profile": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MaxItems:     1,
			AtLeastOneOf: []string{"regular_priority_profile", "spot_priority_profile"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"allocation_strategy": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						Default:  string(fleets.SpotAllocationStrategyPriceCapacityOptimized),
						ValidateFunc: validation.StringInSlice([]string{
							string(fleets.SpotAllocationStrategyPriceCapacityOptimized),
							string(fleets.SpotAllocationStrategyLowestPrice),
							string(fleets.SpotAllocationStrategyCapacityOptimized),
						}, false),
					},

					"eviction_policy": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						Default:  string(fleets.EvictionPolicyDelete),
						ValidateFunc: validation.StringInSlice([]string{
							string(fleets.EvictionPolicyDelete),
							string(fleets.EvictionPolicyDeallocate),
						}, false),
					},

					"maintain_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  true,
					},

					"max_hourly_price_per_vm": {
						Type:         pluginsdk.TypeFloat,
						Optional:     true,
						ForceNew:     true,
						Default:      -1,
						ValidateFunc: computeValidate.SpotMaxPrice,
					},

					"min_capacity": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},

					"capacity": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 10000),
					},
				},
			},
		},

		"tags": commonschema.Tags(),

		"vm_sizes_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"rank": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 65535),
					},
				},
			},
			ConflictsWith: []string{"vm_attributes.0.excluded_vm_sizes"},
			AtLeastOneOf:  []string{"vm_sizes_profile", "vm_attributes"},
		},
	}
}

func vmAttributesSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:         pluginsdk.TypeList,
		Optional:     true,
		ForceNew:     true,
		MaxItems:     1,
		AtLeastOneOf: []string{"vm_sizes_profile", "vm_attributes"},
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"memory_in_gib": {
					Type:     pluginsdk.TypeList,
					Required: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: vmAttributesMaxMinFloatSchema("memory_in_gib"),
					},
				},

				"vcpu_count": {
					Type:     pluginsdk.TypeList,
					Required: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: vmAttributesMaxMinIntegerSchema("vcpu_count"),
					},
				},

				"accelerator_count": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: vmAttributesMaxMinIntegerSchema("accelerator_count"),
					},
				},

				"accelerator_manufacturers": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
						ValidateFunc: validation.StringInSlice(
							fleets.PossibleValuesForAcceleratorManufacturer(), false),
					},
				},

				"accelerator_support": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  string(fleets.VMAttributeSupportExcluded),
					ValidateFunc: validation.StringInSlice(
						fleets.PossibleValuesForVMAttributeSupport(), false),
				},

				"accelerator_types": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
						ValidateFunc: validation.StringInSlice(
							fleets.PossibleValuesForAcceleratorType(), false),
					},
				},

				"architecture_types": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
						ValidateFunc: validation.StringInSlice(
							fleets.PossibleValuesForArchitectureType(), false),
					},
				},

				"burstable_support": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  string(fleets.VMAttributeSupportExcluded),
					ValidateFunc: validation.StringInSlice(
						fleets.PossibleValuesForVMAttributeSupport(), false),
				},

				"cpu_manufacturers": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
						ValidateFunc: validation.StringInSlice(
							fleets.PossibleValuesForCPUManufacturer(), false),
					},
				},

				"data_disk_count": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: vmAttributesMaxMinIntegerSchema("data_disk_count"),
					},
				},

				"excluded_vm_sizes": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					ConflictsWith: []string{"vm_sizes_profile"},
				},

				"local_storage_disk_types": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
						ValidateFunc: validation.StringInSlice(
							fleets.PossibleValuesForLocalStorageDiskType(), false),
					},
				},

				"local_storage_in_gib": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: vmAttributesMaxMinFloatSchema("local_storage_in_gib"),
					},
				},

				"local_storage_support": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  string(fleets.VMAttributeSupportIncluded),
					ValidateFunc: validation.StringInSlice(
						fleets.PossibleValuesForVMAttributeSupport(), false),
				},

				"memory_in_gib_per_vcpu": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: vmAttributesMaxMinFloatSchema("memory_in_gib_per_vcpu"),
					},
				},

				"network_bandwidth_in_mbps": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: vmAttributesMaxMinFloatSchema("network_bandwidth_in_mbps"),
					},
				},

				"network_interface_count": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: vmAttributesMaxMinIntegerSchema("network_interface_count"),
					},
				},

				"rdma_network_interface_count": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: vmAttributesMaxMinIntegerSchema("rdma_network_interface_count"),
					},
				},

				"rdma_support": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  string(fleets.VMAttributeSupportExcluded),
					ValidateFunc: validation.StringInSlice(
						fleets.PossibleValuesForVMAttributeSupport(), false),
				},

				"vm_categories": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
						ValidateFunc: validation.StringInSlice(
							fleets.PossibleValuesForVMCategory(), false),
					},
				},
			},
		},
	}
}

func vmAttributesMaxMinIntegerSchema(parent string) map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"max": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(0),
			AtLeastOneOf: []string{"vm_attributes.0." + parent + ".0.max", "vm_attributes.0." + parent + ".0.min"},
		},

		"min": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(0),
			AtLeastOneOf: []string{"vm_attributes.0." + parent + ".0.max", "vm_attributes.0." + parent + ".0.min"},
		},
	}
}

func vmAttributesMaxMinFloatSchema(parent string) map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"max": {
			Type:         pluginsdk.TypeFloat,
			Optional:     true,
			ValidateFunc: validation.FloatAtLeast(0.0),
			AtLeastOneOf: []string{"vm_attributes.0." + parent + ".0.max", "vm_attributes.0." + parent + ".0.min"},
		},

		"min": {
			Type:         pluginsdk.TypeFloat,
			Optional:     true,
			ValidateFunc: validation.FloatAtLeast(0.0),
			AtLeastOneOf: []string{"vm_attributes.0." + parent + ".0.max", "vm_attributes.0." + parent + ".0.min"},
		},
	}
}

func (r ComputeFleetResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"unique_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r ComputeFleetResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ComputeFleet.ComputeFleetClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model ComputeFleetResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := fleets.NewFleetID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := fleets.Fleet{
				Location: location.Normalize(model.Location),
				Plan:     expandPlanModel(model.Plan),
				Properties: &fleets.FleetProperties{
					RegularPriorityProfile: expandRegularPriorityProfileModel(model.RegularPriorityProfile),
					SpotPriorityProfile:    expandSpotPriorityProfileModel(model.SpotPriorityProfile),
					VMAttributes:           expandVMAttributesModel(model.VMAttributes),
					VMSizesProfile:         pointer.From(expandVMSizeProfileModel(model.VMSizesProfile)),
				},
			}

			if model.Tags != nil {
				properties.Tags = pointer.To(model.Tags)
			}

			if model.Zones != nil {
				properties.Zones = pointer.To(model.Zones)
			}

			expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}
			properties.Identity = expandedIdentity

			additionalLocationsProfileValue, err := expandAdditionalLocationProfileModel(model.AdditionalLocationProfile, metadata.ResourceData)
			if err != nil {
				return err
			}
			properties.Properties.AdditionalLocationsProfile = additionalLocationsProfileValue

			computeProfile := fleets.ComputeProfile{
				AdditionalVirtualMachineCapabilities: &fleets.AdditionalCapabilities{
					HibernationEnabled: pointer.To(model.AdditionalCapabilitiesHibernationEnabled),
					UltraSSDEnabled:    pointer.To(model.AdditionalCapabilitiesUltraSSDEnabled),
				},
				PlatformFaultDomainCount: pointer.To(model.PlatformFaultDomainCount),
			}
			if model.ComputeApiVersion != "" {
				computeProfile.ComputeApiVersion = pointer.To(model.ComputeApiVersion)
			}

			baseVirtualMachineProfileValue, err := expandVirtualMachineProfileModel(model.VirtualMachineProfile, metadata.ResourceData, false, len(model.VMAttributes) > 0)
			if err != nil {
				return err
			}
			computeProfile.BaseVirtualMachineProfile = pointer.From(baseVirtualMachineProfileValue)
			properties.Properties.ComputeProfile = computeProfile

			if err := client.CreateOrUpdateThenPoll(ctx, id, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r ComputeFleetResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ComputeFleet.ComputeFleetClient

			id, err := fleets.ParseFleetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ComputeFleetResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			properties := existing.Model

			// API requires `osProfile.adminPassword` when updating resource but the GET API does not return the sensitive data `osProfile.adminPassword`
			if props := properties.Properties; props != nil {
				if len(model.VirtualMachineProfile[0].OsProfile[0].LinuxConfiguration) > 0 {
					if v := props.ComputeProfile.BaseVirtualMachineProfile.OsProfile; v != nil {
						v.AdminPassword = pointer.To(model.VirtualMachineProfile[0].OsProfile[0].LinuxConfiguration[0].AdminPassword)
					}
				}
				if len(model.VirtualMachineProfile[0].OsProfile[0].WindowsConfiguration) > 0 {
					if v := props.ComputeProfile.BaseVirtualMachineProfile.OsProfile; v != nil {
						v.AdminPassword = pointer.To(model.VirtualMachineProfile[0].OsProfile[0].WindowsConfiguration[0].AdminPassword)
					}
				}

				if a := model.AdditionalLocationProfile; len(a) > 0 && len(a[0].VirtualMachineProfileOverride) > 0 {
					if os := a[0].VirtualMachineProfileOverride[0].OsProfile; len(os) > 0 && len(os[0].WindowsConfiguration) > 0 {
						if v := props.AdditionalLocationsProfile; v != nil && len(v.LocationProfiles) > 0 {
							if vmss := v.LocationProfiles[0].VirtualMachineProfileOverride; vmss != nil && vmss.OsProfile != nil {
								vmss.OsProfile.AdminPassword = pointer.To(os[0].WindowsConfiguration[0].AdminPassword)
							}
						}
					}
					if os := a[0].VirtualMachineProfileOverride[0].OsProfile; len(os) > 0 && len(os[0].LinuxConfiguration) > 0 {
						if v := props.AdditionalLocationsProfile; v != nil && len(v.LocationProfiles) > 0 {
							if vmss := v.LocationProfiles[0].VirtualMachineProfileOverride; vmss != nil && vmss.OsProfile != nil {
								vmss.OsProfile.AdminPassword = pointer.To(os[0].LinuxConfiguration[0].AdminPassword)
							}
						}
					}
				}
			}

			if metadata.ResourceData.HasChange("identity") {
				identityValue, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(model.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				properties.Identity = identityValue
			}

			if metadata.ResourceData.HasChange("plan") {
				properties.Plan = expandPlanModel(model.Plan)
			}

			if metadata.ResourceData.HasChange("additional_location_profile") {
				additionalLocationsProfileValue, err := expandAdditionalLocationProfileModel(model.AdditionalLocationProfile, metadata.ResourceData)
				if err != nil {
					return err
				}
				properties.Properties.AdditionalLocationsProfile = additionalLocationsProfileValue
			}

			if metadata.ResourceData.HasChange("virtual_machine_profile") {
				baseVirtualMachineProfileValue, err := expandVirtualMachineProfileModel(model.VirtualMachineProfile, metadata.ResourceData, false, len(model.VMAttributes) > 0)
				if err != nil {
					return err
				}
				properties.Properties.ComputeProfile.BaseVirtualMachineProfile = pointer.From(baseVirtualMachineProfileValue)
			}

			if metadata.ResourceData.HasChange("compute_api_version") {
				properties.Properties.ComputeProfile.ComputeApiVersion = pointer.To(model.ComputeApiVersion)
			}

			if metadata.ResourceData.HasChange("regular_priority_profile") {
				properties.Properties.RegularPriorityProfile = expandRegularPriorityProfileModel(model.RegularPriorityProfile)
			}

			if metadata.ResourceData.HasChange("spot_priority_profile") {
				properties.Properties.SpotPriorityProfile = expandSpotPriorityProfileModel(model.SpotPriorityProfile)
			}

			if metadata.ResourceData.HasChange("vm_attributes") {
				properties.Properties.VMAttributes = expandVMAttributesModel(model.VMAttributes)
			}

			if metadata.ResourceData.HasChange("vm_sizes_profile") {
				properties.Properties.VMSizesProfile = pointer.From(expandVMSizeProfileModel(model.VMSizesProfile))
			}

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = &model.Tags
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ComputeFleetResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ComputeFleet.ComputeFleetClient

			id, err := fleets.ParseFleetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := ComputeFleetResourceModel{
				Name:              id.FleetName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)

				if model.Identity != nil {
					v, err := identity.FlattenSystemAndUserAssignedMapToModel(pointer.To(identity.SystemAndUserAssignedMap(*model.Identity)))
					if err != nil {
						return err
					}
					state.Identity = pointer.From(v)
				}

				state.Plan = flattenPlanModel(model.Plan)

				if props := model.Properties; props != nil {
					additionalLocationsProfileValue, err := flattenAdditionalLocationProfileModel(props.AdditionalLocationsProfile, metadata)
					if err != nil {
						return err
					}
					state.AdditionalLocationProfile = additionalLocationsProfileValue

					if v := props.ComputeProfile.AdditionalVirtualMachineCapabilities; v != nil {
						state.AdditionalCapabilitiesHibernationEnabled = pointer.From(v.HibernationEnabled)
						state.AdditionalCapabilitiesUltraSSDEnabled = pointer.From(v.UltraSSDEnabled)
					}

					baseVirtualMachineProfileValue, err := flattenVirtualMachineProfileModel(&props.ComputeProfile.BaseVirtualMachineProfile, metadata, false)
					if err != nil {
						return err
					}
					state.VirtualMachineProfile = baseVirtualMachineProfileValue

					state.ComputeApiVersion = pointer.From(props.ComputeProfile.ComputeApiVersion)
					state.PlatformFaultDomainCount = pointer.From(props.ComputeProfile.PlatformFaultDomainCount)

					state.RegularPriorityProfile = flattenRegularPriorityProfileModel(props.RegularPriorityProfile)
					state.SpotPriorityProfile = flattenSpotPriorityProfileModel(props.SpotPriorityProfile)
					state.UniqueId = pointer.From(props.UniqueId)
					state.VMAttributes = flattenVMAttributesModel(props.VMAttributes)
					state.VMSizesProfile = flattenVMSizeProfileModel(&props.VMSizesProfile)
				}
				state.Tags = pointer.From(model.Tags)
				state.Zones = pointer.From(model.Zones)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ComputeFleetResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ComputeFleet.ComputeFleetClient

			id, err := fleets.ParseFleetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func customizeForVmAttributes(state ComputeFleetResourceModel, d sdk.ResourceMetaData) error {
	// Once the properties of vm_attributes are added, they can not be removed
	old, new := d.ResourceDiff.GetChange("vm_attributes.0.accelerator_count")
	if len(old.([]interface{})) > 0 && len(new.([]interface{})) == 0 {
		if err := d.ResourceDiff.ForceNew("vm_attributes.0.accelerator_count"); err != nil {
			return err
		}
	}
	old, new = d.ResourceDiff.GetChange("vm_attributes.0.accelerator_manufacturers")
	if len(old.(*schema.Set).List()) > 0 && len(new.(*schema.Set).List()) == 0 {
		if err := d.ResourceDiff.ForceNew("vm_attributes.0.accelerator_manufacturers"); err != nil {
			return err
		}
	}
	old, new = d.ResourceDiff.GetChange("vm_attributes.0.accelerator_support")
	if old.(string) != "" && new.(string) == "" {
		if err := d.ResourceDiff.ForceNew("vm_attributes.0.accelerator_support"); err != nil {
			return err
		}
	}
	old, new = d.ResourceDiff.GetChange("vm_attributes.0.accelerator_types")
	if len(old.(*schema.Set).List()) > 0 && len(new.(*schema.Set).List()) == 0 {
		if err := d.ResourceDiff.ForceNew("vm_attributes.0.accelerator_types"); err != nil {
			return err
		}
	}
	old, new = d.ResourceDiff.GetChange("vm_attributes.0.architecture_types")
	if len(old.(*schema.Set).List()) > 0 && len(new.(*schema.Set).List()) == 0 {
		if err := d.ResourceDiff.ForceNew("vm_attributes.0.architecture_types"); err != nil {
			return err
		}
	}
	old, new = d.ResourceDiff.GetChange("vm_attributes.0.burstable_support")
	if old.(string) != "" && new.(string) == "" {
		if err := d.ResourceDiff.ForceNew("vm_attributes.0.burstable_support"); err != nil {
			return err
		}
	}
	old, new = d.ResourceDiff.GetChange("vm_attributes.0.cpu_manufacturers")
	if len(old.(*schema.Set).List()) > 0 && len(new.(*schema.Set).List()) == 0 {
		if err := d.ResourceDiff.ForceNew("vm_attributes.0.cpu_manufacturers"); err != nil {
			return err
		}
	}
	old, new = d.ResourceDiff.GetChange("vm_attributes.0.data_disk_count")
	if len(old.([]interface{})) > 0 && len(new.([]interface{})) == 0 {
		if err := d.ResourceDiff.ForceNew("vm_attributes.0.data_disk_count"); err != nil {
			return err
		}
	}
	old, new = d.ResourceDiff.GetChange("vm_attributes.0.excluded_vm_sizes")
	if len(old.(*pluginsdk.Set).List()) > 0 && len(new.(*pluginsdk.Set).List()) == 0 {
		if err := d.ResourceDiff.ForceNew("vm_attributes.0.excluded_vm_sizes"); err != nil {
			return err
		}
	}
	old, new = d.ResourceDiff.GetChange("vm_attributes.0.local_storage_disk_types")
	if len(old.(*schema.Set).List()) > 0 && len(new.(*schema.Set).List()) == 0 {
		if err := d.ResourceDiff.ForceNew("vm_attributes.0.local_storage_disk_types"); err != nil {
			return err
		}
	}
	old, new = d.ResourceDiff.GetChange("vm_attributes.0.local_storage_in_gib")
	if len(old.([]interface{})) > 0 && len(new.([]interface{})) == 0 {
		if err := d.ResourceDiff.ForceNew("vm_attributes.0.local_storage_in_gib"); err != nil {
			return err
		}
	}
	old, new = d.ResourceDiff.GetChange("vm_attributes.0.local_storage_support")
	if old.(string) != "" && new.(string) == "" {
		if err := d.ResourceDiff.ForceNew("vm_attributes.0.local_storage_support"); err != nil {
			return err
		}
	}
	old, new = d.ResourceDiff.GetChange("vm_attributes.0.memory_in_gib_per_vcpu")
	if len(old.([]interface{})) > 0 && len(new.([]interface{})) == 0 {
		if err := d.ResourceDiff.ForceNew("vm_attributes.0.memory_in_gib_per_vcpu"); err != nil {
			return err
		}
	}
	old, new = d.ResourceDiff.GetChange("vm_attributes.0.network_bandwidth_in_mbps")
	if len(old.([]interface{})) > 0 && len(new.([]interface{})) == 0 {
		if err := d.ResourceDiff.ForceNew("vm_attributes.0.network_bandwidth_in_mbps"); err != nil {
			return err
		}
	}
	old, new = d.ResourceDiff.GetChange("vm_attributes.0.network_interface_count")
	if len(old.([]interface{})) > 0 && len(new.([]interface{})) == 0 {
		if err := d.ResourceDiff.ForceNew("vm_attributes.0.network_interface_count"); err != nil {
			return err
		}
	}
	old, new = d.ResourceDiff.GetChange("vm_attributes.0.rdma_network_interface_count")
	if len(old.([]interface{})) > 0 && len(new.([]interface{})) == 0 {
		if err := d.ResourceDiff.ForceNew("vm_attributes.0.rdma_network_interface_count"); err != nil {
			return err
		}
	}
	old, new = d.ResourceDiff.GetChange("vm_attributes.0.rdma_support")
	if old.(string) != "" && new.(string) == "" {
		if err := d.ResourceDiff.ForceNew("vm_attributes.0.rdma_support"); err != nil {
			return err
		}
	}
	old, new = d.ResourceDiff.GetChange("vm_attributes.0.vm_categories")
	if len(old.(*schema.Set).List()) > 0 && len(new.(*schema.Set).List()) == 0 {
		if err := d.ResourceDiff.ForceNew("vm_attributes.0.vm_categories"); err != nil {
			return err
		}
	}

	if v := state.VMAttributes; len(v) > 0 {
		if v[0].AcceleratorSupport == string(fleets.VMAttributeSupportExcluded) {
			if len(v[0].AcceleratorManufacturers) > 0 {
				return fmt.Errorf("`accelerator_manufacturers` cannot be used when `accelerator_support` is specified as `Excluded`")
			}
			if len(v[0].AcceleratorTypes) > 0 {
				return fmt.Errorf("`accelerator_types` cannot be used when `accelerator_support` is specified as `Excluded`")
			}
			if len(v[0].AcceleratorCount) > 0 {
				return fmt.Errorf("`accelerator_count` cannot be used when `accelerator_support` is specified as `Excluded`")
			}
			if len(v[0].VMCategories) > 0 && utils.SliceContainsValue(v[0].VMCategories, string(fleets.VMCategoryGpuAccelerated)) {
				return fmt.Errorf("`GpuAccelerated` cannot be used when `accelerator_support` is specified as `Excluded`")
			}
		}

		if v[0].LocalStorageSupport == string(fleets.VMAttributeSupportExcluded) {
			if len(v[0].LocalStorageInGib) > 0 {
				return fmt.Errorf("`local_storage_in_gib` cannot be used when `local_storage_support` is specified as `Excluded`")
			}
			if len(v[0].LocalStorageDiskTypes) > 0 {
				return fmt.Errorf("`local_storage_disk_types` cannot be used when `local_storage_support` is specified as `Excluded`")
			}
		}

		if v[0].RdmaSupport == string(fleets.VMAttributeSupportExcluded) {
			if len(v[0].RdmaNetworkInterfaceCount) > 0 {
				return fmt.Errorf("`rdma_network_interface_count` cannot be used when `rdma_support` is specified as `Excluded`")
			}
		}
	}

	return nil
}

func (r ComputeFleetResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var state ComputeFleetResourceModel
			if err := metadata.DecodeDiff(&state); err != nil {
				return fmt.Errorf("DecodeDiff: %+v", err)
			}

			err := customizeForVmAttributes(state, metadata)
			if err != nil {
				return err
			}

			if v := state.Identity; len(v) > 0 {
				// whilst the Swagger defines multiple at this time only UAI is supported
				if v[0].Type != identity.TypeUserAssigned {
					return fmt.Errorf("the type `%s` of `identity` is not supported currently", v[0].Type)
				}
			}

			if len(state.SpotPriorityProfile) > 0 && len(state.RegularPriorityProfile) > 0 {
				if state.SpotPriorityProfile[0].Capacity+state.RegularPriorityProfile[0].Capacity > 10000 {
					return fmt.Errorf("the sum of `spot_priority_profile.0.capacity` and `regular_priority_profile.0.capacity` must be between `0` and `10000`, inclusive")
				}
			}

			if len(state.SpotPriorityProfile) > 0 {
				if state.SpotPriorityProfile[0].MaintainEnabled {
					if state.SpotPriorityProfile[0].MinCapacity > 0 {
						return fmt.Errorf("`spot_priority_profile.0.min_capacity` is unable to be specified if `spot_priority_profile.0.maintain_enabled` is enabled")
					}

					if len(state.VMSizesProfile) < 3 {
						return fmt.Errorf("`vm_sizes_profile` must be at least 3 Vm sizes if `spot_priority_profile.0.maintain_enabled` is enabled")
					}

					if len(state.Zones) == 0 {
						return fmt.Errorf("enabling `spot_priority_profile.0.maintain_enabled` requires all qualified availability zones in the region to be supported")
					}
				} else if len(state.RegularPriorityProfile) == 0 {
					// For Spot VMs, you may delete or replace existing VM sizes in your Compute Fleet configuration, if the capacity preference is set to Maintain capacity.
					// In all other scenarios requiring a modification to the running Compute Fleet, you may have to delete the existing Compute Fleet and create a new one.
					if metadata.ResourceDiff.HasChange("vm_sizes_profile") {
						if err := metadata.ResourceDiff.ForceNew("vm_sizes_profile"); err != nil {
							return err
						}
					}
				}

				if state.SpotPriorityProfile[0].MinCapacity > state.SpotPriorityProfile[0].Capacity {
					return fmt.Errorf("`spot_priority_profile.0.min_capacity` must be between `0` and `spot_priority_profile.0.capacity`, inclusive")
				}
			}

			if len(state.RegularPriorityProfile) > 0 {
				if state.RegularPriorityProfile[0].MinCapacity > state.RegularPriorityProfile[0].Capacity {
					return fmt.Errorf("`RegularPriorityProfile.0.min_capacity` must be between `0` and `RegularPriorityProfile.0.capacity`, inclusive")
				}
			}

			if len(state.VMSizesProfile) > 15 {
				return fmt.Errorf("the VM sizes count of `vm_sizes_profile` cannot be greater than `15`")
			}

			if v := state.VirtualMachineProfile[0].DataDisks; len(v) > 0 {
				storageAccountType := v[0].StorageAccountType
				ultraSSDEnabled := state.AdditionalCapabilitiesUltraSSDEnabled

				if !ultraSSDEnabled && storageAccountType == string(fleets.StorageAccountTypesUltraSSDLRS) {
					return fmt.Errorf("`UltraSSD_LRS` storage account type can be used only when `additional_capabilities_ultra_ssd_enabled` is enalbed")
				}

				if v[0].CreateOption == string(fleets.DiskCreateOptionTypesEmpty) {
					if v[0].DiskSizeInGB == 0 {
						return fmt.Errorf("`disk_size_in_gb` is required when`create_option` is `Empty`")
					}

					lunExist := metadata.ResourceDiff.GetRawConfig().AsValueMap()["virtual_machine_profile"].AsValueSlice()[0].AsValueMap()["data_disk"].AsValueSlice()[0].AsValueMap()["lun"]
					if lunExist.IsNull() {
						return fmt.Errorf("`lun` is required when`create_option` is `Empty`")
					}
				}
			}

			vmProfile := state.VirtualMachineProfile[0]
			osProfile := vmProfile.OsProfile[0]
			if len(osProfile.WindowsConfiguration) > 0 && len(osProfile.LinuxConfiguration) > 0 ||
				len(osProfile.WindowsConfiguration) == 0 && len(osProfile.LinuxConfiguration) == 0 {
				return fmt.Errorf("only one of `linux_configuration` and `windows_configuration` in `virtual_machine_profile` must be specified")
			}
			if vmProfile.SourceImageId != "" && len(vmProfile.SourceImageReference) > 0 {
				return fmt.Errorf("only one of `source_image_id` and `source_image_reference` in `virtual_machine_profile` must be specified")
			}

			if v := state.AdditionalLocationProfile; len(v) > 0 && len(v[0].VirtualMachineProfileOverride) > 0 {
				osProfile = state.AdditionalLocationProfile[0].VirtualMachineProfileOverride[0].OsProfile[0]
				if len(osProfile.WindowsConfiguration) > 0 && len(osProfile.LinuxConfiguration) > 0 ||
					len(osProfile.WindowsConfiguration) == 0 && len(osProfile.LinuxConfiguration) == 0 {
					return fmt.Errorf("only one of `linux_configuration` and `windows_configuration` in `virtual_machine_profile_override` must be specified")
				}

				if state.VirtualMachineProfile[0].CapacityReservationGroupId != "" && state.AdditionalLocationProfile[0].VirtualMachineProfileOverride[0].CapacityReservationGroupId == "" {
					return fmt.Errorf("`virtual_machine_profile_override.0.capacity_reservation_group_id` is required when `virtual_machine_profile.0.capacity_reservation_group_id` is specified")
				}

				dataDisks := state.AdditionalLocationProfile[0].VirtualMachineProfileOverride[0].DataDisks
				if len(dataDisks) > 0 {
					if dataDisks[0].CreateOption == string(fleets.DiskCreateOptionTypesEmpty) {
						if dataDisks[0].DiskSizeInGB == 0 {
							return fmt.Errorf("`disk_size_in_gb` is required when`create_option` is `Empty`")
						}
						lunExist := metadata.ResourceDiff.GetRawConfig().AsValueMap()["additional_location_profile"].AsValueSlice()[0].AsValueMap()["virtual_machine_profile_override"].AsValueSlice()[0].AsValueMap()["data_disk"].AsValueSlice()[0].AsValueMap()["lun"]
						if lunExist.IsNull() {
							return fmt.Errorf("`lun` is required when`create_option` is `Empty`")
						}
					}
				}
			}

			err = validateSecuritySetting(state.VirtualMachineProfile)
			if err != nil {
				return err
			}

			err = validateWindowsSetting(state.VirtualMachineProfile, metadata.ResourceDiff, false)
			if err != nil {
				return err
			}

			err = validateLinuxSetting(state.VirtualMachineProfile, metadata.ResourceDiff, false)
			if err != nil {
				return err
			}

			if len(state.AdditionalLocationProfile) > 0 {
				err = validateWindowsSetting(state.AdditionalLocationProfile[0].VirtualMachineProfileOverride, metadata.ResourceDiff, true)
				if err != nil {
					return err
				}

				err = validateLinuxSetting(state.AdditionalLocationProfile[0].VirtualMachineProfileOverride, metadata.ResourceDiff, true)
				if err != nil {
					return err
				}
			}
			for _, v := range state.VirtualMachineProfile[0].Extension {
				if v.ProtectedSettingsJson != "" {
					if len(v.ProtectedSettingsFromKeyVault) > 0 {
						return fmt.Errorf("`protected_settings_from_key_vault` cannot be used with `protected_settings_json`")
					}
				}
			}

			if len(state.Zones) > 0 && state.PlatformFaultDomainCount > 1 {
				return fmt.Errorf("specifying `zones` is not allowed when `platform_fault_domain_count` higher than 1")
			}

			return nil
		},
	}
}

func expandPlanModel(inputList []PlanModel) *fleets.Plan {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := fleets.Plan{
		Name:      input.Name,
		Product:   input.Product,
		Publisher: input.Publisher,
	}

	if input.PromotionCode != "" {
		output.PromotionCode = pointer.To(input.PromotionCode)
	}

	return &output
}

func expandRegularPriorityProfileModel(inputList []RegularPriorityProfileModel) *fleets.RegularPriorityProfile {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.RegularPriorityProfile{
		AllocationStrategy: pointer.To(fleets.RegularPriorityAllocationStrategy(input.AllocationStrategy)),
		Capacity:           pointer.To(input.Capacity),
		MinCapacity:        pointer.To(input.MinCapacity),
	}

	return &output
}

func expandSpotPriorityProfileModel(inputList []SpotPriorityProfileModel) *fleets.SpotPriorityProfile {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := fleets.SpotPriorityProfile{
		AllocationStrategy: pointer.To(fleets.SpotAllocationStrategy(input.AllocationStrategy)),
		Capacity:           pointer.To(input.Capacity),
		EvictionPolicy:     pointer.To(fleets.EvictionPolicy(input.EvictionPolicy)),
		Maintain:           pointer.To(input.MaintainEnabled),
		MinCapacity:        pointer.To(input.MinCapacity),
	}

	if input.MaxHourlyPricePerVM > 0 {
		output.MaxPricePerVM = pointer.To(input.MaxHourlyPricePerVM)
	}
	return &output
}

func expandVMAttributesModel(inputList []VMAttributesModel) *fleets.VMAttributes {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.VMAttributes{
		AcceleratorCount:          expandVMAttributeMinMaxIntegerModel(input.AcceleratorCount),
		AcceleratorManufacturers:  expandAcceleratorManufacturers(input.AcceleratorManufacturers),
		AcceleratorSupport:        pointer.To(fleets.VMAttributeSupport(input.AcceleratorSupport)),
		AcceleratorTypes:          expandAcceleratorTypes(input.AcceleratorTypes),
		ArchitectureTypes:         expandArchitectureTypes(input.ArchitectureTypes),
		BurstableSupport:          pointer.To(fleets.VMAttributeSupport(input.BurstableSupport)),
		CpuManufacturers:          expandCPUManufacturers(input.CpuManufacturers),
		DataDiskCount:             expandVMAttributeMinMaxIntegerModel(input.DataDiskCount),
		LocalStorageDiskTypes:     expandLocalStorageDiskTypes(input.LocalStorageDiskTypes),
		LocalStorageInGiB:         expandVMAttributeMinMaxDoubleModel(input.LocalStorageInGib),
		LocalStorageSupport:       pointer.To(fleets.VMAttributeSupport(input.LocalStorageSupport)),
		MemoryInGiBPerVCPU:        expandVMAttributeMinMaxDoubleModel(input.MemoryInGibPerVCPU),
		NetworkBandwidthInMbps:    expandVMAttributeMinMaxDoubleModel(input.NetworkBandwidthInMbps),
		NetworkInterfaceCount:     expandVMAttributeMinMaxIntegerModel(input.NetworkInterfaceCount),
		RdmaNetworkInterfaceCount: expandVMAttributeMinMaxIntegerModel(input.RdmaNetworkInterfaceCount),
		RdmaSupport:               pointer.To(fleets.VMAttributeSupport(input.RdmaSupport)),
		VMCategories:              expandVMCategories(input.VMCategories),
	}

	if len(input.ExcludedVMSizes) > 0 {
		output.ExcludedVMSizes = pointer.To(input.ExcludedVMSizes)
	}
	output.MemoryInGiB = pointer.From(expandVMAttributeMinMaxDoubleModel(input.MemoryInGib))

	output.VCPUCount = pointer.From(expandVMAttributeMinMaxIntegerModel(input.VCPUCount))

	return &output
}

func expandVMAttributeMinMaxIntegerModel(inputList []VMAttributeMinMaxIntegerModel) *fleets.VMAttributeMinMaxInteger {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.VMAttributeMinMaxInteger{
		Max: pointer.To(input.Max),
		Min: pointer.To(input.Min),
	}

	return &output
}

func expandVMAttributeMinMaxDoubleModel(inputList []VMAttributeMinMaxDoubleModel) *fleets.VMAttributeMinMaxDouble {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := fleets.VMAttributeMinMaxDouble{
		Max: &input.Max,
		Min: &input.Min,
	}

	return &output
}

func expandVMCategories(inputList []string) *[]fleets.VMCategory {
	if len(inputList) == 0 {
		return nil
	}

	outputList := make([]fleets.VMCategory, 0)
	for _, v := range inputList {
		if v != "" {
			outputList = append(outputList, fleets.VMCategory(v))
		}
	}
	return &outputList
}

func expandLocalStorageDiskTypes(inputList []string) *[]fleets.LocalStorageDiskType {
	if len(inputList) == 0 {
		return nil
	}

	outputList := make([]fleets.LocalStorageDiskType, 0)
	for _, v := range inputList {
		if v != "" {
			outputList = append(outputList, fleets.LocalStorageDiskType(v))
		}
	}
	return &outputList
}

func expandAcceleratorManufacturers(inputList []string) *[]fleets.AcceleratorManufacturer {
	if len(inputList) == 0 {
		return nil
	}

	result := make([]fleets.AcceleratorManufacturer, 0)

	for _, v := range inputList {
		if v != "" {
			result = append(result, fleets.AcceleratorManufacturer(v))
		}
	}

	return &result
}

func expandAcceleratorTypes(inputList []string) *[]fleets.AcceleratorType {
	if len(inputList) == 0 {
		return nil
	}

	outputList := make([]fleets.AcceleratorType, 0)
	for _, v := range inputList {
		if v != "" {
			outputList = append(outputList, fleets.AcceleratorType(v))
		}
	}
	return &outputList
}

func expandArchitectureTypes(inputList []string) *[]fleets.ArchitectureType {
	if len(inputList) == 0 {
		return nil
	}
	outputList := make([]fleets.ArchitectureType, 0)
	for _, v := range inputList {
		if v != "" {
			outputList = append(outputList, fleets.ArchitectureType(v))
		}
	}
	return &outputList
}

func expandCPUManufacturers(inputList []string) *[]fleets.CPUManufacturer {
	if len(inputList) == 0 {
		return nil
	}

	outputList := make([]fleets.CPUManufacturer, 0)
	for _, v := range inputList {
		if v != "" {
			outputList = append(outputList, fleets.CPUManufacturer(v))
		}
	}
	return &outputList
}

func expandVMSizeProfileModel(inputList []VMSizeProfileModel) *[]fleets.VMSizeProfile {
	if len(inputList) == 0 {
		return nil
	}

	outputList := make([]fleets.VMSizeProfile, 0)
	for _, v := range inputList {
		input := v
		output := fleets.VMSizeProfile{
			Name: input.Name,
		}
		if input.Rank > 0 {
			output.Rank = pointer.To(input.Rank)
		}
		outputList = append(outputList, output)
	}
	return &outputList
}

func expandAdditionalLocationProfileModel(inputList []AdditionalLocationProfileModel, d *schema.ResourceData) (*fleets.AdditionalLocationsProfile, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	output := fleets.AdditionalLocationsProfile{}
	outputList := make([]fleets.LocationProfile, 0)
	for _, v := range inputList {
		input := v
		output := fleets.LocationProfile{
			Location: input.Location,
		}

		virtualMachineProfileOverrideValue, err := expandVirtualMachineProfileModel(input.VirtualMachineProfileOverride, d, true, true)
		if err != nil {
			return nil, err
		}

		output.VirtualMachineProfileOverride = virtualMachineProfileOverrideValue

		outputList = append(outputList, output)
	}

	output.LocationProfiles = outputList

	return &output, nil
}

func flattenPlanModel(input *fleets.Plan) []PlanModel {
	outputList := make([]PlanModel, 0)
	if input == nil {
		return outputList
	}
	output := PlanModel{
		Name:      input.Name,
		Product:   input.Product,
		Publisher: input.Publisher,
	}

	output.PromotionCode = pointer.From(input.PromotionCode)

	return append(outputList, output)
}

func flattenAdditionalLocationProfileModel(input *fleets.AdditionalLocationsProfile, metadata sdk.ResourceMetaData) ([]AdditionalLocationProfileModel, error) {
	outputList := make([]AdditionalLocationProfileModel, 0)
	if input == nil {
		return outputList, nil
	}

	for _, input := range input.LocationProfiles {
		output := AdditionalLocationProfileModel{
			Location: input.Location,
		}
		virtualMachineProfileOverrideValue, err := flattenVirtualMachineProfileModel(input.VirtualMachineProfileOverride, metadata, true)
		if err != nil {
			return nil, err
		}

		output.VirtualMachineProfileOverride = virtualMachineProfileOverrideValue
		outputList = append(outputList, output)
	}

	return outputList, nil
}

func flattenSubResourceId(inputList []fleets.SubResource) []string {
	outputList := make([]string, 0)
	if len(inputList) == 0 {
		return outputList
	}
	for _, input := range inputList {
		output := pointer.From(input.Id)
		outputList = append(outputList, output)
	}
	return outputList
}

func flattenIPTagModel(inputList *[]fleets.VirtualMachineScaleSetIPTag) []IPTagModel {
	outputList := make([]IPTagModel, 0)
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		output := IPTagModel{}

		output.Type = pointer.From(input.IPTagType)
		output.Tag = pointer.From(input.Tag)
		outputList = append(outputList, output)
	}
	return outputList
}

func flattenRegularPriorityProfileModel(input *fleets.RegularPriorityProfile) []RegularPriorityProfileModel {
	outputList := make([]RegularPriorityProfileModel, 0)
	if input == nil {
		return outputList
	}

	output := RegularPriorityProfileModel{}
	output.AllocationStrategy = string(pointer.From(input.AllocationStrategy))
	output.Capacity = pointer.From(input.Capacity)
	output.MinCapacity = pointer.From(input.MinCapacity)

	return append(outputList, output)
}

func flattenSpotPriorityProfileModel(input *fleets.SpotPriorityProfile) []SpotPriorityProfileModel {
	outputList := make([]SpotPriorityProfileModel, 0)
	if input == nil {
		return outputList
	}

	output := SpotPriorityProfileModel{}
	output.AllocationStrategy = string(pointer.From(input.AllocationStrategy))
	output.Capacity = pointer.From(input.Capacity)
	output.EvictionPolicy = string(pointer.From(input.EvictionPolicy))
	output.MaintainEnabled = pointer.From(input.Maintain)

	// defaulted since MaxHourlyPricePerVM isn't returned if it's unset
	maxHourlyPricePerVM := float64(-1.0)
	if input.MaxPricePerVM != nil {
		maxHourlyPricePerVM = pointer.From(input.MaxPricePerVM)
	}
	output.MaxHourlyPricePerVM = maxHourlyPricePerVM

	output.MinCapacity = pointer.From(input.MinCapacity)

	return append(outputList, output)
}

func flattenVMAttributesModel(input *fleets.VMAttributes) []VMAttributesModel {
	outputList := make([]VMAttributesModel, 0)
	if input == nil {
		return outputList
	}
	output := VMAttributesModel{
		AcceleratorCount:          flattenVMAttributeMinMaxIntegerModel(input.AcceleratorCount),
		AcceleratorManufacturers:  flattenToStringSlice(input.AcceleratorManufacturers),
		AcceleratorTypes:          flattenToStringSlice(input.AcceleratorTypes),
		ArchitectureTypes:         flattenToStringSlice(input.ArchitectureTypes),
		CpuManufacturers:          flattenToStringSlice(input.CpuManufacturers),
		DataDiskCount:             flattenVMAttributeMinMaxIntegerModel(input.DataDiskCount),
		LocalStorageDiskTypes:     flattenToStringSlice(input.LocalStorageDiskTypes),
		LocalStorageInGib:         flattenVMAttributeMinMaxDoubleModel(input.LocalStorageInGiB),
		MemoryInGib:               flattenVMAttributeMinMaxDoubleModel(&input.MemoryInGiB),
		MemoryInGibPerVCPU:        flattenVMAttributeMinMaxDoubleModel(input.MemoryInGiBPerVCPU),
		NetworkBandwidthInMbps:    flattenVMAttributeMinMaxDoubleModel(input.NetworkBandwidthInMbps),
		NetworkInterfaceCount:     flattenVMAttributeMinMaxIntegerModel(input.NetworkInterfaceCount),
		RdmaNetworkInterfaceCount: flattenVMAttributeMinMaxIntegerModel(input.RdmaNetworkInterfaceCount),
		VCPUCount:                 flattenVMAttributeMinMaxIntegerModel(&input.VCPUCount),
		VMCategories:              flattenToStringSlice(input.VMCategories),
	}

	output.AcceleratorSupport = string(pointer.From(input.AcceleratorSupport))
	output.BurstableSupport = string(pointer.From(input.BurstableSupport))
	output.ExcludedVMSizes = pointer.From(input.ExcludedVMSizes)
	output.LocalStorageSupport = string(pointer.From(input.LocalStorageSupport))
	output.RdmaSupport = string(pointer.From(input.RdmaSupport))

	return append(outputList, output)
}

func flattenVMAttributeMinMaxIntegerModel(input *fleets.VMAttributeMinMaxInteger) []VMAttributeMinMaxIntegerModel {
	outputList := make([]VMAttributeMinMaxIntegerModel, 0)
	if input == nil {
		return outputList
	}
	output := VMAttributeMinMaxIntegerModel{}
	output.Max = pointer.From(input.Max)
	output.Min = pointer.From(input.Min)

	return append(outputList, output)
}

func flattenVMAttributeMinMaxDoubleModel(input *fleets.VMAttributeMinMaxDouble) []VMAttributeMinMaxDoubleModel {
	outputList := make([]VMAttributeMinMaxDoubleModel, 0)
	if input == nil {
		return outputList
	}
	output := VMAttributeMinMaxDoubleModel{}
	output.Max = pointer.From(input.Max)
	output.Min = pointer.From(input.Min)

	return append(outputList, output)
}

func flattenToStringSlice[T any](inputList *[]T) []string {
	outputList := make([]string, 0)
	if inputList == nil {
		return outputList
	}

	result := make([]string, len(*inputList))
	for i, v := range *inputList {
		result[i] = fmt.Sprintf("%v", v)
	}

	return result
}

func flattenVMSizeProfileModel(inputList *[]fleets.VMSizeProfile) []VMSizeProfileModel {
	outputList := make([]VMSizeProfileModel, 0)
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		output := VMSizeProfileModel{
			Name: input.Name,
		}
		output.Rank = pointer.From(input.Rank)
		outputList = append(outputList, output)
	}
	return outputList
}
