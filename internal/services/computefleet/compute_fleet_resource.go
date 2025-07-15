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
	RegularPriorityProfile   []RegularPriorityProfileModel              `tfschema:"regular_priority_profile"`
	SpotPriorityProfile      []SpotPriorityProfileModel                 `tfschema:"spot_priority_profile"`
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

type RegularPriorityProfileModel struct {
	AllocationStrategy string `tfschema:"allocation_strategy"`
	Capacity           int64  `tfschema:"capacity"`
	MinCapacity        int64  `tfschema:"min_capacity"`
}

type SpotPriorityProfileModel struct {
	AllocationStrategy      string  `tfschema:"allocation_strategy"`
	Capacity                int64   `tfschema:"capacity"`
	EvictionPolicy          string  `tfschema:"eviction_policy"`
	MaintainCapacityEnabled bool    `tfschema:"maintain_enabled"`
	MaxHourlyPricePerVM     float64 `tfschema:"max_hourly_price_per_vm"`
	MinCapacity             int64   `tfschema:"min_capacity"`
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

		"vm_sizes_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
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
						ValidateFunc: validation.IntBetween(0, 65535),
					},
				},
			},
		},

		"additional_capabilities": virtualMachineAdditionalCapabilitiesSchema(),

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

		"zones": commonschema.ZonesMultipleOptionalForceNew(),

		"compute_api_version": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Computed:     true,
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
					"capacity": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(1, 10000),
					},

					"allocation_strategy": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						Default:      string(fleets.RegularPriorityAllocationStrategyLowestPrice),
						ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForRegularPriorityAllocationStrategy(), false),
					},

					"min_capacity": {
						Type:         pluginsdk.TypeInt,
						ForceNew:     true,
						Optional:     true,
						Default:      0,
						ValidateFunc: validation.IntAtLeast(0),
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
					"capacity": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(1, 10000),
					},

					"allocation_strategy": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						Default:      string(fleets.SpotAllocationStrategyPriceCapacityOptimized),
						ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForSpotAllocationStrategy(), false),
					},

					"eviction_policy": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						Default:      string(fleets.EvictionPolicyDelete),
						ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForEvictionPolicy(), false),
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
						Default:      0,
						ValidateFunc: validation.IntAtLeast(0),
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func virtualMachineAdditionalCapabilitiesSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
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
				},
				"ultra_ssd_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					ForceNew: true,
					Default:  false,
				},
			},
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

			computeProfile := fleets.ComputeProfile{
				AdditionalVirtualMachineCapabilities: expandAdditionalCapabilities(model.AdditionalCapabilities),
				PlatformFaultDomainCount:             pointer.To(model.PlatformFaultDomainCount),
			}
			if model.ComputeApiVersion != "" {
				computeProfile.ComputeApiVersion = pointer.To(model.ComputeApiVersion)
			}

			baseVirtualMachineProfileValue, err := expandVirtualMachineProfileModel(model.VirtualMachineProfile, metadata.ResourceData)
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
				properties.Plan = expandPlanModel(model.Plan)
			}

			if metadata.ResourceData.HasChange("regular_priority_profile") {
				properties.Properties.RegularPriorityProfile = expandRegularPriorityProfileModel(model.RegularPriorityProfile)
			}

			if metadata.ResourceData.HasChange("spot_priority_profile") {
				properties.Properties.SpotPriorityProfile = expandSpotPriorityProfileModel(model.SpotPriorityProfile)
			}

			if metadata.ResourceData.HasChange("vm_sizes_profile") {
				properties.Properties.VMSizesProfile = pointer.From(expandVMSizeProfileModel(model.VMSizesProfile))
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

				state.Plan = flattenPlanModel(model.Plan)

				if props := model.Properties; props != nil {
					state.AdditionalCapabilities = flattenAdditionalCapabilities(props.ComputeProfile.AdditionalVirtualMachineCapabilities)

					baseVirtualMachineProfileValue, err := flattenVirtualMachineProfileModel(&props.ComputeProfile.BaseVirtualMachineProfile, metadata)
					if err != nil {
						return err
					}
					state.VirtualMachineProfile = baseVirtualMachineProfileValue

					state.ComputeApiVersion = pointer.From(props.ComputeProfile.ComputeApiVersion)
					state.PlatformFaultDomainCount = pointer.From(props.ComputeProfile.PlatformFaultDomainCount)
					state.RegularPriorityProfile = flattenRegularPriorityProfileModel(props.RegularPriorityProfile)
					state.SpotPriorityProfile = flattenSpotPriorityProfileModel(props.SpotPriorityProfile)
					state.UniqueId = pointer.From(props.UniqueId)
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

			if len(state.SpotPriorityProfile) > 0 && len(state.RegularPriorityProfile) > 0 {
				if state.SpotPriorityProfile[0].Capacity+state.RegularPriorityProfile[0].Capacity > 10000 {
					return fmt.Errorf("the sum of `spot_priority_profile.0.capacity` and `regular_priority_profile.0.capacity` must be between `0` and `10000`, inclusive")
				}
			}

			if len(state.SpotPriorityProfile) > 0 {
				if state.SpotPriorityProfile[0].MaintainCapacityEnabled {
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

			if v := state.VirtualMachineProfile[0].DataDisks; len(v) > 0 {
				storageAccountType := v[0].StorageAccountType
				ultraSSDEnabled := false
				if ac := state.AdditionalCapabilities; len(ac) > 0 {
					ultraSSDEnabled = ac[0].UltraSsdEnabled
				}

				if !ultraSSDEnabled && storageAccountType == string(fleets.StorageAccountTypesUltraSSDLRS) {
					return fmt.Errorf("`UltraSSD_LRS` storage account type can be used only when `ultra_ssd_enabled` is enalbed")
				}

				if v[0].CreateOption == string(fleets.DiskCreateOptionTypesEmpty) {
					if v[0].DiskSizeInGiB == 0 {
						return fmt.Errorf("`disk_size_in_gib` is required when`create_option` is `Empty`")
					}

					lunExist := metadata.ResourceDiff.GetRawConfig().AsValueMap()["virtual_machine_profile"].AsValueSlice()[0].AsValueMap()["data_disk"].AsValueSlice()[0].AsValueMap()["lun"]
					if lunExist.IsNull() {
						return fmt.Errorf("`lun` is required when`create_option` is `Empty`")
					}
				}
			}

			vmProfile := state.VirtualMachineProfile[0]
			if vmProfile.SourceImageId != "" && len(vmProfile.SourceImageReference) > 0 {
				return fmt.Errorf("only one of `source_image_id` and `source_image_reference` in `virtual_machine_profile` must be specified")
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
					return fmt.Errorf("`protected_settings_from_key_vault` cannot be used with `protected_settings_json`")
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
		Maintain:           pointer.To(input.MaintainCapacityEnabled),
		MinCapacity:        pointer.To(input.MinCapacity),
	}

	if input.MaxHourlyPricePerVM > 0 || input.MaxHourlyPricePerVM == -1 {
		output.MaxPricePerVM = pointer.To(input.MaxHourlyPricePerVM)
	}
	return &output
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

func expandAdditionalCapabilities(inputList []AdditionalCapabilitiesModel) *fleets.AdditionalCapabilities {
	if len(inputList) == 0 {
		return nil
	}

	capabilities := fleets.AdditionalCapabilities{
		UltraSSDEnabled:    pointer.To(inputList[0].UltraSsdEnabled),
		HibernationEnabled: pointer.To(inputList[0].HibernationEnabled),
	}

	return &capabilities
}

func flattenPlanModel(input *fleets.Plan) []PlanModel {
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

func flattenAdditionalCapabilities(input *fleets.AdditionalCapabilities) []AdditionalCapabilitiesModel {
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

func flattenRegularPriorityProfileModel(input *fleets.RegularPriorityProfile) []RegularPriorityProfileModel {
	outputList := make([]RegularPriorityProfileModel, 0)
	if input == nil {
		return outputList
	}

	output := RegularPriorityProfileModel{
		AllocationStrategy: string(pointer.From(input.AllocationStrategy)),
		Capacity:           pointer.From(input.Capacity),
		MinCapacity:        pointer.From(input.MinCapacity),
	}

	return append(outputList, output)
}

func flattenSpotPriorityProfileModel(input *fleets.SpotPriorityProfile) []SpotPriorityProfileModel {
	outputList := make([]SpotPriorityProfileModel, 0)
	if input == nil {
		return outputList
	}

	output := SpotPriorityProfileModel{
		AllocationStrategy:      string(pointer.From(input.AllocationStrategy)),
		Capacity:                pointer.From(input.Capacity),
		EvictionPolicy:          string(pointer.From(input.EvictionPolicy)),
		MaintainCapacityEnabled: pointer.From(input.Maintain),
		MinCapacity:             pointer.From(input.MinCapacity),
	}

	// defaulted since MaxHourlyPricePerVM isn't returned if it's unset
	maxHourlyPricePerVM := float64(-1.0)
	if input.MaxPricePerVM != nil {
		maxHourlyPricePerVM = pointer.From(input.MaxPricePerVM)
	}
	output.MaxHourlyPricePerVM = maxHourlyPricePerVM

	return append(outputList, output)
}

func flattenVMSizeProfileModel(inputList *[]fleets.VMSizeProfile) []VMSizeProfileModel {
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
