package compute

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2021-07-01/compute"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/base64"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// TODO: confirm locking as appropriate

func resourceWindowsVirtualMachine() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceWindowsVirtualMachineCreate,
		Read:   resourceWindowsVirtualMachineRead,
		Update: resourceWindowsVirtualMachineUpdate,
		Delete: resourceWindowsVirtualMachineDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.VirtualMachineID(id)
			return err
		}, importVirtualMachine(compute.OperatingSystemTypesWindows, "azurerm_windows_virtual_machine")),

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

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			// Required
			"admin_password": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				Sensitive:        true,
				DiffSuppressFunc: adminPasswordDiffSuppressFunc,
			},

			"admin_username": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
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

			"additional_unattend_content": additionalUnattendContentSchema(),

			"allow_extension_operations": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"availability_set_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: computeValidate.AvailabilitySetID,
				// the Compute/VM API is broken and returns the Resource Group name in UPPERCASE :shrug:
				DiffSuppressFunc: suppress.CaseDifference,
				// TODO: raise a GH issue for the broken API
				// availability_set_id:                 "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG-200122113424880096/providers/Microsoft.Compute/availabilitySets/ACCTESTAVSET-200122113424880096" => "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG-200122113424880096/providers/Microsoft.Compute/availabilitySets/acctestavset-200122113424880096" (forces new resource)
				ConflictsWith: []string{
					"virtual_machine_scale_set_id",
					"zone",
				},
			},

			"boot_diagnostics": bootDiagnosticsSchema(),

			"computer_name": {
				Type:     pluginsdk.TypeString,
				Optional: true,

				// Computed since we reuse the VM name if one's not specified
				Computed: true,
				ForceNew: true,

				ValidateFunc: computeValidate.WindowsComputerNameFull,
			},

			"custom_data": base64.OptionalSchema(true),

			"dedicated_host_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: computeValidate.DedicatedHostID,
				// the Compute/VM API is broken and returns the Resource Group name in UPPERCASE :shrug:
				// same for `dedicated_host_group_id`
				DiffSuppressFunc: suppress.CaseDifference,
				ConflictsWith: []string{
					"dedicated_host_group_id",
				},
				// TODO: raise a GH issue for the broken API
				// /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/TOM-MANUAL/providers/Microsoft.Compute/hostGroups/tom-hostgroup/hosts/tom-manual-host
			},

			"dedicated_host_group_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: computeValidate.DedicatedHostGroupID,
				// the Compute/VM API is broken and returns the Resource Group name in UPPERCASE
				DiffSuppressFunc: suppress.CaseDifference,
				ConflictsWith: []string{
					"dedicated_host_id",
				},
			},

			"enable_automatic_updates": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true, // updating this is not allowed "Changing property 'windowsConfiguration.enableAutomaticUpdates' is not allowed." Target="windowsConfiguration.enableAutomaticUpdates"
				Default:  true,
			},

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
					// NOTE: whilst Delete is an option here, it's only applicable for VMSS
					string(compute.VirtualMachineEvictionPolicyTypesDeallocate),
				}, false),
			},

			"extensions_time_budget": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      "PT1H30M",
				ValidateFunc: azValidate.ISO8601DurationBetween("PT15M", "PT2H"),
			},

			"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

			"license_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"None",
					"Windows_Client",
					"Windows_Server",
				}, false),
				DiffSuppressFunc: func(_, old, new string, _ *pluginsdk.ResourceData) bool {
					if old == "None" && new == "" || old == "" && new == "None" {
						return true
					}

					return false
				},
			},

			"max_bid_price": {
				Type:         pluginsdk.TypeFloat,
				Optional:     true,
				Default:      -1,
				ValidateFunc: validation.FloatAtLeast(-1.0),
			},

			"patch_mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(compute.WindowsVMGuestPatchModeAutomaticByOS),
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.WindowsVMGuestPatchModeAutomaticByOS),
					string(compute.WindowsVMGuestPatchModeAutomaticByPlatform),
					string(compute.WindowsVMGuestPatchModeManual),
				}, false),
			},

			"hotpatching_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
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

			"proximity_placement_group_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: computeValidate.ProximityPlacementGroupID,
				// the Compute/VM API is broken and returns the Resource Group name in UPPERCASE :shrug:
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"secret": windowsSecretSchema(),

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
					computeValidate.ImageID,
					computeValidate.SharedImageID,
					computeValidate.SharedImageVersionID,
				),
			},

			"source_image_reference": sourceImageReferenceSchema(true),

			"tags": tags.Schema(),

			"timezone": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: computeValidate.VirtualMachineTimeZone(),
			},

			"virtual_machine_scale_set_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				ConflictsWith: []string{
					"availability_set_id",
				},
				ValidateFunc: computeValidate.VirtualMachineScaleSetID,
			},

			"platform_fault_domain": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      -1,
				ForceNew:     true,
				RequiredWith: []string{"virtual_machine_scale_set_id"},
				ValidateFunc: validation.IntAtLeast(-1),
			},

			"user_data": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsBase64,
			},

			"vtpm_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"winrm_listener": winRmListenerSchema(),

			"zone": func() *pluginsdk.Schema {
				if !features.ThreePointOhBeta() {
					return &pluginsdk.Schema{
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						// this has to be computed because when you are trying to assign this VM to a VMSS in VMO mode with zones,
						// the VMO mode VMSS will assign a zone for each of its instance.
						// and if the VMSS in not zonal, this value should be left empty
						Computed: true,
						ConflictsWith: []string{
							"availability_set_id",
						},
					}
				}

				return commonschema.ZoneSingleOptionalForceNew()
			}(),

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

func resourceWindowsVirtualMachineCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VMClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewVirtualMachineID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	locks.ByName(id.Name, virtualMachineResourceName)
	defer locks.UnlockByName(id.Name, virtualMachineResourceName)

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, compute.InstanceViewTypesUserData)
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("checking for existing Windows %s: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(resp.Response) {
		return tf.ImportAsExistsError("azurerm_windows_virtual_machine", *resp.ID)
	}

	additionalCapabilitiesRaw := d.Get("additional_capabilities").([]interface{})
	additionalCapabilities := expandVirtualMachineAdditionalCapabilities(additionalCapabilitiesRaw)

	additionalUnattendContentRaw := d.Get("additional_unattend_content").([]interface{})
	additionalUnattendContent := expandAdditionalUnattendContent(additionalUnattendContentRaw)

	adminPassword := d.Get("admin_password").(string)
	adminUsername := d.Get("admin_username").(string)
	allowExtensionOperations := d.Get("allow_extension_operations").(bool)

	bootDiagnosticsRaw := d.Get("boot_diagnostics").([]interface{})
	bootDiagnostics := expandBootDiagnostics(bootDiagnosticsRaw)

	var computerName string
	if v, ok := d.GetOk("computer_name"); ok && len(v.(string)) > 0 {
		computerName = v.(string)
	} else {
		_, errs := computeValidate.WindowsComputerNameFull(d.Get("name"), "computer_name")
		if len(errs) > 0 {
			return fmt.Errorf("unable to assume default computer name %s. Please adjust the %q, or specify an explicit %q", errs[0], "name", "computer_name")
		}
		computerName = id.Name
	}
	enableAutomaticUpdates := d.Get("enable_automatic_updates").(bool)
	location := azure.NormalizeLocation(d.Get("location").(string))

	identity, err := expandVirtualMachineIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	planRaw := d.Get("plan").([]interface{})
	plan := expandPlan(planRaw)

	priority := compute.VirtualMachinePriorityTypes(d.Get("priority").(string))
	provisionVMAgent := d.Get("provision_vm_agent").(bool)
	patchMode := d.Get("patch_mode").(string)
	hotPatch := d.Get("hotpatching_enabled").(bool)
	size := d.Get("size").(string)
	t := d.Get("tags").(map[string]interface{})

	networkInterfaceIdsRaw := d.Get("network_interface_ids").([]interface{})
	networkInterfaceIds := expandVirtualMachineNetworkInterfaceIDs(networkInterfaceIdsRaw)

	osDiskRaw := d.Get("os_disk").([]interface{})
	osDisk := expandVirtualMachineOSDisk(osDiskRaw, compute.OperatingSystemTypesWindows)

	secretsRaw := d.Get("secret").([]interface{})
	secrets := expandWindowsSecrets(secretsRaw)

	sourceImageReferenceRaw := d.Get("source_image_reference").([]interface{})
	sourceImageId := d.Get("source_image_id").(string)
	sourceImageReference, err := expandSourceImageReference(sourceImageReferenceRaw, sourceImageId)
	if err != nil {
		return err
	}

	winRmListenersRaw := d.Get("winrm_listener").(*pluginsdk.Set).List()
	winRmListeners := expandWinRMListener(winRmListenersRaw)

	params := compute.VirtualMachine{
		Name:     utils.String(id.Name),
		Location: utils.String(location),
		Identity: identity,
		Plan:     plan,
		VirtualMachineProperties: &compute.VirtualMachineProperties{
			HardwareProfile: &compute.HardwareProfile{
				VMSize: compute.VirtualMachineSizeTypes(size),
			},
			OsProfile: &compute.OSProfile{
				AdminPassword:            utils.String(adminPassword),
				AdminUsername:            utils.String(adminUsername),
				ComputerName:             utils.String(computerName),
				AllowExtensionOperations: utils.Bool(allowExtensionOperations),
				WindowsConfiguration: &compute.WindowsConfiguration{
					ProvisionVMAgent:       utils.Bool(provisionVMAgent),
					EnableAutomaticUpdates: utils.Bool(enableAutomaticUpdates),
					WinRM:                  winRmListeners,
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

	if !provisionVMAgent && allowExtensionOperations {
		return fmt.Errorf("`allow_extension_operations` cannot be set to `true` when `provision_vm_agent` is set to `false`")
	}

	if len(additionalUnattendContentRaw) > 0 {
		params.OsProfile.WindowsConfiguration.AdditionalUnattendContent = additionalUnattendContent
	}

	isHotpatchImage := isValidHotPatchSourceImageReference(sourceImageReferenceRaw, sourceImageId)

	// Validate VM Guest Patch Mode configuration
	if patchMode == string(compute.WindowsVMGuestPatchModeAutomaticByPlatform) && !provisionVMAgent {
		return fmt.Errorf("%q cannot be set to %q when %q is set to %q", "patch_mode", "AutomaticByPlatform", "provision_vm_agent", "false")
	}

	if isHotpatchImage && patchMode != string(compute.WindowsVMGuestPatchModeAutomaticByPlatform) {
		return fmt.Errorf("%q must always be set to %q when %q points to a hotpatch enabled image", "patch_mode", "AutomaticByPlatform", "source_image_reference")
	}

	// hot patching can only be enabled if the patch_mode is set to "AutomaticByPlatform"
	// and if the image reference is using one of the following skus:
	// 2022-datacenter-azure-edition-core or 2022-datacenter-azure-edition-core-smalldisk
	if hotPatch {
		if patchMode != string(compute.WindowsVMGuestPatchModeAutomaticByPlatform) {
			return fmt.Errorf("%q cannot be set to %q when %q is set to %q", "hotpatching_enabled", "true", "patch_mode", patchMode)
		}

		if !provisionVMAgent {
			return fmt.Errorf("%q cannot be set to %q when %q is set to %q", "hotpatching_enabled", "true", "provisionVMAgent", "false")
		}

		if !isHotpatchImage {
			if sourceImageId != "" {
				return fmt.Errorf("the %q field is not supported if referencing the image via the %q field", "hotpatching_enabled", "source_image_id")
			}

			return fmt.Errorf("%q is currently only supported on %q or %q image reference skus", "hotpatching_enabled", "2022-datacenter-azure-edition-core", "2022-datacenter-azure-edition-core-smalldisk")
		}
	}

	params.OsProfile.WindowsConfiguration.PatchSettings = &compute.PatchSettings{
		PatchMode:         compute.WindowsVMGuestPatchMode(patchMode),
		EnableHotpatching: utils.Bool(hotPatch),
	}

	if v, ok := d.GetOk("availability_set_id"); ok {
		params.AvailabilitySet = &compute.SubResource{
			ID: utils.String(v.(string)),
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

	if encryptionAtHostEnabled, ok := d.GetOk("encryption_at_host_enabled"); ok {
		if params.SecurityProfile == nil {
			params.SecurityProfile = &compute.SecurityProfile{}
		}
		params.SecurityProfile.EncryptionAtHost = utils.Bool(encryptionAtHostEnabled.(bool))
	}

	if securebootEnabled, ok := d.GetOk("secure_boot_enabled"); ok {
		if params.SecurityProfile == nil {
			params.SecurityProfile = &compute.SecurityProfile{}
		}

		if params.SecurityProfile.UefiSettings == nil {
			params.SecurityProfile.UefiSettings = &compute.UefiSettings{}
		}

		secureboot := d.Get("secure_boot_enabled").(bool)
		if secureboot {
			params.SecurityProfile.SecurityType = compute.SecurityTypesTrustedLaunch
		}

		params.SecurityProfile.UefiSettings.SecureBootEnabled = utils.Bool(securebootEnabled.(bool))
	}

	if vtpmEnabled, ok := d.GetOk("vtpm_enabled"); ok {
		if params.SecurityProfile == nil {
			params.SecurityProfile = &compute.SecurityProfile{}
		}

		if params.SecurityProfile.UefiSettings == nil {
			params.SecurityProfile.UefiSettings = &compute.UefiSettings{}
		}
		vtpm := d.Get("vtpm_enabled").(bool)
		if vtpm {
			params.SecurityProfile.SecurityType = compute.SecurityTypesTrustedLaunch
		}

		params.SecurityProfile.UefiSettings.VTpmEnabled = utils.Bool(vtpmEnabled.(bool))
	}

	if evictionPolicyRaw, ok := d.GetOk("eviction_policy"); ok {
		if params.Priority != compute.VirtualMachinePriorityTypesSpot {
			return fmt.Errorf("an `eviction_policy` can only be specified when `priority` is set to `Spot`")
		}

		params.EvictionPolicy = compute.VirtualMachineEvictionPolicyTypes(evictionPolicyRaw.(string))
	} else if priority == compute.VirtualMachinePriorityTypesSpot {
		return fmt.Errorf("an `eviction_policy` must be specified when `priority` is set to `Spot`")
	}

	if v, ok := d.GetOk("license_type"); ok {
		params.LicenseType = utils.String(v.(string))
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

	if v, ok := d.GetOk("timezone"); ok {
		params.VirtualMachineProperties.OsProfile.WindowsConfiguration.TimeZone = utils.String(v.(string))
	}

	if v, ok := d.GetOk("user_data"); ok {
		params.UserData = utils.String(v.(string))
	}

	if v, ok := d.GetOk("zone"); ok {
		params.Zones = &[]string{
			v.(string),
		}
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, params)
	if err != nil {
		return fmt.Errorf("creating Windows %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Windows %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceWindowsVirtualMachineRead(d, meta)
}

func resourceWindowsVirtualMachineRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VMClient
	disksClient := meta.(*clients.Client).Compute.DisksClient
	networkInterfacesClient := meta.(*clients.Client).Network.InterfacesClient
	publicIPAddressesClient := meta.(*clients.Client).Network.PublicIPsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualMachineID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, compute.InstanceViewTypesUserData)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Windows Virtual Machine %q was not found in Resource Group %q - removing from state!", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Windows Virtual Machine %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

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
		return fmt.Errorf("retrieving Windows Virtual Machine %q (Resource Group %q): `properties` was nil", id.Name, id.ResourceGroup)
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

	if err := d.Set("boot_diagnostics", flattenBootDiagnostics(props.DiagnosticsProfile)); err != nil {
		return fmt.Errorf("setting `boot_diagnostics`: %+v", err)
	}

	d.Set("eviction_policy", string(props.EvictionPolicy))
	if profile := props.HardwareProfile; profile != nil {
		d.Set("size", string(profile.VMSize))
	}
	d.Set("license_type", props.LicenseType)

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

		if config := profile.WindowsConfiguration; config != nil {
			if err := d.Set("additional_unattend_content", flattenAdditionalUnattendContent(config.AdditionalUnattendContent, d)); err != nil {
				return fmt.Errorf("setting `additional_unattend_content`: %+v", err)
			}

			d.Set("enable_automatic_updates", config.EnableAutomaticUpdates)

			d.Set("provision_vm_agent", config.ProvisionVMAgent)

			if patchSettings := config.PatchSettings; patchSettings != nil {
				d.Set("patch_mode", patchSettings.PatchMode)
				d.Set("hotpatching_enabled", patchSettings.EnableHotpatching)
			}

			d.Set("timezone", config.TimeZone)

			if err := d.Set("winrm_listener", flattenWinRMListener(config.WinRM)); err != nil {
				return fmt.Errorf("setting `winrm_listener`: %+v", err)
			}
		}

		if err := d.Set("secret", flattenWindowsSecrets(profile.Secrets)); err != nil {
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
		d.Set("source_image_id", storageImageId)

		if err := d.Set("source_image_reference", flattenSourceImageReference(profile.ImageReference)); err != nil {
			return fmt.Errorf("setting `source_image_reference`: %+v", err)
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

func resourceWindowsVirtualMachineUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VMClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualMachineID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.Name, virtualMachineResourceName)
	defer locks.UnlockByName(id.Name, virtualMachineResourceName)

	log.Printf("[DEBUG] Retrieving Windows Virtual Machine %q (Resource Group %q)..", id.Name, id.ResourceGroup)
	existing, err := client.Get(ctx, id.ResourceGroup, id.Name, compute.InstanceViewTypesUserData)
	if err != nil {
		if utils.ResponseWasNotFound(existing.Response) {
			return nil
		}

		return fmt.Errorf("retrieving Windows Virtual Machine %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	log.Printf("[DEBUG] Retrieving InstanceView for Windows Virtual Machine %q (Resource Group %q)..", id.Name, id.ResourceGroup)
	instanceView, err := client.InstanceView(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving InstanceView for Windows Virtual Machine %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
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
			profile.Secrets = expandWindowsSecrets(secretsRaw)
		}

		update.VirtualMachineProperties.OsProfile = &profile
	}

	if d.HasChange("allow_extension_operations") {
		allowExtensionOperations := d.Get("allow_extension_operations").(bool)

		shouldUpdate = true

		if update.OsProfile == nil {
			update.OsProfile = &compute.OSProfile{}
		}

		update.OsProfile.AllowExtensionOperations = utils.Bool(allowExtensionOperations)
	}

	if d.HasChange("patch_mode") {
		shouldUpdate = true

		if update.OsProfile == nil {
			update.OsProfile = &compute.OSProfile{}
		}

		if update.OsProfile.WindowsConfiguration == nil {
			update.OsProfile.WindowsConfiguration = &compute.WindowsConfiguration{}
		}

		if update.OsProfile.WindowsConfiguration.PatchSettings == nil {
			update.OsProfile.WindowsConfiguration.PatchSettings = &compute.PatchSettings{}
		}

		update.OsProfile.WindowsConfiguration.PatchSettings.PatchMode = compute.WindowsVMGuestPatchMode(d.Get("patch_mode").(string))
	}

	if d.HasChange("hotpatching_enabled") {
		shouldUpdate = true

		if update.OsProfile == nil {
			update.OsProfile = &compute.OSProfile{}
		}

		if update.OsProfile.WindowsConfiguration == nil {
			update.OsProfile.WindowsConfiguration = &compute.WindowsConfiguration{}
		}

		if update.OsProfile.WindowsConfiguration.PatchSettings == nil {
			update.OsProfile.WindowsConfiguration.PatchSettings = &compute.PatchSettings{}
		}

		update.OsProfile.WindowsConfiguration.PatchSettings.EnableHotpatching = utils.Bool(d.Get("hotpatching_enabled").(bool))
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

	if d.HasChange("os_disk") {
		shouldUpdate = true

		// Code="Conflict" Message="Disk resizing is allowed only when creating a VM or when the VM is deallocated." Target="disk.diskSizeGB"
		shouldShutDown = true
		shouldDeallocate = true

		osDiskRaw := d.Get("os_disk").([]interface{})
		osDisk := expandVirtualMachineOSDisk(osDiskRaw, compute.OperatingSystemTypesWindows)
		update.VirtualMachineProperties.StorageProfile = &compute.StorageProfile{
			OsDisk: osDisk,
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
		sizes, err := client.ListAvailableSizes(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return fmt.Errorf("retrieving available sizes for Windows Virtual Machine %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
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

	if d.HasChange("tags") {
		shouldUpdate = true

		tagsRaw := d.Get("tags").(map[string]interface{})
		update.Tags = tags.Expand(tagsRaw)
	}

	if d.HasChange("additional_capabilities") {
		shouldUpdate = true

		if d.HasChange("additional_capabilities.0.ultra_ssd_enabled") {
			shouldShutDown = true
			shouldDeallocate = true
		}

		additionalCapabilitiesRaw := d.Get("additional_capabilities").([]interface{})
		update.VirtualMachineProperties.AdditionalCapabilities = expandVirtualMachineAdditionalCapabilities(additionalCapabilitiesRaw)
	}

	if d.HasChange("encryption_at_host_enabled") {
		shouldUpdate = true
		shouldDeallocate = true // API returns the following error if not deallocate: 'securityProfile.encryptionAtHost' can be updated only when VM is in deallocated state
		if update.SecurityProfile == nil {
			update.SecurityProfile = &compute.SecurityProfile{}
		}
		update.SecurityProfile.EncryptionAtHost = utils.Bool(d.Get("encryption_at_host_enabled").(bool))
	}

	if d.HasChange("license_type") {
		shouldUpdate = true

		license := d.Get("license_type").(string)
		if license == "" {
			// Only for create no specification is possible in the API. API does not allow empty string in update.
			// So removing attribute license_type from Terraform configuration if it was set to value other than 'None' would lead to an endless loop in apply.
			// To allow updating in this case set value explicitly to 'None'.
			license = "None"
		}
		update.VirtualMachineProperties.LicenseType = &license
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
		log.Printf("[DEBUG] Shutting Down Windows Virtual Machine %q (Resource Group %q)..", id.Name, id.ResourceGroup)
		forceShutdown := false
		future, err := client.PowerOff(ctx, id.ResourceGroup, id.Name, utils.Bool(forceShutdown))
		if err != nil {
			return fmt.Errorf("sending Power Off to Windows Virtual Machine %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}

		if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for Power Off of Windows Virtual Machine %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}

		log.Printf("[DEBUG] Shut Down Windows Virtual Machine %q (Resource Group %q)..", id.Name, id.ResourceGroup)
	}

	if shouldDeallocate {
		if !hasEphemeralOSDisk {
			log.Printf("[DEBUG] Deallocating Windows Virtual Machine %q (Resource Group %q)..", id.Name, id.ResourceGroup)
			// Upgrading to the 2021-07-01 exposed a new hibernate parameter in the Deallocate method
			future, err := client.Deallocate(ctx, id.ResourceGroup, id.Name, utils.Bool(false))
			if err != nil {
				return fmt.Errorf("deallocating Windows Virtual Machine %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for Deallocation of Windows Virtual Machine %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
			}

			log.Printf("[DEBUG] Deallocated Windows Virtual Machine %q (Resource Group %q)..", id.Name, id.ResourceGroup)
		} else {
			// Code="OperationNotAllowed" Message="Operation 'deallocate' is not supported for VMs or VM Scale Set instances using an ephemeral OS disk."
			log.Printf("[DEBUG] Skipping deallocation for Windows Virtual Machine %q (Resource Group %q) since cannot deallocate a Virtual Machine with an Ephemeral OS Disk", id.Name, id.ResourceGroup)
		}
	}

	// now the VM's shutdown/deallocated we can update the disk which can't be done via the VM API:
	// Code="ResizeDiskError" Message="Managed disk resize via Virtual Machine [name] is not allowed. Please resize disk resource at [id]."
	// Portal: "Disks can be resized or account type changed only when they are unattached or the owner VM is deallocated."
	if d.HasChange("os_disk.0.disk_size_gb") {
		diskName := d.Get("os_disk.0.name").(string)
		newSize := d.Get("os_disk.0.disk_size_gb").(int)
		log.Printf("[DEBUG] Resizing OS Disk %q for Windows Virtual Machine %q (Resource Group %q) to %dGB..", diskName, id.Name, id.ResourceGroup, newSize)

		disksClient := meta.(*clients.Client).Compute.DisksClient

		update := compute.DiskUpdate{
			DiskUpdateProperties: &compute.DiskUpdateProperties{
				DiskSizeGB: utils.Int32(int32(newSize)),
			},
		}

		future, err := disksClient.Update(ctx, id.ResourceGroup, diskName, update)
		if err != nil {
			return fmt.Errorf("resizing OS Disk %q for Windows Virtual Machine %q (Resource Group %q): %+v", diskName, id.Name, id.ResourceGroup, err)
		}

		if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for resize of OS Disk %q for Windows Virtual Machine %q (Resource Group %q): %+v", diskName, id.Name, id.ResourceGroup, err)
		}

		log.Printf("[DEBUG] Resized OS Disk %q for Windows Virtual Machine %q (Resource Group %q) to %dGB.", diskName, id.Name, id.ResourceGroup, newSize)
	}

	if d.HasChange("os_disk.0.disk_encryption_set_id") {
		if diskEncryptionSetId := d.Get("os_disk.0.disk_encryption_set_id").(string); diskEncryptionSetId != "" {
			diskName := d.Get("os_disk.0.name").(string)
			log.Printf("[DEBUG] Updating encryption settings of OS Disk %q for Windows Virtual Machine %q (Resource Group %q) to %q..", diskName, id.Name, id.ResourceGroup, diskEncryptionSetId)

			encryptionType, err := retrieveDiskEncryptionSetEncryptionType(ctx, meta.(*clients.Client).Compute.DiskEncryptionSetsClient, diskEncryptionSetId)
			if err != nil {
				return err
			}

			disksClient := meta.(*clients.Client).Compute.DisksClient

			update := compute.DiskUpdate{
				DiskUpdateProperties: &compute.DiskUpdateProperties{
					Encryption: &compute.Encryption{
						Type:                *encryptionType,
						DiskEncryptionSetID: utils.String(diskEncryptionSetId),
					},
				},
			}

			future, err := disksClient.Update(ctx, id.ResourceGroup, diskName, update)
			if err != nil {
				return fmt.Errorf("updating encryption settings of OS Disk %q for Windows Virtual Machine %q (Resource Group %q): %+v", diskName, id.Name, id.ResourceGroup, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting to update encryption settings of OS Disk %q for Windows Virtual Machine %q (Resource Group %q): %+v", diskName, id.Name, id.ResourceGroup, err)
			}

			log.Printf("[DEBUG] Updating encryption settings of OS Disk %q for Windows Virtual Machine %q (Resource Group %q) to %q.", diskName, id.Name, id.ResourceGroup, diskEncryptionSetId)
		} else {
			return fmt.Errorf("once a customer-managed key is used, you can’t change the selection back to a platform-managed key")
		}
	}

	if shouldUpdate {
		log.Printf("[DEBUG] Updating Windows Virtual Machine %q (Resource Group %q)..", id.Name, id.ResourceGroup)
		future, err := client.Update(ctx, id.ResourceGroup, id.Name, update)
		if err != nil {
			return fmt.Errorf("updating Windows Virtual Machine %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}

		if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for update of Windows Virtual Machine %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}

		log.Printf("[DEBUG] Updated Windows Virtual Machine %q (Resource Group %q).", id.Name, id.ResourceGroup)
	}

	// if we've shut it down and it was turned off, let's boot it back up
	if shouldTurnBackOn && shouldShutDown {
		log.Printf("[DEBUG] Starting Windows Virtual Machine %q (Resource Group %q)..", id.Name, id.ResourceGroup)
		future, err := client.Start(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return fmt.Errorf("starting Windows Virtual Machine %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}

		if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for start of Windows Virtual Machine %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}

		log.Printf("[DEBUG] Started Windows Virtual Machine %q (Resource Group %q)..", id.Name, id.ResourceGroup)
	}

	return resourceWindowsVirtualMachineRead(d, meta)
}

func resourceWindowsVirtualMachineDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VMClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualMachineID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.Name, virtualMachineResourceName)
	defer locks.UnlockByName(id.Name, virtualMachineResourceName)

	log.Printf("[DEBUG] Retrieving Windows Virtual Machine %q (Resource Group %q)..", id.Name, id.ResourceGroup)
	existing, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(existing.Response) {
			return nil
		}

		return fmt.Errorf("retrieving Windows Virtual Machine %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if !meta.(*clients.Client).Features.VirtualMachine.SkipShutdownAndForceDelete {
		// If the VM was in a Failed state we can skip powering off, since that'll fail
		if strings.EqualFold(*existing.ProvisioningState, "failed") {
			log.Printf("[DEBUG] Powering Off Windows Virtual Machine was skipped because the VM was in %q state %q (Resource Group %q).", *existing.ProvisioningState, id.Name, id.ResourceGroup)
		} else {
			// ISSUE: 4920
			// shutting down the Virtual Machine prior to removing it means users are no longer charged for some Azure resources
			// thus this can be a large cost-saving when deleting larger instances
			// https://docs.microsoft.com/en-us/azure/virtual-machines/states-lifecycle
			log.Printf("[DEBUG] Powering Off Windows Virtual Machine %q (Resource Group %q)..", id.Name, id.ResourceGroup)
			skipShutdown := !meta.(*clients.Client).Features.VirtualMachine.GracefulShutdown
			powerOffFuture, err := client.PowerOff(ctx, id.ResourceGroup, id.Name, utils.Bool(skipShutdown))
			if err != nil {
				return fmt.Errorf("powering off Windows Virtual Machine %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
			}
			if err := powerOffFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for power off of Windows Virtual Machine %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
			}
			log.Printf("[DEBUG] Powered Off Windows Virtual Machine %q (Resource Group %q).", id.Name, id.ResourceGroup)
		}
	}

	log.Printf("[DEBUG] Deleting Windows Virtual Machine %q (Resource Group %q)..", id.Name, id.ResourceGroup)

	// Force Delete is in an opt-in Preview and can only be specified (true/false) if the feature is enabled
	// as such we default this to `nil` which matches the previous behaviour (where this isn't sent) and
	// conditionally set this if required
	var forceDeletion *bool = nil
	if meta.(*clients.Client).Features.VirtualMachine.SkipShutdownAndForceDelete {
		forceDeletion = utils.Bool(true)
	}
	deleteFuture, err := client.Delete(ctx, id.ResourceGroup, id.Name, forceDeletion)
	if err != nil {
		return fmt.Errorf("deleting Windows Virtual Machine %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	if err := deleteFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of Windows Virtual Machine %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	log.Printf("[DEBUG] Deleted Windows Virtual Machine %q (Resource Group %q).", id.Name, id.ResourceGroup)

	deleteOSDisk := meta.(*clients.Client).Features.VirtualMachine.DeleteOSDiskOnDeletion
	if deleteOSDisk {
		log.Printf("[DEBUG] Deleting OS Disk from Windows Virtual Machine %q (Resource Group %q)..", id.Name, id.ResourceGroup)
		disksClient := meta.(*clients.Client).Compute.DisksClient
		managedDiskId := ""
		if props := existing.VirtualMachineProperties; props != nil && props.StorageProfile != nil && props.StorageProfile.OsDisk != nil {
			if disk := props.StorageProfile.OsDisk.ManagedDisk; disk != nil && disk.ID != nil {
				managedDiskId = *disk.ID
			}
		}

		if managedDiskId != "" {
			diskId, err := parse.ManagedDiskID(managedDiskId)
			if err != nil {
				return err
			}

			diskDeleteFuture, err := disksClient.Delete(ctx, diskId.ResourceGroup, diskId.DiskName)
			if err != nil {
				if !response.WasNotFound(diskDeleteFuture.Response()) {
					return fmt.Errorf("deleting OS Disk %q (Resource Group %q) for Windows Virtual Machine %q (Resource Group %q): %+v", diskId.DiskName, diskId.ResourceGroup, id.Name, id.ResourceGroup, err)
				}
			}
			if !response.WasNotFound(diskDeleteFuture.Response()) {
				if err := diskDeleteFuture.WaitForCompletionRef(ctx, disksClient.Client); err != nil {
					return fmt.Errorf("OS Disk %q (Resource Group %q) for Windows Virtual Machine %q (Resource Group %q): %+v", diskId.DiskName, diskId.ResourceGroup, id.Name, id.ResourceGroup, err)
				}
			}

			log.Printf("[DEBUG] Deleted OS Disk from Windows Virtual Machine %q (Resource Group %q).", diskId.DiskName, diskId.ResourceGroup)
		} else {
			log.Printf("[DEBUG] Skipping Deleting OS Disk from Windows Virtual Machine %q (Resource Group %q) - cannot determine OS Disk ID.", id.Name, id.ResourceGroup)
		}
	} else {
		log.Printf("[DEBUG] Skipping Deleting OS Disk from Windows Virtual Machine %q (Resource Group %q)..", id.Name, id.ResourceGroup)
	}

	// Need to add a get and a state wait to avoid bug in network API where the attached disk(s) are not actually deleted
	// Service team indicated that we need to do a get after VM delete call returns to verify that the VM and all attached
	// disks have actually been deleted.

	log.Printf("[INFO] verifying Windows Virtual Machine %q has been deleted", id.Name)
	virtualMachine, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil && !utils.ResponseWasNotFound(virtualMachine.Response) {
		return fmt.Errorf("verifying Windows Virtual Machine %q (Resource Group %q) has been deleted: %+v", id.Name, id.ResourceGroup, err)
	}

	if !utils.ResponseWasNotFound(virtualMachine.Response) {
		log.Printf("[INFO] Windows Virtual Machine still exists, waiting on Windows Virtual Machine %q to be deleted", id.Name)

		deleteWait := &pluginsdk.StateChangeConf{
			Pending:    []string{"200"},
			Target:     []string{"404"},
			MinTimeout: 30 * time.Second,
			Timeout:    d.Timeout(pluginsdk.TimeoutDelete),
			Refresh: func() (interface{}, string, error) {
				log.Printf("[INFO] checking on state of Windows Virtual Machine %q", id.Name)
				resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
				if err != nil {
					if utils.ResponseWasNotFound(resp.Response) {
						return resp, strconv.Itoa(resp.StatusCode), nil
					}
					return nil, "nil", fmt.Errorf("polling for the status of Windows Virtual Machine %q (Resource Group %q): %v", id.Name, id.ResourceGroup, err)
				}
				return resp, strconv.Itoa(resp.StatusCode), nil
			},
		}

		if _, err := deleteWait.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("waiting for the deletion of Windows Virtual Machine %q (Resource Group %q): %v", id.Name, id.ResourceGroup, err)
		}
	}

	return nil
}
