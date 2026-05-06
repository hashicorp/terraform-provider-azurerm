// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package computefleet

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurefleet/2024-11-01/fleets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservationgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/images"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ComputeFleetResourceModel struct {
	Name                     string                                     `tfschema:"name"`
	ResourceGroupName        string                                     `tfschema:"resource_group_name"`
	Identity                 []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Location                 string                                     `tfschema:"location"`
	Plan                     []PlanModel                                `tfschema:"plan"`
	AdditionalCapabilities   []AdditionalCapabilitiesModel              `tfschema:"additional_capabilities"`
	VirtualMachineProfile    []VirtualMachineProfileModel               `tfschema:"virtual_machine_profile"`
	ComputeApiVersion        string                                     `tfschema:"compute_api_version"`
	PlatformFaultDomainCount int64                                      `tfschema:"platform_fault_domain_count"`
	OnDemandCapacity         []OnDemandCapacityModel                    `tfschema:"on_demand_capacity"`
	SpotCapacity             []SpotCapacityModel                        `tfschema:"spot_capacity"`
	UniqueId                 string                                     `tfschema:"unique_id"`
	VMSizesProfile           []VMSizeProfileModel                       `tfschema:"vm_sizes_profile"`
	Tags                     map[string]string                          `tfschema:"tags"`
	Zones                    []string                                   `tfschema:"zones"`
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
	AdminPassword                     string             `tfschema:"admin_password"`
	AdminUsername                     string             `tfschema:"admin_username"`
	ComputerNamePrefix                string             `tfschema:"computer_name_prefix"`
	Secret                            []LinuxSecretModel `tfschema:"secret"`
	PasswordAuthenticationEnabled     bool               `tfschema:"password_authentication_enabled"`
	VMAgentPlatformUpdatesEnabled     bool               `tfschema:"vm_agent_platform_updates_enabled"`
	PatchMode                         string             `tfschema:"patch_mode"`
	BypassPlatformSafetyChecksEnabled bool               `tfschema:"bypass_platform_safety_checks_enabled"`
	PatchRebooting                    string             `tfschema:"patch_rebooting"`
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
	PatchMode                         string                           `tfschema:"patch_mode"`
	BypassPlatformSafetyChecksEnabled bool                             `tfschema:"bypass_platform_safety_checks_enabled"`
	PatchRebooting                    string                           `tfschema:"patch_rebooting"`
	HotPatchingEnabled                bool                             `tfschema:"hot_patching_enabled"`
	ProvisionVMAgentEnabled           bool                             `tfschema:"provision_vm_agent_enabled"`
	TimeZone                          string                           `tfschema:"time_zone"`
	WinRM                             []WinRMModel                     `tfschema:"winrm_listener"`
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
	AllocationStrategy      string  `tfschema:"allocation_strategy"`
	TargetCapacity          int64   `tfschema:"target_capacity"`
	EvictionPolicy          string  `tfschema:"eviction_policy"`
	MaintainCapacityEnabled bool    `tfschema:"maintain_capacity_enabled"`
	MaxHourlyPricePerVM     float64 `tfschema:"max_hourly_price_per_vm"`
	MinimumCapacity         int64   `tfschema:"minimum_capacity"`
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
			// From Azure REST API, `name` can be up to 80 characters long. It must begin with a word character, and it must end with a word character or with '_'. The name may contain word characters or '.', '-', '_'. Instead, the validation is implemented according to Azure portal behavior
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9.-]{0,62}[a-zA-Z0-9]$|^[a-zA-Z0-9]$"),
				"The fleet name can only start or end with a number or a letter, and can contain only letters, numbers, periods (.), hyphens (-), up to 64 characters",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"virtual_machine_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"network_interface": networkInterfaceSchema(),

					"os_profile": osProfileSchema(),

					"boot_diagnostic_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  false,
					},

					"boot_diagnostic_storage_account_endpoint": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					},

					"capacity_reservation_group_id": commonschema.ResourceIDReferenceOptionalForceNew(&capacityreservationgroups.CapacityReservationGroupId{}),

					"data_disk": storageProfileDataDiskSchema(),

					"encryption_at_host_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  false,
					},

					"extension": extensionSchema(),

					"extension_operations_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  true,
					},

					"extensions_time_budget_duration": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: azValidate.ISO8601DurationBetween("PT15M", "PT2H"),
					},

					"gallery_application": galleryApplicationSchema(),

					"license_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						ValidateFunc: validation.StringInSlice([]string{
							"RHEL_BYOS",
							"SLES_BYOS",
							"Windows_Client",
							"Windows_Server",
						}, false),
					},

					"network_api_version": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						// `Default` set according to portal behavior
						Default: string(fleets.NetworkApiVersionTwoZeroTwoZeroNegativeOneOneNegativeZeroOne),
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile(`^\d{4}-\d{2}-\d{2}(-preview)?$`),
							"`network_api_version` must be in the format `YYYY-MM-DD` with an optional `-preview` suffix",
						),
					},

					"os_disk": storageProfileOsDiskSchema(),

					"scheduled_event_os_image_timeout_duration": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						ValidateFunc: validation.StringInSlice([]string{
							"PT15M",
						}, false),
					},

					"scheduled_event_termination_timeout_duration": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: azValidate.ISO8601DurationBetween("PT5M", "PT15M"),
					},

					"secure_boot_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  false,
					},

					"source_image_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						ValidateFunc: validation.Any(
							images.ValidateImageID,
							computeValidate.SharedImageID,
							computeValidate.SharedImageVersionID,
							computeValidate.CommunityGalleryImageID,
							computeValidate.CommunityGalleryImageVersionID,
							computeValidate.SharedGalleryImageID,
							computeValidate.SharedGalleryImageVersionID,
						),
						ExactlyOneOf: []string{
							"virtual_machine_profile.0.source_image_reference",
							"virtual_machine_profile.0.source_image_id",
						},
					},

					"source_image_reference": storageProfileSourceImageReferenceSchema(),

					"user_data_base64": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						ValidateFunc: validation.All(
							validation.StringIsBase64,
							validation.StringLenBetween(1, 349528),
						),
					},

					"vtpm_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  false,
					},
				},
			},
		},

		"vm_sizes_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			// In Azure portal, maximum number of items is 15 while in Azure REST API, it is 10
			// Leaving `MaxItems` as 15 while confirming with service team whether the limit can be fixed
			MaxItems: 15,
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
						ValidateFunc: validation.IntAtLeast(0),
					},
				},
			},
		},

		"additional_capabilities": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"hibernation_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  false,
						AtLeastOneOf: []string{
							"additional_capabilities.0.hibernation_enabled",
							"additional_capabilities.0.ultra_ssd_enabled",
						},
					},
					"ultra_ssd_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  false,
						AtLeastOneOf: []string{
							"additional_capabilities.0.hibernation_enabled",
							"additional_capabilities.0.ultra_ssd_enabled",
						},
					},
				},
			},
		},

		"compute_api_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// NOTE: O+C When creating new resource without providing this property, Azure will return a default value. As the default value is not documented anywhere, user experience will be poor to make this a `Required` property
			Computed: true,
			// This property is tested to be non-updatable
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^\d{4}-\d{2}-\d{2}(-preview)?$`),
				"`compute_api_version` must be in the format `YYYY-MM-DD` with an optional `-preview` suffix",
			),
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"on_demand_capacity": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MaxItems:     1,
			AtLeastOneOf: []string{"on_demand_capacity", "spot_capacity"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"target_capacity": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 10000),
					},

					"allocation_strategy": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(fleets.RegularPriorityAllocationStrategyLowestPrice),
						ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForRegularPriorityAllocationStrategy(), false),
					},

					"minimum_starting_capacity": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      0,
						ValidateFunc: validation.IntAtLeast(0),
					},
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
						Type:         pluginsdk.TypeString,
						ForceNew:     true,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"product": {
						Type:         pluginsdk.TypeString,
						ForceNew:     true,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"publisher": {
						Type:         pluginsdk.TypeString,
						ForceNew:     true,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"promotion_code": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"platform_fault_domain_count": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      1,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, 2),
		},

		"spot_capacity": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MaxItems:     1,
			AtLeastOneOf: []string{"on_demand_capacity", "spot_capacity"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					// `maintain_capacity_enabled` is `Required` as multiple properties depend on its value
					// `maintain_capacity_enabled` can be `true` using Azure REST API even if less than 3 `vm_sizes_profile` blocks are defined, which is different from Azure portal behavior
					// `vm_sizes_profile` can be updated in Azure portal even if `maintain_capacity_enabled` is `false`
					"maintain_capacity_enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},

					"target_capacity": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 10000),
					},

					"allocation_strategy": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(fleets.SpotAllocationStrategyPriceCapacityOptimized),
						ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForSpotAllocationStrategy(), false),
					},

					"eviction_policy": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(fleets.EvictionPolicyDelete),
						ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForEvictionPolicy(), false),
					},

					"max_hourly_price_per_vm": {
						Type:         pluginsdk.TypeFloat,
						Optional:     true,
						Default:      -1,
						ValidateFunc: computeValidate.SpotMaxPrice,
					},

					"minimum_capacity": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      0,
						ValidateFunc: validation.IntAtLeast(0),
					},
				},
			},
		},

		"zones": commonschema.ZonesMultipleOptionalForceNew(),

		"tags": commonschema.Tags(),
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
				Plan:     r.expandPlanModel(model.Plan),
				Properties: &fleets.FleetProperties{
					RegularPriorityProfile: r.expandOnDemandCapacityModel(model.OnDemandCapacity),
					SpotPriorityProfile:    r.expandSpotCapacityModel(model.SpotCapacity),
					VMSizesProfile:         pointer.From(r.expandVMSizeProfileModel(model.VMSizesProfile, metadata)),
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

			computeProfile := fleets.ComputeProfile{
				AdditionalVirtualMachineCapabilities: r.expandAdditionalCapabilities(model.AdditionalCapabilities),
				PlatformFaultDomainCount:             pointer.To(model.PlatformFaultDomainCount),
			}
			if model.ComputeApiVersion != "" {
				computeProfile.ComputeApiVersion = pointer.To(model.ComputeApiVersion)
			}

			baseVirtualMachineProfileValue, err := r.expandVirtualMachineProfileModel(model.VirtualMachineProfile, &model, metadata.ResourceData)
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
				if v := props.ComputeProfile.BaseVirtualMachineProfile.OsProfile; v != nil {
					if len(model.VirtualMachineProfile[0].OsProfile[0].LinuxConfiguration) > 0 {
						v.AdminPassword = pointer.To(model.VirtualMachineProfile[0].OsProfile[0].LinuxConfiguration[0].AdminPassword)
					}
					if len(model.VirtualMachineProfile[0].OsProfile[0].WindowsConfiguration) > 0 {
						v.AdminPassword = pointer.To(model.VirtualMachineProfile[0].OsProfile[0].WindowsConfiguration[0].AdminPassword)
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
				properties.Plan = r.expandPlanModel(model.Plan)
			}

			if metadata.ResourceData.HasChange("on_demand_capacity") {
				properties.Properties.RegularPriorityProfile = r.expandOnDemandCapacityModel(model.OnDemandCapacity)
			}

			if metadata.ResourceData.HasChange("spot_capacity") {
				properties.Properties.SpotPriorityProfile = r.expandSpotCapacityModel(model.SpotCapacity)
			}

			if metadata.ResourceData.HasChange("vm_sizes_profile") {
				properties.Properties.VMSizesProfile = pointer.From(r.expandVMSizeProfileModel(model.VMSizesProfile, metadata))
			}

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = pointer.To(model.Tags)
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

				state.Plan = r.flattenPlanModel(model.Plan)

				if props := model.Properties; props != nil {
					state.AdditionalCapabilities = r.flattenAdditionalCapabilities(props.ComputeProfile.AdditionalVirtualMachineCapabilities)

					baseVirtualMachineProfileValue, err := r.flattenVirtualMachineProfileModel(&props.ComputeProfile.BaseVirtualMachineProfile, metadata)
					if err != nil {
						return err
					}
					state.VirtualMachineProfile = baseVirtualMachineProfileValue

					state.ComputeApiVersion = pointer.From(props.ComputeProfile.ComputeApiVersion)
					state.PlatformFaultDomainCount = pointer.From(props.ComputeProfile.PlatformFaultDomainCount)
					state.OnDemandCapacity = r.flattenOnDemandCapacityModel(props.RegularPriorityProfile)
					state.SpotCapacity = r.flattenSpotCapacityModel(props.SpotPriorityProfile)
					state.UniqueId = pointer.From(props.UniqueId)
					state.VMSizesProfile = r.flattenVMSizeProfileModel(&props.VMSizesProfile)
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

func (r ComputeFleetResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var state ComputeFleetResourceModel
			if err := metadata.DecodeDiff(&state); err != nil {
				return fmt.Errorf("DecodeDiff: %+v", err)
			}

			if v := state.Identity; len(v) > 0 {
				// whilst the Swagger defines multiple at this time only UAI is supported
				if v[0].Type != identity.TypeUserAssigned {
					return fmt.Errorf("the type `%s` of `identity` is not supported currently", v[0].Type)
				}
			}

			if len(state.SpotCapacity) > 0 && len(state.OnDemandCapacity) > 0 {
				if state.SpotCapacity[0].TargetCapacity+state.OnDemandCapacity[0].TargetCapacity > 10000 {
					return errors.New("the sum of `spot_capacity.0.target_capacity` and `on_demand_capacity.0.target_capacity` must be between `0` and `10000`, inclusive")
				}
			}

			if len(state.SpotCapacity) > 0 {
				if state.SpotCapacity[0].MaintainCapacityEnabled {
					if state.SpotCapacity[0].MinimumCapacity > 0 {
						return errors.New("`spot_capacity.0.minimum_capacity` is unable to be specified if `spot_capacity.0.maintain_capacity_enabled` is enabled")
					}

					if len(state.VMSizesProfile) < 3 {
						return errors.New("`vm_sizes_profile` must be at least 3 Vm sizes if `spot_capacity.0.maintain_capacity_enabled` is enabled")
					}

					if len(state.Zones) == 0 {
						return errors.New("enabling `spot_capacity.0.maintain_capacity_enabled` requires all qualified availability zones in the region to be supported")
					}
				}

				if state.SpotCapacity[0].MinimumCapacity > state.SpotCapacity[0].TargetCapacity {
					return errors.New("`spot_capacity.0.minimum_capacity` must be between `0` and `spot_capacity.0.target_capacity`, inclusive")
				}
			}

			// Relax `ForceNew` behavior of `spot_capacity` properties to allow updating resource from without to with `spot_capacity` properties
			// Apart from this, other changes are not allowed
			forceNewSpotCapacityProperties := []string{
				"maintain_capacity_enabled",
				"allocation_strategy",
				"eviction_policy",
				"max_hourly_price_per_vm",
				"minimum_capacity",
			}
			oldSpotCapacity, _ := metadata.ResourceDiff.GetChange("spot_capacity")
			if len(oldSpotCapacity.([]interface{})) > 0 {
				if len(state.SpotCapacity) == 0 {
					metadata.ResourceDiff.ForceNew("spot_capacity")
				}

				for _, spotCapacityProperty := range forceNewSpotCapacityProperties {
					spotCapacityProperty = fmt.Sprintf("spot_capacity.0.%s", spotCapacityProperty)
					if metadata.ResourceDiff.HasChange(spotCapacityProperty) {
						metadata.ResourceDiff.ForceNew(spotCapacityProperty)
					}
				}
			}

			if len(state.OnDemandCapacity) > 0 {
				if state.OnDemandCapacity[0].MinimumStartingCapacity > state.OnDemandCapacity[0].TargetCapacity {
					return errors.New("`on_demand_capacity.0.minimum_starting_capacity` must be between `0` and `on_demand_capacity.0.target_capacity`, inclusive")
				}
			}

			// Relax `ForceNew` behavior of `on_demand_capacity` properties to allow updating resource from without to with `on_demand_capacity` properties
			// Apart from this, other changes are not allowed
			forceNewOnDemandCapacityProperties := []string{
				"allocation_strategy",
				"minimum_starting_capacity",
			}
			oldOnDemandCapacity, _ := metadata.ResourceDiff.GetChange("on_demand_capacity")
			if len(oldOnDemandCapacity.([]interface{})) > 0 {
				if len(state.OnDemandCapacity) == 0 {
					metadata.ResourceDiff.ForceNew("on_demand_capacity")
				}

				for _, onDemandCapacityProperty := range forceNewOnDemandCapacityProperties {
					onDemandCapacityProperty = fmt.Sprintf("on_demand_capacity.0.%s", onDemandCapacityProperty)
					if metadata.ResourceDiff.HasChange(onDemandCapacityProperty) {
						metadata.ResourceDiff.ForceNew(onDemandCapacityProperty)
					}
				}
			}

			if v := state.VirtualMachineProfile[0].DataDisks; len(v) > 0 {
				ultraSSDEnabled := false
				if ac := state.AdditionalCapabilities; len(ac) > 0 {
					ultraSSDEnabled = ac[0].UltraSsdEnabled
				}

				for i, dataDisk := range v {
					storageAccountType := dataDisk.StorageAccountType
					if !ultraSSDEnabled && storageAccountType == string(fleets.StorageAccountTypesUltraSSDLRS) {
						return errors.New("`UltraSSD_LRS` storage account type can be used only when `ultra_ssd_enabled` is enabled")
					}

					if dataDisk.CreateOption == string(fleets.DiskCreateOptionTypesEmpty) {
						if dataDisk.DiskSizeInGiB == 0 {
							return fmt.Errorf("`virtual_machine_profile.0.data_disk.%d.disk_size_in_gib` is required when `create_option` is `Empty`", i)
						}

						lunExist := metadata.ResourceDiff.GetRawConfig().AsValueMap()["virtual_machine_profile"].AsValueSlice()[0].AsValueMap()["data_disk"].AsValueSlice()[i].AsValueMap()["lun"]
						if lunExist.IsNull() {
							return fmt.Errorf("`virtual_machine_profile.0.data_disk.%d.lun` is required when `create_option` is `Empty`", i)
						}
					}

					if dataDisk.Caching != "" && dataDisk.DiskSizeInGiB > 4095 {
						return fmt.Errorf("`virtual_machine_profile.0.data_disk.%d.disk_size_in_gib` cannot be greater than 4095 GiB when `caching` is specified", i)
					}
				}
			}

			for i, vmSizesProfile := range state.VMSizesProfile {
				if !metadata.ResourceDiff.GetRawConfig().AsValueMap()["vm_sizes_profile"].AsValueSlice()[i].AsValueMap()["rank"].IsNull() && (len(state.OnDemandCapacity) == 0 || state.OnDemandCapacity[0].AllocationStrategy != string(fleets.RegularPriorityAllocationStrategyPrioritized)) {
					return errors.New("`on_demand_capacity` should be specified with `allocation_strategy` equal to `Prioritized` when `virtual_machine_profile.0.vm_sizes_profile.#.rank` is specified")
				}

				if vmSizesProfile.Rank >= int64(len(state.VMSizesProfile)) {
					return fmt.Errorf("`virtual_machine_profile.0.vm_sizes_profile.%d.rank` must be less than the number of `vm_sizes_profile` blocks", i)
				}
			}

			vmSizesProfileNames := make([]interface{}, len(state.VMSizesProfile))
			for i, v := range state.VMSizesProfile {
				vmSizesProfileNames[i] = v.Name
			}

			vmSizesProfileNamesWithoutDuplicate := schema.NewSet(pluginsdk.HashString, vmSizesProfileNames)
			if len(state.VMSizesProfile) > vmSizesProfileNamesWithoutDuplicate.Len() {
				return errors.New("`virtual_machine_profile.0.vm_sizes_profile` blocks with duplicated `name` are not allowed")
			}

			err := validateSecuritySetting(state.VirtualMachineProfile)
			if err != nil {
				return err
			}

			err = validateWindowsSetting(state.VirtualMachineProfile, metadata.ResourceDiff)
			if err != nil {
				return err
			}

			err = validateLinuxSetting(state.VirtualMachineProfile, metadata.ResourceDiff)
			if err != nil {
				return err
			}

			for _, v := range state.VirtualMachineProfile[0].Extension {
				if v.ProtectedSettingsJson != "" && len(v.ProtectedSettingsFromKeyVault) > 0 {
					return errors.New("`protected_settings_from_key_vault` cannot be used with `protected_settings_json`")
				}
			}

			if len(state.Zones) > 0 && state.PlatformFaultDomainCount > 1 {
				return errors.New("specifying `zones` is not allowed when `platform_fault_domain_count` higher than 1")
			}

			return nil
		},
	}
}

func (r ComputeFleetResource) expandPlanModel(inputList []PlanModel) *fleets.Plan {
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

func (r ComputeFleetResource) expandOnDemandCapacityModel(inputList []OnDemandCapacityModel) *fleets.RegularPriorityProfile {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := fleets.RegularPriorityProfile{
		AllocationStrategy: pointer.ToEnum[fleets.RegularPriorityAllocationStrategy](input.AllocationStrategy),
		Capacity:           pointer.To(input.TargetCapacity),
		MinCapacity:        pointer.To(input.MinimumStartingCapacity),
	}

	return &output
}

func (r ComputeFleetResource) expandSpotCapacityModel(inputList []SpotCapacityModel) *fleets.SpotPriorityProfile {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := fleets.SpotPriorityProfile{
		AllocationStrategy: pointer.ToEnum[fleets.SpotAllocationStrategy](input.AllocationStrategy),
		Capacity:           pointer.To(input.TargetCapacity),
		EvictionPolicy:     pointer.ToEnum[fleets.EvictionPolicy](input.EvictionPolicy),
		Maintain:           pointer.To(input.MaintainCapacityEnabled),
		MinCapacity:        pointer.To(input.MinimumCapacity),
	}

	if input.MaxHourlyPricePerVM > 0 || input.MaxHourlyPricePerVM == -1 {
		output.MaxPricePerVM = pointer.To(input.MaxHourlyPricePerVM)
	}
	return &output
}

func (r ComputeFleetResource) expandVMSizeProfileModel(inputList []VMSizeProfileModel, metadata sdk.ResourceMetaData) *[]fleets.VMSizeProfile {
	if len(inputList) == 0 {
		return nil
	}

	outputList := make([]fleets.VMSizeProfile, 0)
	for i, v := range inputList {
		input := v
		output := fleets.VMSizeProfile{
			Name: input.Name,
		}
		if !metadata.ResourceData.GetRawConfig().AsValueMap()["vm_sizes_profile"].AsValueSlice()[i].AsValueMap()["rank"].IsNull() {
			output.Rank = pointer.To(input.Rank)
		}
		outputList = append(outputList, output)
	}
	return &outputList
}

func (r ComputeFleetResource) expandAdditionalCapabilities(inputList []AdditionalCapabilitiesModel) *fleets.AdditionalCapabilities {
	if len(inputList) == 0 {
		return nil
	}

	capabilities := fleets.AdditionalCapabilities{
		UltraSSDEnabled:    pointer.To(inputList[0].UltraSsdEnabled),
		HibernationEnabled: pointer.To(inputList[0].HibernationEnabled),
	}

	return &capabilities
}

func (r ComputeFleetResource) flattenPlanModel(input *fleets.Plan) []PlanModel {
	outputList := make([]PlanModel, 0)
	if input == nil {
		return outputList
	}
	output := PlanModel{
		Name:          input.Name,
		Product:       input.Product,
		Publisher:     input.Publisher,
		PromotionCode: pointer.From(input.PromotionCode),
	}

	return append(outputList, output)
}

func (r ComputeFleetResource) flattenAdditionalCapabilities(input *fleets.AdditionalCapabilities) []AdditionalCapabilitiesModel {
	outputList := make([]AdditionalCapabilitiesModel, 0)
	if input == nil {
		return outputList
	}
	output := AdditionalCapabilitiesModel{
		UltraSsdEnabled:    pointer.From(input.UltraSSDEnabled),
		HibernationEnabled: pointer.From(input.HibernationEnabled),
	}

	return append(outputList, output)
}

func (r ComputeFleetResource) flattenOnDemandCapacityModel(input *fleets.RegularPriorityProfile) []OnDemandCapacityModel {
	outputList := make([]OnDemandCapacityModel, 0)
	if input == nil {
		return outputList
	}

	output := OnDemandCapacityModel{
		AllocationStrategy:      string(pointer.From(input.AllocationStrategy)),
		TargetCapacity:          pointer.From(input.Capacity),
		MinimumStartingCapacity: pointer.From(input.MinCapacity),
	}

	return append(outputList, output)
}

func (r ComputeFleetResource) flattenSpotCapacityModel(input *fleets.SpotPriorityProfile) []SpotCapacityModel {
	outputList := make([]SpotCapacityModel, 0)
	if input == nil {
		return outputList
	}

	output := SpotCapacityModel{
		AllocationStrategy:      string(pointer.From(input.AllocationStrategy)),
		TargetCapacity:          pointer.From(input.Capacity),
		EvictionPolicy:          string(pointer.From(input.EvictionPolicy)),
		MaintainCapacityEnabled: pointer.From(input.Maintain),
		MinimumCapacity:         pointer.From(input.MinCapacity),
	}

	// defaulted since MaxHourlyPricePerVM isn't returned if it's unset
	maxHourlyPricePerVM := float64(-1.0)
	if input.MaxPricePerVM != nil {
		maxHourlyPricePerVM = pointer.From(input.MaxPricePerVM)
	}
	output.MaxHourlyPricePerVM = maxHourlyPricePerVM

	return append(outputList, output)
}

func (r ComputeFleetResource) flattenVMSizeProfileModel(inputList *[]fleets.VMSizeProfile) []VMSizeProfileModel {
	outputList := make([]VMSizeProfileModel, 0)
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		output := VMSizeProfileModel{
			Name: input.Name,
			Rank: pointer.From(input.Rank),
		}
		outputList = append(outputList, output)
	}
	return outputList
}
