// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservationgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/images"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/proximityplacementgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-04-02/disks"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/base64"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/compute/2023-03-01/compute"
)

func resourceLinuxVirtualMachine() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLinuxVirtualMachineCreate,
		Read:   resourceLinuxVirtualMachineRead,
		Update: resourceLinuxVirtualMachineUpdate,
		Delete: resourceLinuxVirtualMachineDelete,
		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := commonids.ParseVirtualMachineID(id)
			return err
		}, importVirtualMachine(compute.OperatingSystemTypesLinux, "azurerm_linux_virtual_machine")),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(45 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(45 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(45 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: computeValidate.VirtualMachineName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			// Required
			"admin_username": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: computeValidate.LinuxAdminUsername,
			},

			"network_interface_ids": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: networkValidate.NetworkInterfaceID,
				},
			},

			"os_disk": virtualMachineOSDiskSchema(),

			"size": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			// Optional
			"additional_capabilities": virtualMachineAdditionalCapabilitiesSchema(),

			"admin_password": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				ForceNew:         true,
				Sensitive:        true,
				DiffSuppressFunc: adminPasswordDiffSuppressFunc,
				ValidateFunc:     computeValidate.LinuxAdminPassword,
			},

			"admin_ssh_key": SSHKeysSchema(true),

			"allow_extension_operations": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"availability_set_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateAvailabilitySetID,
				// the Compute/VM API is broken and returns the Availability Set name in UPPERCASE :shrug:
				// tracked by https://github.com/Azure/azure-rest-api-specs/issues/19424
				DiffSuppressFunc: suppress.CaseDifference,
				ConflictsWith: []string{
					"capacity_reservation_group_id",
					"virtual_machine_scale_set_id",
					"zone",
				},
			},

			"boot_diagnostics": bootDiagnosticsSchema(),

			"bypass_platform_safety_checks_on_user_schedule_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"capacity_reservation_group_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				// the Compute/VM API is broken and returns the Resource Group name in UPPERCASE
				// tracked by https://github.com/Azure/azure-rest-api-specs/issues/19424
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     capacityreservationgroups.ValidateCapacityReservationGroupID,
				ConflictsWith: []string{
					"availability_set_id",
					"proximity_placement_group_id",
				},
			},

			"computer_name": {
				Type:     pluginsdk.TypeString,
				Optional: true,

				// Computed since we reuse the VM name if one's not specified
				Computed: true,
				ForceNew: true,

				ValidateFunc: computeValidate.LinuxComputerNameFull,
			},

			"custom_data": base64.OptionalSchema(true),

			"dedicated_host_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: commonids.ValidateDedicatedHostID,
				// the Compute/VM API is broken and returns the Resource Group name in UPPERCASE :shrug:
				// tracked by https://github.com/Azure/azure-rest-api-specs/issues/19424
				DiffSuppressFunc: suppress.CaseDifference,
				ConflictsWith: []string{
					"dedicated_host_group_id",
				},
			},

			"dedicated_host_group_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: commonids.ValidateDedicatedHostGroupID,
				// the Compute/VM API is broken and returns the Resource Group name in UPPERCASE
				// tracked by https://github.com/Azure/azure-rest-api-specs/issues/19424
				DiffSuppressFunc: suppress.CaseDifference,
				ConflictsWith: []string{
					"dedicated_host_id",
				},
			},

			"disable_password_authentication": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},

			"disk_controller_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.DiskControllerTypesNVMe),
					string(compute.DiskControllerTypesSCSI),
				}, false),
			},

			"edge_zone": commonschema.EdgeZoneOptionalForceNew(),

			"encryption_at_host_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"eviction_policy": {
				// only applicable when `priority` is set to `Spot`
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.VirtualMachineEvictionPolicyTypesDeallocate),
					string(compute.VirtualMachineEvictionPolicyTypesDelete),
				}, false),
			},

			"extensions_time_budget": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      "PT1H30M",
				ValidateFunc: azValidate.ISO8601DurationBetween("PT15M", "PT2H"),
			},

			"gallery_application": VirtualMachineGalleryApplicationSchema(),

			"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

			"license_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"RHEL_BYOS",
					"SLES_BYOS",
				}, false),
			},

			"max_bid_price": {
				Type:         pluginsdk.TypeFloat,
				Optional:     true,
				Default:      -1,
				ValidateFunc: validation.FloatAtLeast(-1.0),
			},

			"plan": planSchema(),

			"priority": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(compute.VirtualMachinePriorityTypesRegular),
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.VirtualMachinePriorityTypesRegular),
					string(compute.VirtualMachinePriorityTypesSpot),
				}, false),
			},

			"provision_vm_agent": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},

			"patch_mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.LinuxVMGuestPatchModeAutomaticByPlatform),
					string(compute.LinuxVMGuestPatchModeImageDefault),
				}, false),
				Default: string(compute.LinuxVMGuestPatchModeImageDefault),
			},

			"patch_assessment_mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(compute.LinuxPatchAssessmentModeImageDefault),
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.LinuxPatchAssessmentModeAutomaticByPlatform),
					string(compute.LinuxPatchAssessmentModeImageDefault),
				}, false),
			},

			"proximity_placement_group_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: proximityplacementgroups.ValidateProximityPlacementGroupID,
				// the Compute/VM API is broken and returns the Resource Group name in UPPERCASE :shrug:
				// tracked by https://github.com/Azure/azure-rest-api-specs/issues/19424
				DiffSuppressFunc: suppress.CaseDifference,
				ConflictsWith: []string{
					"capacity_reservation_group_id",
				},
			},

			"reboot_setting": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.LinuxVMGuestPatchAutomaticByPlatformRebootSettingAlways),
					string(compute.LinuxVMGuestPatchAutomaticByPlatformRebootSettingIfRequired),
					string(compute.LinuxVMGuestPatchAutomaticByPlatformRebootSettingNever),
				}, false),
			},

			"secret": linuxSecretSchema(),

			"secure_boot_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
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
					"source_image_id",
					"source_image_reference",
				},
			},

			"source_image_reference": sourceImageReferenceSchema(true),

			"virtual_machine_scale_set_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ConflictsWith: []string{
					"availability_set_id",
				},
				ValidateFunc: commonids.ValidateVirtualMachineScaleSetID,
			},

			"vm_agent_platform_updates_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"vtpm_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"platform_fault_domain": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      -1,
				ForceNew:     true,
				RequiredWith: []string{"virtual_machine_scale_set_id"},
				ValidateFunc: validation.IntAtLeast(-1),
			},

			"tags": tags.Schema(),

			"os_image_notification": virtualMachineOsImageNotificationSchema(),

			"termination_notification": virtualMachineTerminationNotificationSchema(),

			"user_data": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsBase64,
			},

			"zone": commonschema.ZoneSingleOptionalForceNew(),

			// Computed
			"private_ip_address": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"private_ip_addresses": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
			"public_ip_address": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"public_ip_addresses": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
			"virtual_machine_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceLinuxVirtualMachineCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VMClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewVirtualMachineID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	locks.ByName(id.VirtualMachineName, VirtualMachineResourceName)
	defer locks.UnlockByName(id.VirtualMachineName, VirtualMachineResourceName)

	resp, err := client.Get(ctx, id.ResourceGroupName, id.VirtualMachineName, "")
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(resp.Response) {
		return tf.ImportAsExistsError("azurerm_linux_virtual_machine", *resp.ID)
	}

	additionalCapabilitiesRaw := d.Get("additional_capabilities").([]interface{})
	additionalCapabilities := expandVirtualMachineAdditionalCapabilities(additionalCapabilitiesRaw)

	adminUsername := d.Get("admin_username").(string)
	allowExtensionOperations := d.Get("allow_extension_operations").(bool)

	bootDiagnosticsRaw := d.Get("boot_diagnostics").([]interface{})
	bootDiagnostics := expandBootDiagnostics(bootDiagnosticsRaw)

	var computerName string
	if v, ok := d.GetOk("computer_name"); ok && len(v.(string)) > 0 {
		computerName = v.(string)
	} else {
		_, errs := computeValidate.LinuxComputerNameFull(d.Get("name"), "computer_name")
		if len(errs) > 0 {
			return fmt.Errorf("unable to assume default computer name %s Please adjust the %q, or specify an explicit %q", errs[0], "name", "computer_name")
		}
		computerName = id.VirtualMachineName
	}
	disablePasswordAuthentication := d.Get("disable_password_authentication").(bool)
	vmAgentPlatformUpdatesEnabled := d.Get("vm_agent_platform_updates_enabled").(bool)
	location := azure.NormalizeLocation(d.Get("location").(string))
	identityRaw := d.Get("identity").([]interface{})
	identity, err := expandVirtualMachineIdentity(identityRaw)
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}
	planRaw := d.Get("plan").([]interface{})
	plan := expandPlan(planRaw)
	priority := compute.VirtualMachinePriorityTypes(d.Get("priority").(string))
	provisionVMAgent := d.Get("provision_vm_agent").(bool)
	size := d.Get("size").(string)
	t := d.Get("tags").(map[string]interface{})

	networkInterfaceIdsRaw := d.Get("network_interface_ids").([]interface{})
	networkInterfaceIds := expandVirtualMachineNetworkInterfaceIDs(networkInterfaceIdsRaw)

	osDiskRaw := d.Get("os_disk").([]interface{})
	osDisk, err := expandVirtualMachineOSDisk(osDiskRaw, compute.OperatingSystemTypesLinux)
	if err != nil {
		return fmt.Errorf("expanding `os_disk`: %+v", err)
	}
	securityEncryptionType := osDiskRaw[0].(map[string]interface{})["security_encryption_type"].(string)

	secretsRaw := d.Get("secret").([]interface{})
	secrets := expandLinuxSecrets(secretsRaw)

	sourceImageReferenceRaw := d.Get("source_image_reference").([]interface{})
	sourceImageId := d.Get("source_image_id").(string)
	sourceImageReference := expandSourceImageReference(sourceImageReferenceRaw, sourceImageId)

	sshKeysRaw := d.Get("admin_ssh_key").(*pluginsdk.Set).List()
	sshKeys := ExpandSSHKeys(sshKeysRaw)

	params := compute.VirtualMachine{
		Name:             utils.String(id.VirtualMachineName),
		ExtendedLocation: expandEdgeZone(d.Get("edge_zone").(string)),
		Location:         utils.String(location),
		Identity:         identity,
		Plan:             plan,
		VirtualMachineProperties: &compute.VirtualMachineProperties{
			ApplicationProfile: &compute.ApplicationProfile{
				GalleryApplications: expandVirtualMachineGalleryApplication(d.Get("gallery_application").([]interface{})),
			},
			HardwareProfile: &compute.HardwareProfile{
				VMSize: compute.VirtualMachineSizeTypes(size),
			},
			OsProfile: &compute.OSProfile{
				AdminUsername:            utils.String(adminUsername),
				ComputerName:             utils.String(computerName),
				AllowExtensionOperations: utils.Bool(allowExtensionOperations),
				LinuxConfiguration: &compute.LinuxConfiguration{
					DisablePasswordAuthentication: utils.Bool(disablePasswordAuthentication),
					EnableVMAgentPlatformUpdates:  utils.Bool(vmAgentPlatformUpdatesEnabled),
					ProvisionVMAgent:              utils.Bool(provisionVMAgent),
					SSH: &compute.SSHConfiguration{
						PublicKeys: &sshKeys,
					},
				},
				Secrets: secrets,
			},
			NetworkProfile: &compute.NetworkProfile{
				NetworkInterfaces: &networkInterfaceIds,
			},
			Priority: priority,
			StorageProfile: &compute.StorageProfile{
				ImageReference: sourceImageReference,
				OsDisk:         osDisk,

				// Data Disks are instead handled via the Association resource - as such we can send an empty value here
				// but for Updates this'll need to be nil, else any associations will be overwritten
				DataDisks: &[]compute.DataDisk{},
			},

			// Optional
			AdditionalCapabilities: additionalCapabilities,
			DiagnosticsProfile:     bootDiagnostics,
			ExtensionsTimeBudget:   utils.String(d.Get("extensions_time_budget").(string)),
		},
		Tags: tags.Expand(t),
	}

	if diskControllerType, ok := d.GetOk("disk_controller_type"); ok {
		params.StorageProfile.DiskControllerType = compute.DiskControllerTypes(diskControllerType.(string))
	}

	if encryptionAtHostEnabled, ok := d.GetOk("encryption_at_host_enabled"); ok {
		if encryptionAtHostEnabled.(bool) {
			if compute.SecurityEncryptionTypesDiskWithVMGuestState == compute.SecurityEncryptionTypes(securityEncryptionType) {
				return fmt.Errorf("`encryption_at_host_enabled` cannot be set to `true` when `os_disk.0.security_encryption_type` is set to `DiskWithVMGuestState`")
			}
		}

		params.VirtualMachineProperties.SecurityProfile = &compute.SecurityProfile{
			EncryptionAtHost: utils.Bool(encryptionAtHostEnabled.(bool)),
		}
	}

	if v, ok := d.GetOk("license_type"); ok {
		params.VirtualMachineProperties.LicenseType = utils.String(v.(string))
	}

	patchMode := d.Get("patch_mode").(string)
	if patchMode != "" {
		if patchMode == string(compute.LinuxVMGuestPatchModeAutomaticByPlatform) && !provisionVMAgent {
			return fmt.Errorf("%q cannot be set to %q when %q is set to %q", "patch_mode", "AutomaticByPlatform", "provision_vm_agent", "false")
		}

		params.VirtualMachineProperties.OsProfile.LinuxConfiguration.PatchSettings = &compute.LinuxPatchSettings{
			PatchMode: compute.LinuxVMGuestPatchMode(patchMode),
		}
	}

	if v, ok := d.GetOk("patch_assessment_mode"); ok {
		if v.(string) == string(compute.LinuxPatchAssessmentModeAutomaticByPlatform) && !provisionVMAgent {
			return fmt.Errorf("`provision_vm_agent` must be set to `true` when `patch_assessment_mode` is set to `AutomaticByPlatform`")
		}

		if params.VirtualMachineProperties.OsProfile.LinuxConfiguration.PatchSettings == nil {
			params.VirtualMachineProperties.OsProfile.LinuxConfiguration.PatchSettings = &compute.LinuxPatchSettings{}
		}
		params.VirtualMachineProperties.OsProfile.LinuxConfiguration.PatchSettings.AssessmentMode = compute.LinuxPatchAssessmentMode(v.(string))
	}

	if d.Get("bypass_platform_safety_checks_on_user_schedule_enabled").(bool) {
		if patchMode != string(compute.LinuxVMGuestPatchModeAutomaticByPlatform) {
			return fmt.Errorf("`patch_mode` must be set to `AutomaticByPlatform` when `bypass_platform_safety_checks_on_user_schedule_enabled` is set to `true`")
		}

		if params.VirtualMachineProperties.OsProfile.LinuxConfiguration.PatchSettings == nil {
			params.VirtualMachineProperties.OsProfile.LinuxConfiguration.PatchSettings = &compute.LinuxPatchSettings{}
		}

		if params.VirtualMachineProperties.OsProfile.LinuxConfiguration.PatchSettings.AutomaticByPlatformSettings == nil {
			params.VirtualMachineProperties.OsProfile.LinuxConfiguration.PatchSettings.AutomaticByPlatformSettings = &compute.LinuxVMGuestPatchAutomaticByPlatformSettings{}
		}

		params.VirtualMachineProperties.OsProfile.LinuxConfiguration.PatchSettings.AutomaticByPlatformSettings.BypassPlatformSafetyChecksOnUserSchedule = pointer.To(true)
	}

	if v, ok := d.GetOk("reboot_setting"); ok {
		if patchMode != string(compute.LinuxVMGuestPatchModeAutomaticByPlatform) {
			return fmt.Errorf("`patch_mode` must be set to `AutomaticByPlatform` when `reboot_setting` is specified")
		}

		if params.VirtualMachineProperties.OsProfile.LinuxConfiguration.PatchSettings == nil {
			params.VirtualMachineProperties.OsProfile.LinuxConfiguration.PatchSettings = &compute.LinuxPatchSettings{}
		}

		if params.VirtualMachineProperties.OsProfile.LinuxConfiguration.PatchSettings.AutomaticByPlatformSettings == nil {
			params.VirtualMachineProperties.OsProfile.LinuxConfiguration.PatchSettings.AutomaticByPlatformSettings = &compute.LinuxVMGuestPatchAutomaticByPlatformSettings{}
		}

		params.VirtualMachineProperties.OsProfile.LinuxConfiguration.PatchSettings.AutomaticByPlatformSettings.RebootSetting = compute.LinuxVMGuestPatchAutomaticByPlatformRebootSetting(v.(string))
	}

	secureBootEnabled := d.Get("secure_boot_enabled").(bool)
	vtpmEnabled := d.Get("vtpm_enabled").(bool)
	if securityEncryptionType != "" {
		if compute.SecurityEncryptionTypesDiskWithVMGuestState == compute.SecurityEncryptionTypes(securityEncryptionType) && !secureBootEnabled {
			return fmt.Errorf("`secure_boot_enabled` must be set to `true` when `os_disk.0.security_encryption_type` is set to `DiskWithVMGuestState`")
		}
		if !vtpmEnabled {
			return fmt.Errorf("`vtpm_enabled` must be set to `true` when `os_disk.0.security_encryption_type` is specified")
		}

		if params.VirtualMachineProperties.SecurityProfile == nil {
			params.VirtualMachineProperties.SecurityProfile = &compute.SecurityProfile{}
		}
		params.VirtualMachineProperties.SecurityProfile.SecurityType = compute.SecurityTypesConfidentialVM

		if params.VirtualMachineProperties.SecurityProfile.UefiSettings == nil {
			params.VirtualMachineProperties.SecurityProfile.UefiSettings = &compute.UefiSettings{}
		}
		params.VirtualMachineProperties.SecurityProfile.UefiSettings.SecureBootEnabled = utils.Bool(secureBootEnabled)
		params.VirtualMachineProperties.SecurityProfile.UefiSettings.VTpmEnabled = utils.Bool(vtpmEnabled)
	} else {
		if secureBootEnabled {
			if params.VirtualMachineProperties.SecurityProfile == nil {
				params.VirtualMachineProperties.SecurityProfile = &compute.SecurityProfile{}
			}
			if params.VirtualMachineProperties.SecurityProfile.UefiSettings == nil {
				params.VirtualMachineProperties.SecurityProfile.UefiSettings = &compute.UefiSettings{}
			}
			params.VirtualMachineProperties.SecurityProfile.SecurityType = compute.SecurityTypesTrustedLaunch
			params.VirtualMachineProperties.SecurityProfile.UefiSettings.SecureBootEnabled = utils.Bool(secureBootEnabled)
		}

		if vtpmEnabled {
			if params.VirtualMachineProperties.SecurityProfile == nil {
				params.VirtualMachineProperties.SecurityProfile = &compute.SecurityProfile{}
			}
			if params.VirtualMachineProperties.SecurityProfile.UefiSettings == nil {
				params.VirtualMachineProperties.SecurityProfile.UefiSettings = &compute.UefiSettings{}
			}
			params.VirtualMachineProperties.SecurityProfile.SecurityType = compute.SecurityTypesTrustedLaunch
			params.VirtualMachineProperties.SecurityProfile.UefiSettings.VTpmEnabled = utils.Bool(vtpmEnabled)
		}
	}

	var osImageNotificationProfile *compute.OSImageNotificationProfile
	var terminateNotificationProfile *compute.TerminateNotificationProfile

	if v, ok := d.GetOk("os_image_notification"); ok {
		osImageNotificationProfile = expandOsImageNotificationProfile(v.([]interface{}))
	}

	if v, ok := d.GetOk("termination_notification"); ok {
		terminateNotificationProfile = expandTerminateNotificationProfile(v.([]interface{}))
	}

	if terminateNotificationProfile != nil || osImageNotificationProfile != nil {
		params.VirtualMachineProperties.ScheduledEventsProfile = &compute.ScheduledEventsProfile{
			OsImageNotificationProfile:   osImageNotificationProfile,
			TerminateNotificationProfile: terminateNotificationProfile,
		}
	}

	if !provisionVMAgent && allowExtensionOperations {
		return fmt.Errorf("`allow_extension_operations` cannot be set to `true` when `provision_vm_agent` is set to `false`")
	}

	if v, ok := d.GetOk("availability_set_id"); ok {
		params.AvailabilitySet = &compute.SubResource{
			ID: utils.String(v.(string)),
		}
	}

	if v, ok := d.GetOk("capacity_reservation_group_id"); ok {
		params.CapacityReservation = &compute.CapacityReservationProfile{
			CapacityReservationGroup: &compute.SubResource{
				ID: utils.String(v.(string)),
			},
		}
	}

	if v, ok := d.GetOk("custom_data"); ok {
		params.OsProfile.CustomData = utils.String(v.(string))
	}

	if v, ok := d.GetOk("dedicated_host_id"); ok {
		params.Host = &compute.SubResource{
			ID: utils.String(v.(string)),
		}
	}

	if v, ok := d.GetOk("dedicated_host_group_id"); ok {
		params.HostGroup = &compute.SubResource{
			ID: utils.String(v.(string)),
		}
	}

	if evictionPolicyRaw, ok := d.GetOk("eviction_policy"); ok {
		if params.Priority != compute.VirtualMachinePriorityTypesSpot {
			return fmt.Errorf("an `eviction_policy` can only be specified when `priority` is set to `Spot`")
		}

		params.EvictionPolicy = compute.VirtualMachineEvictionPolicyTypes(evictionPolicyRaw.(string))
	} else if priority == compute.VirtualMachinePriorityTypesSpot {
		return fmt.Errorf("an `eviction_policy` must be specified when `priority` is set to `Spot`")
	}

	if v, ok := d.Get("max_bid_price").(float64); ok && v > 0 {
		if priority != compute.VirtualMachinePriorityTypesSpot {
			return fmt.Errorf("`max_bid_price` can only be configured when `priority` is set to `Spot`")
		}

		params.BillingProfile = &compute.BillingProfile{
			MaxPrice: utils.Float(v),
		}
	}

	if v, ok := d.GetOk("proximity_placement_group_id"); ok {
		params.ProximityPlacementGroup = &compute.SubResource{
			ID: utils.String(v.(string)),
		}
	}

	if v, ok := d.GetOk("virtual_machine_scale_set_id"); ok {
		params.VirtualMachineScaleSet = &compute.SubResource{
			ID: utils.String(v.(string)),
		}
	}

	platformFaultDomain := d.Get("platform_fault_domain").(int)
	if platformFaultDomain != -1 {
		params.PlatformFaultDomain = utils.Int32(int32(platformFaultDomain))
	}

	if v, ok := d.GetOk("user_data"); ok {
		params.UserData = utils.String(v.(string))
	}

	if v, ok := d.GetOk("zone"); ok {
		params.Zones = &[]string{
			v.(string),
		}
	}

	// "Authentication using either SSH or by user name and password must be enabled in Linux profile." Target="linuxConfiguration"
	adminPassword := d.Get("admin_password").(string)
	if disablePasswordAuthentication && len(sshKeys) == 0 {
		return fmt.Errorf("at least one `admin_ssh_key` must be specified when `disable_password_authentication` is set to `true`")
	} else if !disablePasswordAuthentication {
		if adminPassword == "" {
			return fmt.Errorf("an `admin_password` must be specified if `disable_password_authentication` is set to `false`")
		}

		params.OsProfile.AdminPassword = utils.String(adminPassword)
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroupName, id.VirtualMachineName, params)
	if err != nil {
		return fmt.Errorf("creating Linux %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Linux %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceLinuxVirtualMachineRead(d, meta)
}

func resourceLinuxVirtualMachineRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VMClient
	disksClient := meta.(*clients.Client).Compute.DisksClient
	networkInterfacesClient := meta.(*clients.Client).Network.InterfacesClient
	publicIPAddressesClient := meta.(*clients.Client).Network.PublicIPsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseVirtualMachineID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroupName, id.VirtualMachineName, compute.InstanceViewTypesUserData)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Linux %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Linux %s: %+v", id, err)
	}

	d.Set("name", id.VirtualMachineName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("location", location.NormalizeNilable(resp.Location))
	d.Set("edge_zone", flattenEdgeZone(resp.ExtendedLocation))

	identity, err := flattenVirtualMachineIdentity(resp.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	if err := d.Set("plan", flattenPlan(resp.Plan)); err != nil {
		return fmt.Errorf("setting `plan`: %+v", err)
	}

	if resp.VirtualMachineProperties == nil {
		return fmt.Errorf("retrieving Linux %s: `properties` was nil", id)
	}

	props := *resp.VirtualMachineProperties
	if err := d.Set("additional_capabilities", flattenVirtualMachineAdditionalCapabilities(props.AdditionalCapabilities)); err != nil {
		return fmt.Errorf("setting `additional_capabilities`: %+v", err)
	}

	availabilitySetId := ""
	if props.AvailabilitySet != nil && props.AvailabilitySet.ID != nil {
		availabilitySetId = *props.AvailabilitySet.ID
	}
	d.Set("availability_set_id", availabilitySetId)

	capacityReservationGroupId := ""
	if props.CapacityReservation != nil && props.CapacityReservation.CapacityReservationGroup != nil && props.CapacityReservation.CapacityReservationGroup.ID != nil {
		capacityReservationGroupId = *props.CapacityReservation.CapacityReservationGroup.ID
	}
	d.Set("capacity_reservation_group_id", capacityReservationGroupId)

	if props.ApplicationProfile != nil && props.ApplicationProfile.GalleryApplications != nil {
		d.Set("gallery_application", flattenVirtualMachineGalleryApplication(props.ApplicationProfile.GalleryApplications))
	}

	licenseType := ""
	if props.LicenseType != nil {
		licenseType = *props.LicenseType
	}
	d.Set("license_type", licenseType)

	if err := d.Set("boot_diagnostics", flattenBootDiagnostics(props.DiagnosticsProfile)); err != nil {
		return fmt.Errorf("setting `boot_diagnostics`: %+v", err)
	}

	d.Set("eviction_policy", string(props.EvictionPolicy))
	if profile := props.HardwareProfile; profile != nil {
		d.Set("size", string(profile.VMSize))
	}

	extensionsTimeBudget := "PT1H30M"
	if props.ExtensionsTimeBudget != nil {
		extensionsTimeBudget = *props.ExtensionsTimeBudget
	}
	d.Set("extensions_time_budget", extensionsTimeBudget)

	// defaulted since BillingProfile isn't returned if it's unset
	maxBidPrice := float64(-1.0)
	if props.BillingProfile != nil && props.BillingProfile.MaxPrice != nil {
		maxBidPrice = *props.BillingProfile.MaxPrice
	}
	d.Set("max_bid_price", maxBidPrice)

	if profile := props.NetworkProfile; profile != nil {
		if err := d.Set("network_interface_ids", flattenVirtualMachineNetworkInterfaceIDs(props.NetworkProfile.NetworkInterfaces)); err != nil {
			return fmt.Errorf("setting `network_interface_ids`: %+v", err)
		}
	}

	dedicatedHostId := ""
	if props.Host != nil && props.Host.ID != nil {
		dedicatedHostId = *props.Host.ID
	}
	d.Set("dedicated_host_id", dedicatedHostId)

	dedicatedHostGroupId := ""
	if props.HostGroup != nil && props.HostGroup.ID != nil {
		dedicatedHostGroupId = *props.HostGroup.ID
	}
	d.Set("dedicated_host_group_id", dedicatedHostGroupId)

	virtualMachineScaleSetId := ""
	if props.VirtualMachineScaleSet != nil && props.VirtualMachineScaleSet.ID != nil {
		virtualMachineScaleSetId = *props.VirtualMachineScaleSet.ID
	}
	d.Set("virtual_machine_scale_set_id", virtualMachineScaleSetId)
	platformFaultDomain := -1
	if props.PlatformFaultDomain != nil {
		platformFaultDomain = int(*props.PlatformFaultDomain)
	}
	d.Set("platform_fault_domain", platformFaultDomain)

	if profile := props.OsProfile; profile != nil {
		d.Set("admin_username", profile.AdminUsername)
		d.Set("allow_extension_operations", profile.AllowExtensionOperations)
		d.Set("computer_name", profile.ComputerName)

		if config := profile.LinuxConfiguration; config != nil {
			d.Set("disable_password_authentication", config.DisablePasswordAuthentication)
			d.Set("provision_vm_agent", config.ProvisionVMAgent)
			d.Set("vm_agent_platform_updates_enabled", config.EnableVMAgentPlatformUpdates)

			flattenedSSHKeys, err := FlattenSSHKeys(config.SSH)
			if err != nil {
				return fmt.Errorf("flattening `admin_ssh_key`: %+v", err)
			}
			if err := d.Set("admin_ssh_key", pluginsdk.NewSet(SSHKeySchemaHash, *flattenedSSHKeys)); err != nil {
				return fmt.Errorf("setting `admin_ssh_key`: %+v", err)
			}
			patchMode := string(compute.LinuxVMGuestPatchModeImageDefault)
			if patchSettings := config.PatchSettings; patchSettings != nil && patchSettings.PatchMode != "" {
				patchMode = string(patchSettings.PatchMode)
			}
			d.Set("patch_mode", patchMode)

			assessmentMode := string(compute.LinuxPatchAssessmentModeImageDefault)
			if patchSettings := config.PatchSettings; patchSettings != nil && patchSettings.AssessmentMode != "" {
				assessmentMode = string(patchSettings.AssessmentMode)
			}
			d.Set("patch_assessment_mode", assessmentMode)

			bypassPlatformSafetyChecksOnUserScheduleEnabled := false
			rebootSetting := ""
			if patchSettings := config.PatchSettings; patchSettings != nil && patchSettings.AutomaticByPlatformSettings != nil {
				bypassPlatformSafetyChecksOnUserScheduleEnabled = pointer.From(patchSettings.AutomaticByPlatformSettings.BypassPlatformSafetyChecksOnUserSchedule)
				rebootSetting = string(patchSettings.AutomaticByPlatformSettings.RebootSetting)
			}
			d.Set("bypass_platform_safety_checks_on_user_schedule_enabled", bypassPlatformSafetyChecksOnUserScheduleEnabled)
			d.Set("reboot_setting", rebootSetting)
		}

		if err := d.Set("secret", flattenLinuxSecrets(profile.Secrets)); err != nil {
			return fmt.Errorf("setting `secret`: %+v", err)
		}
	}
	// Resources created with azurerm_virtual_machine have priority set to ""
	// We need to treat "" as equal to "Regular" to allow migration azurerm_virtual_machine -> azurerm_linux_virtual_machine
	priority := string(compute.VirtualMachinePriorityTypesRegular)
	if props.Priority != "" {
		priority = string(props.Priority)
	}
	d.Set("priority", priority)
	proximityPlacementGroupId := ""
	if props.ProximityPlacementGroup != nil && props.ProximityPlacementGroup.ID != nil {
		proximityPlacementGroupId = *props.ProximityPlacementGroup.ID
	}
	d.Set("proximity_placement_group_id", proximityPlacementGroupId)

	if profile := props.StorageProfile; profile != nil {
		d.Set("disk_controller_type", string(props.StorageProfile.DiskControllerType))

		// the storage_account_type isn't returned so we need to look it up
		flattenedOSDisk, err := flattenVirtualMachineOSDisk(ctx, disksClient, profile.OsDisk)
		if err != nil {
			return fmt.Errorf("flattening `os_disk`: %+v", err)
		}
		if err := d.Set("os_disk", flattenedOSDisk); err != nil {
			return fmt.Errorf("settings `os_disk`: %+v", err)
		}

		var storageImageId string
		if profile.ImageReference != nil && profile.ImageReference.ID != nil {
			storageImageId = *profile.ImageReference.ID
		}
		if profile.ImageReference != nil && profile.ImageReference.CommunityGalleryImageID != nil {
			storageImageId = *profile.ImageReference.CommunityGalleryImageID
		}
		if profile.ImageReference != nil && profile.ImageReference.SharedGalleryImageID != nil {
			storageImageId = *profile.ImageReference.SharedGalleryImageID
		}

		d.Set("source_image_id", storageImageId)

		if err := d.Set("source_image_reference", flattenSourceImageReference(profile.ImageReference, storageImageId != "")); err != nil {
			return fmt.Errorf("setting `source_image_reference`: %+v", err)
		}
	}

	if scheduleProfile := props.ScheduledEventsProfile; scheduleProfile != nil {
		if err := d.Set("os_image_notification", flattenOsImageNotificationProfile(scheduleProfile.OsImageNotificationProfile)); err != nil {
			return fmt.Errorf("setting `termination_notification`: %+v", err)
		}

		if err := d.Set("termination_notification", flattenTerminateNotificationProfile(scheduleProfile.TerminateNotificationProfile)); err != nil {
			return fmt.Errorf("setting `termination_notification`: %+v", err)
		}
	}

	encryptionAtHostEnabled := false
	vtpmEnabled := false
	secureBootEnabled := false

	if secprofile := props.SecurityProfile; secprofile != nil {
		if secprofile.EncryptionAtHost != nil {
			encryptionAtHostEnabled = *secprofile.EncryptionAtHost
		}
		if uefi := props.SecurityProfile.UefiSettings; uefi != nil {
			if uefi.VTpmEnabled != nil {
				vtpmEnabled = *uefi.VTpmEnabled
			}
			if uefi.SecureBootEnabled != nil {
				secureBootEnabled = *uefi.SecureBootEnabled
			}
		}
	}

	d.Set("encryption_at_host_enabled", encryptionAtHostEnabled)
	d.Set("vtpm_enabled", vtpmEnabled)
	d.Set("secure_boot_enabled", secureBootEnabled)

	d.Set("virtual_machine_id", props.VMID)

	d.Set("user_data", props.UserData)

	zone := ""
	if resp.Zones != nil {
		if zones := *resp.Zones; len(zones) > 0 {
			zone = zones[0]
		}
	}
	d.Set("zone", zone)

	connectionInfo := retrieveConnectionInformation(ctx, networkInterfacesClient, publicIPAddressesClient, resp.VirtualMachineProperties)
	d.Set("private_ip_address", connectionInfo.primaryPrivateAddress)
	d.Set("private_ip_addresses", connectionInfo.privateAddresses)
	d.Set("public_ip_address", connectionInfo.primaryPublicAddress)
	d.Set("public_ip_addresses", connectionInfo.publicAddresses)
	isWindows := false
	setConnectionInformation(d, connectionInfo, isWindows)

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceLinuxVirtualMachineUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VMClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseVirtualMachineID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.VirtualMachineName, VirtualMachineResourceName)
	defer locks.UnlockByName(id.VirtualMachineName, VirtualMachineResourceName)

	log.Printf("[DEBUG] Retrieving Linux %s", id)
	existing, err := client.Get(ctx, id.ResourceGroupName, id.VirtualMachineName, compute.InstanceViewTypesUserData)
	if err != nil {
		return fmt.Errorf("retrieving Linux %s: %+v", id, err)
	}

	log.Printf("[DEBUG] Retrieving InstanceView for Linux %s.", id)
	instanceView, err := client.InstanceView(ctx, id.ResourceGroupName, id.VirtualMachineName)
	if err != nil {
		return fmt.Errorf("retrieving InstanceView for Linux %s: %+v", id, err)
	}

	shouldTurnBackOn := virtualMachineShouldBeStarted(instanceView)
	hasEphemeralOSDisk := false
	if props := existing.VirtualMachineProperties; props != nil {
		if storage := props.StorageProfile; storage != nil {
			if disk := storage.OsDisk; disk != nil {
				if settings := disk.DiffDiskSettings; settings != nil {
					hasEphemeralOSDisk = settings.Option == compute.DiffDiskOptionsLocal
				}
			}
		}
	}

	shouldUpdate := false
	shouldShutDown := false
	shouldDeallocate := false

	update := compute.VirtualMachineUpdate{
		VirtualMachineProperties: &compute.VirtualMachineProperties{},
	}

	if d.HasChange("boot_diagnostics") {
		shouldUpdate = true

		bootDiagnosticsRaw := d.Get("boot_diagnostics").([]interface{})
		update.VirtualMachineProperties.DiagnosticsProfile = expandBootDiagnostics(bootDiagnosticsRaw)
	}

	if d.HasChange("secret") {
		shouldUpdate = true

		profile := compute.OSProfile{}

		if d.HasChange("secret") {
			secretsRaw := d.Get("secret").([]interface{})
			profile.Secrets = expandLinuxSecrets(secretsRaw)
		}

		update.VirtualMachineProperties.OsProfile = &profile
	}

	if d.HasChange("identity") {
		shouldUpdate = true

		identityRaw := d.Get("identity").([]interface{})
		identity, err := expandVirtualMachineIdentity(identityRaw)
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		update.Identity = identity
	}

	if d.HasChange("license_type") {
		shouldUpdate = true

		if v, ok := d.GetOk("license_type"); ok {
			update.LicenseType = utils.String(v.(string))
		}
	}

	if d.HasChange("capacity_reservation_group_id") {
		shouldUpdate = true
		shouldDeallocate = true

		if v, ok := d.GetOk("capacity_reservation_group_id"); ok {
			update.CapacityReservation = &compute.CapacityReservationProfile{
				CapacityReservationGroup: &compute.SubResource{
					ID: utils.String(v.(string)),
				},
			}
		} else {
			update.CapacityReservation = &compute.CapacityReservationProfile{
				CapacityReservationGroup: &compute.SubResource{},
			}
		}
	}

	if d.HasChange("dedicated_host_id") {
		shouldUpdate = true

		// Code="PropertyChangeNotAllowed" Message="Updating Host of VM 'VMNAME' is not allowed as the VM is currently allocated. Please Deallocate the VM and retry the operation."
		shouldDeallocate = true

		if v, ok := d.GetOk("dedicated_host_id"); ok {
			update.Host = &compute.SubResource{
				ID: utils.String(v.(string)),
			}
		} else {
			update.Host = &compute.SubResource{}
		}
	}

	if d.HasChange("dedicated_host_group_id") {
		shouldUpdate = true

		// Code="PropertyChangeNotAllowed" Message="Updating Host of VM 'VMNAME' is not allowed as the VM is currently allocated. Please Deallocate the VM and retry the operation."
		shouldDeallocate = true

		if v, ok := d.GetOk("dedicated_host_group_id"); ok {
			update.HostGroup = &compute.SubResource{
				ID: utils.String(v.(string)),
			}
		} else {
			update.HostGroup = &compute.SubResource{}
		}
	}

	if d.HasChange("extensions_time_budget") {
		shouldUpdate = true
		update.ExtensionsTimeBudget = utils.String(d.Get("extensions_time_budget").(string))
	}

	if d.HasChange("gallery_application") {
		shouldUpdate = true
		update.ApplicationProfile = &compute.ApplicationProfile{
			GalleryApplications: expandVirtualMachineGalleryApplication(d.Get("gallery_application").([]interface{})),
		}
	}

	if d.HasChange("max_bid_price") {
		shouldUpdate = true

		// Code="OperationNotAllowed" Message="Max price change is not allowed. For more information, see http://aka.ms/AzureSpot/errormessages"
		shouldShutDown = true

		// "code":"OperationNotAllowed"
		// "message": "Max price change is not allowed when the VM [name] is currently allocated.
		//			   Please deallocate and try again. For more information, see http://aka.ms/AzureSpot/errormessages"
		shouldDeallocate = true

		maxBidPrice := d.Get("max_bid_price").(float64)
		update.VirtualMachineProperties.BillingProfile = &compute.BillingProfile{
			MaxPrice: utils.Float(maxBidPrice),
		}
	}

	if d.HasChange("network_interface_ids") {
		shouldUpdate = true

		// Code="CannotAddOrRemoveNetworkInterfacesFromARunningVirtualMachine"
		// Message="Secondary network interfaces cannot be added or removed from a running virtual machine.
		shouldShutDown = true

		// @tombuildsstuff: after testing shutting it down isn't sufficient - we need a full deallocation
		shouldDeallocate = true

		networkInterfaceIdsRaw := d.Get("network_interface_ids").([]interface{})
		networkInterfaceIds := expandVirtualMachineNetworkInterfaceIDs(networkInterfaceIdsRaw)

		update.VirtualMachineProperties.NetworkProfile = &compute.NetworkProfile{
			NetworkInterfaces: &networkInterfaceIds,
		}
	}

	if d.HasChange("disk_controller_type") {
		shouldUpdate = true
		shouldDeallocate = true

		if update.VirtualMachineProperties.StorageProfile == nil {
			update.VirtualMachineProperties.StorageProfile = &compute.StorageProfile{}
		}

		update.VirtualMachineProperties.StorageProfile.DiskControllerType = compute.DiskControllerTypes(d.Get("disk_controller_type").(string))
	}

	if d.HasChange("os_disk") {
		shouldUpdate = true

		// Code="Conflict" Message="Disk resizing is allowed only when creating a VM or when the VM is deallocated." Target="disk.diskSizeGB"
		shouldShutDown = true
		shouldDeallocate = true

		osDiskRaw := d.Get("os_disk").([]interface{})
		osDisk, err := expandVirtualMachineOSDisk(osDiskRaw, compute.OperatingSystemTypesLinux)
		if err != nil {
			return fmt.Errorf("expanding `os_disk`: %+v", err)
		}

		if update.VirtualMachineProperties.StorageProfile == nil {
			update.VirtualMachineProperties.StorageProfile = &compute.StorageProfile{}
		}

		update.VirtualMachineProperties.StorageProfile.OsDisk = osDisk
	}

	if d.HasChange("virtual_machine_scale_set_id") {
		shouldUpdate = true

		if vmssIDRaw, ok := d.GetOk("virtual_machine_scale_set_id"); ok {
			update.VirtualMachineProperties.VirtualMachineScaleSet = &compute.SubResource{
				ID: utils.String(vmssIDRaw.(string)),
			}
		} else {
			update.VirtualMachineProperties.VirtualMachineScaleSet = &compute.SubResource{}
		}
	}

	if d.HasChange("proximity_placement_group_id") {
		shouldUpdate = true

		// Code="OperationNotAllowed" Message="Updating proximity placement group of VM is not allowed while the VM is running. Please stop/deallocate the VM and retry the operation."
		shouldShutDown = true
		shouldDeallocate = true

		if ppgIDRaw, ok := d.GetOk("proximity_placement_group_id"); ok {
			update.VirtualMachineProperties.ProximityPlacementGroup = &compute.SubResource{
				ID: utils.String(ppgIDRaw.(string)),
			}
		} else {
			update.VirtualMachineProperties.ProximityPlacementGroup = &compute.SubResource{}
		}
	}

	if d.HasChange("size") {
		shouldUpdate = true

		// this is kind of superflurious since Azure can do this for us, but if we do this we can subsequently
		// deallocate the VM to switch hosts if required
		shouldShutDown = true
		vmSize := d.Get("size").(string)

		// Azure will auto-reboot this for us, providing this machine will fit on this host
		// otherwise we need to shut down the VM to move it to another host to be able to use this size
		availableOnThisHost := false
		sizes, err := client.ListAvailableSizes(ctx, id.ResourceGroupName, id.VirtualMachineName)
		if err != nil {
			return fmt.Errorf("retrieving available sizes for Linux %s: %+v", id, err)
		}

		if sizes.Value != nil {
			for _, size := range *sizes.Value {
				if size.Name == nil {
					continue
				}

				if strings.EqualFold(*size.Name, vmSize) {
					availableOnThisHost = true
					break
				}
			}
		}

		if !availableOnThisHost {
			log.Printf("[DEBUG] Requested VM Size isn't available on the Host - must switch host to resize..")
			// Code="OperationNotAllowed"
			// Message="Unable to resize the VM [name] because the requested size Standard_F4s_v2 is not available in the current hardware cluster.
			//         The available sizes in this cluster are: [list]. The requested size might be available in other clusters of this region.
			//         Read more on VM resizing strategy at https://aka.ms/azure-resizevm."
			shouldDeallocate = true
		}

		update.VirtualMachineProperties.HardwareProfile = &compute.HardwareProfile{
			VMSize: compute.VirtualMachineSizeTypes(vmSize),
		}
	}

	if d.HasChange("vm_agent_platform_updates_enabled") {
		shouldUpdate = true
		if update.VirtualMachineProperties.OsProfile == nil {
			update.VirtualMachineProperties.OsProfile = &compute.OSProfile{}
		}

		if update.VirtualMachineProperties.OsProfile.LinuxConfiguration == nil {
			update.VirtualMachineProperties.OsProfile.LinuxConfiguration = &compute.LinuxConfiguration{}
		}

		update.VirtualMachineProperties.OsProfile.LinuxConfiguration.EnableVMAgentPlatformUpdates = utils.Bool(d.Get("vm_agent_platform_updates_enabled").(bool))
	}

	if d.HasChange("patch_mode") {
		shouldUpdate = true
		patchSettings := &compute.LinuxPatchSettings{}

		if patchMode, ok := d.GetOk("patch_mode"); ok {
			patchSettings.PatchMode = compute.LinuxVMGuestPatchMode(patchMode.(string))
		} else {
			patchSettings.PatchMode = compute.LinuxVMGuestPatchModeImageDefault
		}

		if update.VirtualMachineProperties.OsProfile == nil {
			update.VirtualMachineProperties.OsProfile = &compute.OSProfile{}
		}

		if update.VirtualMachineProperties.OsProfile.LinuxConfiguration == nil {
			update.VirtualMachineProperties.OsProfile.LinuxConfiguration = &compute.LinuxConfiguration{}
		}

		update.VirtualMachineProperties.OsProfile.LinuxConfiguration.PatchSettings = patchSettings
	}

	if d.HasChange("patch_assessment_mode") {
		assessmentMode := d.Get("patch_assessment_mode").(string)
		if assessmentMode == string(compute.LinuxPatchAssessmentModeAutomaticByPlatform) && !d.Get("provision_vm_agent").(bool) {
			return fmt.Errorf("`provision_vm_agent` must be set to `true` when `patch_assessment_mode` is set to `AutomaticByPlatform`")
		}

		shouldUpdate = true

		if update.VirtualMachineProperties.OsProfile == nil {
			update.VirtualMachineProperties.OsProfile = &compute.OSProfile{}
		}

		if update.VirtualMachineProperties.OsProfile.LinuxConfiguration == nil {
			update.VirtualMachineProperties.OsProfile.LinuxConfiguration = &compute.LinuxConfiguration{}
		}

		if update.VirtualMachineProperties.OsProfile.LinuxConfiguration.PatchSettings == nil {
			update.VirtualMachineProperties.OsProfile.LinuxConfiguration.PatchSettings = &compute.LinuxPatchSettings{}
		}

		update.VirtualMachineProperties.OsProfile.LinuxConfiguration.PatchSettings.AssessmentMode = compute.LinuxPatchAssessmentMode(assessmentMode)
	}

	isPatchModeAutomaticByPlatform := d.Get("patch_mode") == string(compute.LinuxVMGuestPatchModeAutomaticByPlatform)
	bypassPlatformSafetyChecksOnUserScheduleEnabled := d.Get("bypass_platform_safety_checks_on_user_schedule_enabled").(bool)
	if bypassPlatformSafetyChecksOnUserScheduleEnabled && !isPatchModeAutomaticByPlatform {
		return fmt.Errorf("`patch_mode` must be set to `AutomaticByPlatform` when `bypass_platform_safety_checks_on_user_schedule_enabled` is set to `true`")
	}
	if d.HasChange("bypass_platform_safety_checks_on_user_schedule_enabled") {
		shouldUpdate = true

		if update.VirtualMachineProperties.OsProfile == nil {
			update.VirtualMachineProperties.OsProfile = &compute.OSProfile{}
		}

		if update.VirtualMachineProperties.OsProfile.LinuxConfiguration == nil {
			update.VirtualMachineProperties.OsProfile.LinuxConfiguration = &compute.LinuxConfiguration{}
		}

		if update.VirtualMachineProperties.OsProfile.LinuxConfiguration.PatchSettings == nil {
			update.VirtualMachineProperties.OsProfile.LinuxConfiguration.PatchSettings = &compute.LinuxPatchSettings{}
		}

		if isPatchModeAutomaticByPlatform {
			if update.VirtualMachineProperties.OsProfile.LinuxConfiguration.PatchSettings.AutomaticByPlatformSettings == nil {
				update.VirtualMachineProperties.OsProfile.LinuxConfiguration.PatchSettings.AutomaticByPlatformSettings = &compute.LinuxVMGuestPatchAutomaticByPlatformSettings{}
			}

			update.VirtualMachineProperties.OsProfile.LinuxConfiguration.PatchSettings.AutomaticByPlatformSettings.BypassPlatformSafetyChecksOnUserSchedule = pointer.To(bypassPlatformSafetyChecksOnUserScheduleEnabled)
		}
	}

	rebootSetting := d.Get("reboot_setting").(string)
	if rebootSetting != "" && !isPatchModeAutomaticByPlatform {
		return fmt.Errorf("`patch_mode` must be set to `AutomaticByPlatform` when `reboot_setting` is specified")
	}
	if d.HasChange("reboot_setting") {
		shouldUpdate = true

		if update.VirtualMachineProperties.OsProfile == nil {
			update.VirtualMachineProperties.OsProfile = &compute.OSProfile{}
		}

		if update.VirtualMachineProperties.OsProfile.LinuxConfiguration == nil {
			update.VirtualMachineProperties.OsProfile.LinuxConfiguration = &compute.LinuxConfiguration{}
		}

		if update.VirtualMachineProperties.OsProfile.LinuxConfiguration.PatchSettings == nil {
			update.VirtualMachineProperties.OsProfile.LinuxConfiguration.PatchSettings = &compute.LinuxPatchSettings{}
		}

		if isPatchModeAutomaticByPlatform {
			if update.VirtualMachineProperties.OsProfile.LinuxConfiguration.PatchSettings.AutomaticByPlatformSettings == nil {
				update.VirtualMachineProperties.OsProfile.LinuxConfiguration.PatchSettings.AutomaticByPlatformSettings = &compute.LinuxVMGuestPatchAutomaticByPlatformSettings{}
			}

			update.VirtualMachineProperties.OsProfile.LinuxConfiguration.PatchSettings.AutomaticByPlatformSettings.RebootSetting = compute.LinuxVMGuestPatchAutomaticByPlatformRebootSetting(rebootSetting)
		}
	}

	if d.HasChange("allow_extension_operations") {
		allowExtensionOperations := d.Get("allow_extension_operations").(bool)

		shouldUpdate = true

		if update.OsProfile == nil {
			update.OsProfile = &compute.OSProfile{}
		}

		update.OsProfile.AllowExtensionOperations = utils.Bool(allowExtensionOperations)
	}

	var osImageNotificationProfile *compute.OSImageNotificationProfile
	var terminateNotificationProfile *compute.TerminateNotificationProfile

	if d.HasChange("os_image_notification") {
		shouldUpdate = true
		osImageNotificationProfile = expandOsImageNotificationProfile(d.Get("os_image_notification").([]interface{}))
	}

	if d.HasChange("termination_notification") {
		shouldUpdate = true
		terminateNotificationProfile = expandTerminateNotificationProfile(d.Get("termination_notification").([]interface{}))
	}

	if osImageNotificationProfile != nil || terminateNotificationProfile != nil {
		update.ScheduledEventsProfile = &compute.ScheduledEventsProfile{
			OsImageNotificationProfile:   osImageNotificationProfile,
			TerminateNotificationProfile: terminateNotificationProfile,
		}
	}

	if d.HasChange("tags") {
		shouldUpdate = true

		tagsRaw := d.Get("tags").(map[string]interface{})
		update.Tags = tags.Expand(tagsRaw)
	}

	if d.HasChange("additional_capabilities") {
		shouldUpdate = true

		n, _ := d.GetChange("additional_capabilities")
		if len(n.([]interface{})) == 0 || d.HasChange("additional_capabilities.0.ultra_ssd_enabled") {
			shouldShutDown = true
			shouldDeallocate = true
		}

		additionalCapabilitiesRaw := d.Get("additional_capabilities").([]interface{})
		update.VirtualMachineProperties.AdditionalCapabilities = expandVirtualMachineAdditionalCapabilities(additionalCapabilitiesRaw)
	}

	if d.HasChange("encryption_at_host_enabled") {
		if d.Get("encryption_at_host_enabled").(bool) {
			osDiskRaw := d.Get("os_disk").([]interface{})
			securityEncryptionType := osDiskRaw[0].(map[string]interface{})["security_encryption_type"].(string)
			if compute.SecurityEncryptionTypesDiskWithVMGuestState == compute.SecurityEncryptionTypes(securityEncryptionType) {
				return fmt.Errorf("`encryption_at_host_enabled` cannot be set to `true` when `os_disk.0.security_encryption_type` is set to `DiskWithVMGuestState`")
			}
		}

		shouldUpdate = true
		shouldDeallocate = true // API returns the following error if not deallocate: 'securityProfile.encryptionAtHost' can be updated only when VM is in deallocated state

		update.VirtualMachineProperties.SecurityProfile = &compute.SecurityProfile{
			EncryptionAtHost: utils.Bool(d.Get("encryption_at_host_enabled").(bool)),
		}
	}

	if d.HasChange("user_data") {
		shouldUpdate = true
		update.UserData = utils.String(d.Get("user_data").(string))
	}

	if instanceView.Statuses != nil {
		for _, status := range *instanceView.Statuses {
			if status.Code == nil {
				continue
			}

			// could also be the provisioning state which we're not bothered with here
			state := strings.ToLower(*status.Code)
			if !strings.HasPrefix(state, "powerstate/") {
				continue
			}

			state = strings.TrimPrefix(state, "powerstate/")
			switch strings.ToLower(state) {
			case "deallocated":
				// VM already deallocated, no shutdown and deallocation needed anymore
				shouldShutDown = false
				shouldDeallocate = false
			case "deallocating":
				// VM is deallocating
				// To make sure we do not start updating before this action has finished,
				// only skip the shutdown and send another deallocation request if shouldDeallocate == true
				shouldShutDown = false
			case "stopped":
				// VM already stopped, no shutdown needed anymore
				shouldShutDown = false
			}
		}
	}

	if shouldShutDown {
		log.Printf("[DEBUG] Shutting Down Linux %s", id)
		forceShutdown := false
		future, err := client.PowerOff(ctx, id.ResourceGroupName, id.VirtualMachineName, utils.Bool(forceShutdown))
		if err != nil {
			return fmt.Errorf("sending Power Off to Linux %s: %+v", id, err)
		}

		if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for Power Off of Linux %s: %+v", id, err)
		}

		log.Printf("[DEBUG] Shut Down Linux %s.", id)
	}

	if shouldDeallocate {
		if !hasEphemeralOSDisk {
			log.Printf("[DEBUG] Deallocating Linux %s", id)
			// Upgrade to 2021-07-01 added a hibernate parameter to this call defaulting to false
			future, err := client.Deallocate(ctx, id.ResourceGroupName, id.VirtualMachineName, utils.Bool(false))
			if err != nil {
				return fmt.Errorf("deallocating Linux %s: %+v", id, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for Deallocation of Linux %s: %+v", id, err)
			}

			log.Printf("[DEBUG] Deallocated Linux %s", id)
		} else {
			// Code="OperationNotAllowed" Message="Operation 'deallocate' is not supported for VMs or VM Scale Set instances using an ephemeral OS disk."
			log.Printf("[DEBUG] Skipping deallocation for Linux %s since cannot deallocate a Virtual Machine with an Ephemeral OS Disk", id)
		}
	}

	// now the VM's shutdown/deallocated we can update the disk which can't be done via the VM API:
	// Code="ResizeDiskError" Message="Managed disk resize via Virtual Machine [name] is not allowed. Please resize disk resource at [id]."
	// Portal: "Disks can be resized or account type changed only when they are unattached or the owner VM is deallocated."
	if d.HasChange("os_disk.0.disk_size_gb") {
		diskName := d.Get("os_disk.0.name").(string)
		newSize := d.Get("os_disk.0.disk_size_gb").(int)
		log.Printf("[DEBUG] Resizing OS Disk %q for Linux %s to %dGB..", diskName, id, newSize)

		disksClient := meta.(*clients.Client).Compute.DisksClient
		subscriptionId := meta.(*clients.Client).Account.SubscriptionId
		id := commonids.NewManagedDiskID(subscriptionId, id.ResourceGroupName, diskName)

		update := disks.DiskUpdate{
			Properties: &disks.DiskUpdateProperties{
				DiskSizeGB: utils.Int64(int64(newSize)),
			},
		}

		err := disksClient.UpdateThenPoll(ctx, id, update)
		if err != nil {
			return fmt.Errorf("resizing OS Disk %q for Linux Virtual Machine %q (Resource Group %q): %+v", diskName, id.DiskName, id.ResourceGroupName, err)
		}

		log.Printf("[DEBUG] Resized OS Disk %q for Linux Virtual Machine %q (Resource Group %q) to %dGB.", diskName, id.DiskName, id.ResourceGroupName, newSize)
	}

	if d.HasChange("os_disk.0.disk_encryption_set_id") {
		if diskEncryptionSetId := d.Get("os_disk.0.disk_encryption_set_id").(string); diskEncryptionSetId != "" {
			diskName := d.Get("os_disk.0.name").(string)
			log.Printf("[DEBUG] Updating encryption settings of OS Disk %q for Linux %s to %q..", diskName, id, diskEncryptionSetId)

			encryptionType, err := retrieveDiskEncryptionSetEncryptionType(ctx, meta.(*clients.Client).Compute.DiskEncryptionSetsClient, diskEncryptionSetId)
			if err != nil {
				return err
			}

			disksClient := meta.(*clients.Client).Compute.DisksClient
			subscriptionId := meta.(*clients.Client).Account.SubscriptionId
			id := commonids.NewManagedDiskID(subscriptionId, id.ResourceGroupName, diskName)

			update := disks.DiskUpdate{
				Properties: &disks.DiskUpdateProperties{
					Encryption: &disks.Encryption{
						Type:                encryptionType,
						DiskEncryptionSetId: utils.String(diskEncryptionSetId),
					},
				},
			}

			err = disksClient.UpdateThenPoll(ctx, id, update)
			if err != nil {
				return fmt.Errorf("updating encryption settings of OS Disk %q for Linux Virtual Machine %q (Resource Group %q): %+v", diskName, id.DiskName, id.ResourceGroupName, err)
			}

			log.Printf("[DEBUG] Updating encryption settings of OS Disk %q for Linux Virtual Machine %q (Resource Group %q) to %q.", diskName, id.DiskName, id.ResourceGroupName, diskEncryptionSetId)
		} else {
			return fmt.Errorf("once a customer-managed key is used, you can’t change the selection back to a platform-managed key")
		}
	}

	if shouldUpdate {
		log.Printf("[DEBUG] Updating Linux %s", id)
		future, err := client.Update(ctx, id.ResourceGroupName, id.VirtualMachineName, update)
		if err != nil {
			return fmt.Errorf("updating Linux %s: %+v", id, err)
		}

		if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for update of Linux %s: %+v", id, err)
		}

		log.Printf("[DEBUG] Updated Linux %s", id)
	}

	// if we've shut it down and it was turned off, let's boot it back up
	if shouldTurnBackOn && (shouldShutDown || shouldDeallocate) {
		log.Printf("[DEBUG] Starting Linux %s", id)
		future, err := client.Start(ctx, id.ResourceGroupName, id.VirtualMachineName)
		if err != nil {
			return fmt.Errorf("starting Linux %s: %+v", id, err)
		}

		if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for start of Linux %s: %+v", id, err)
		}

		log.Printf("[DEBUG] Started Linux %s", id)
	}

	return resourceLinuxVirtualMachineRead(d, meta)
}

func resourceLinuxVirtualMachineDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VMClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseVirtualMachineID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.VirtualMachineName, VirtualMachineResourceName)
	defer locks.UnlockByName(id.VirtualMachineName, VirtualMachineResourceName)

	log.Printf("[DEBUG] Retrieving Linux %s", id)
	existing, err := client.Get(ctx, id.ResourceGroupName, id.VirtualMachineName, "")
	if err != nil {
		if utils.ResponseWasNotFound(existing.Response) {
			return nil
		}

		return fmt.Errorf("retrieving Linux %s: %+v", id, err)
	}

	if !meta.(*clients.Client).Features.VirtualMachine.SkipShutdownAndForceDelete {
		// If the VM was in a Failed state we can skip powering off, since that'll fail
		if strings.EqualFold(*existing.ProvisioningState, "failed") {
			log.Printf("[DEBUG] Powering Off Linux Virtual Machine was skipped because the VM was in %q state %s", *existing.ProvisioningState, id)
		} else {
			// ISSUE: 4920
			// shutting down the Virtual Machine prior to removing it means users are no longer charged for some Azure resources
			// thus this can be a large cost-saving when deleting larger instances
			// https://docs.microsoft.com/en-us/azure/virtual-machines/states-lifecycle
			log.Printf("[DEBUG] Powering Off Linux %s", id)
			skipShutdown := !meta.(*clients.Client).Features.VirtualMachine.GracefulShutdown
			powerOffFuture, err := client.PowerOff(ctx, id.ResourceGroupName, id.VirtualMachineName, utils.Bool(skipShutdown))
			if err != nil {
				return fmt.Errorf("powering off Linux %s: %+v", id, err)
			}
			if err := powerOffFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for power off of Linux %s: %+v", id, err)
			}
			log.Printf("[DEBUG] Powered Off Linux %s", id)
		}
	}

	log.Printf("[DEBUG] Deleting Linux %s", id)

	// Force Delete is in an opt-in Preview and can only be specified (true/false) if the feature is enabled
	// as such we default this to `nil` which matches the previous behaviour (where this isn't sent) and
	// conditionally set this if required
	var forceDeletion *bool = nil
	if meta.(*clients.Client).Features.VirtualMachine.SkipShutdownAndForceDelete {
		forceDeletion = utils.Bool(true)
	}
	deleteFuture, err := client.Delete(ctx, id.ResourceGroupName, id.VirtualMachineName, forceDeletion)
	if err != nil {
		return fmt.Errorf("deleting Linux %s: %+v", id, err)
	}
	if err := deleteFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of Linux %s: %+v", id, err)
	}
	log.Printf("[DEBUG] Deleted Linux %s", id)

	deleteOSDisk := meta.(*clients.Client).Features.VirtualMachine.DeleteOSDiskOnDeletion
	if deleteOSDisk {
		log.Printf("[DEBUG] Deleting OS Disk from Linux %s", id)
		disksClient := meta.(*clients.Client).Compute.DisksClient
		managedDiskId := ""
		if props := existing.VirtualMachineProperties; props != nil && props.StorageProfile != nil && props.StorageProfile.OsDisk != nil {
			if disk := props.StorageProfile.OsDisk.ManagedDisk; disk != nil && disk.ID != nil {
				managedDiskId = *disk.ID
			}
		}

		if managedDiskId != "" {
			diskId, err := commonids.ParseManagedDiskID(managedDiskId)
			if err != nil {
				return err
			}

			if err := disksClient.DeleteThenPoll(ctx, *diskId); err != nil {
				return fmt.Errorf("deleting %s for Linux %s: %+v", diskId, id, err)
			}

			log.Printf("[DEBUG] Deleted %s for Linux %s", diskId, id)
		} else {
			log.Printf("[DEBUG] Skipping Deleting OS Disk from Linux %s - cannot determine OS Disk ID.", id)
		}
	} else {
		log.Printf("[DEBUG] Skipping Deleting OS Disk from Linux %s", id)
	}

	// Need to add a get and a state wait to avoid bug in network API where the attached disk(s) are not actually deleted
	// Service team indicated that we need to do a get after VM delete call returns to verify that the VM and all attached
	// disks have actually been deleted.

	log.Printf("[INFO] verifying Linux %s has been deleted", id)
	virtualMachine, err := client.Get(ctx, id.ResourceGroupName, id.VirtualMachineName, "")
	if err != nil && !utils.ResponseWasNotFound(virtualMachine.Response) {
		return fmt.Errorf("verifying Linux %s has been deleted: %+v", id, err)
	}

	if !utils.ResponseWasNotFound(virtualMachine.Response) {
		log.Printf("[INFO] %s still exists, waiting on vm to be deleted", id)

		deleteWait := &pluginsdk.StateChangeConf{
			Pending:    []string{"200"},
			Target:     []string{"404"},
			MinTimeout: 30 * time.Second,
			Timeout:    d.Timeout(pluginsdk.TimeoutDelete),
			Refresh: func() (interface{}, string, error) {
				log.Printf("[INFO] checking on state of Linux %s", id)
				resp, err := client.Get(ctx, id.ResourceGroupName, id.VirtualMachineName, "")
				if err != nil {
					if utils.ResponseWasNotFound(resp.Response) {
						return resp, strconv.Itoa(resp.StatusCode), nil
					}
					return nil, "nil", fmt.Errorf("polling for the status of Linux %s: %v", id, err)
				}
				return resp, strconv.Itoa(resp.StatusCode), nil
			},
		}

		if _, err := deleteWait.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("waiting for the deletion of Linux %s: %v", id, err)
		}
	}

	return nil
}
