package compute

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	azValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
	computeValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/validate"
	networkValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/base64"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// TODO: confirm locking as appropriate

func resourceLinuxVirtualMachine() *schema.Resource {
	return &schema.Resource{
		Create: resourceLinuxVirtualMachineCreate,
		Read:   resourceLinuxVirtualMachineRead,
		Update: resourceLinuxVirtualMachineUpdate,
		Delete: resourceLinuxVirtualMachineDelete,
		Importer: azSchema.ValidateResourceIDPriorToImportThen(func(id string) error {
			_, err := parse.VirtualMachineID(id)
			return err
		}, importVirtualMachine(compute.Linux, "azurerm_linux_virtual_machine")),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(45 * time.Minute),
			Delete: schema.DefaultTimeout(45 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateVmName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			// Required
			"admin_username": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"network_interface_ids": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: networkValidate.NetworkInterfaceID,
				},
			},

			"os_disk": virtualMachineOSDiskSchema(),

			"size": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			// Optional
			"additional_capabilities": virtualMachineAdditionalCapabilitiesSchema(),

			"admin_password": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Sensitive:        true,
				DiffSuppressFunc: adminPasswordDiffSuppressFunc,
			},

			"admin_ssh_key": SSHKeysSchema(true),

			"allow_extension_operations": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"availability_set_id": {
				Type:         schema.TypeString,
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
				Type:     schema.TypeString,
				Optional: true,

				// Computed since we reuse the VM name if one's not specified
				Computed: true,
				ForceNew: true,

				ValidateFunc: ValidateLinuxComputerNameFull,
			},

			"custom_data": base64.OptionalSchema(true),

			"dedicated_host_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: computeValidate.DedicatedHostID,
				// the Compute/VM API is broken and returns the Resource Group name in UPPERCASE :shrug:
				DiffSuppressFunc: suppress.CaseDifference,
				// TODO: raise a GH issue for the broken API
				// /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/TOM-MANUAL/providers/Microsoft.Compute/hostGroups/tom-hostgroup/hosts/tom-manual-host
			},

			"disable_password_authentication": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},

			"encryption_at_host_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"eviction_policy": {
				// only applicable when `priority` is set to `Spot`
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					// NOTE: whilst Delete is an option here, it's only applicable for VMSS
					string(compute.Deallocate),
				}, false),
			},

			"extensions_time_budget": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "PT1H30M",
				ValidateFunc: azValidate.ISO8601DurationBetween("PT15M", "PT2H"),
			},

			"identity": virtualMachineIdentitySchema(),

			"max_bid_price": {
				Type:         schema.TypeFloat,
				Optional:     true,
				Default:      -1,
				ValidateFunc: validation.FloatAtLeast(-1.0),
			},

			"plan": planSchema(),

			"priority": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(compute.Regular),
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.Regular),
					string(compute.Spot),
				}, false),
			},

			"provision_vm_agent": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},

			"proximity_placement_group_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: computeValidate.ProximityPlacementGroupID,
				// the Compute/VM API is broken and returns the Resource Group name in UPPERCASE :shrug:
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"secret": linuxSecretSchema(),

			"source_image_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.Any(
					computeValidate.ImageID,
					computeValidate.SharedImageID,
					computeValidate.SharedImageVersionID,
				),
			},

			"source_image_reference": sourceImageReferenceSchema(true),

			"virtual_machine_scale_set_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ConflictsWith: []string{
					"availability_set_id",
				},
				ValidateFunc: computeValidate.VirtualMachineScaleSetID,
			},

			"tags": tags.Schema(),

			"zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				// this has to be computed because when you are trying to assign this VM to a VMSS in VMO mode with zones,
				// the VMO mode VMSS will assign a zone for each of its instance.
				// and if the VMSS in not zonal, this value should be left empty
				Computed: true,
				ConflictsWith: []string{
					"availability_set_id",
				},
			},

			// Computed
			"private_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"public_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"virtual_machine_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceLinuxVirtualMachineCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VMClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	locks.ByName(name, virtualMachineResourceName)
	defer locks.UnlockByName(name, virtualMachineResourceName)

	resp, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("checking for existing Linux Virtual Machine %q (Resource Group %q): %+v", name, resourceGroup, err)
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
		_, errs := ValidateLinuxComputerNameFull(d.Get("name"), "computer_name")
		if len(errs) > 0 {
			return fmt.Errorf("unable to assume default computer name %s Please adjust the %q, or specify an explicit %q", errs[0], "name", "computer_name")
		}
		computerName = name
	}
	disablePasswordAuthentication := d.Get("disable_password_authentication").(bool)
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
	osDisk := expandVirtualMachineOSDisk(osDiskRaw, compute.Linux)

	secretsRaw := d.Get("secret").([]interface{})
	secrets := expandLinuxSecrets(secretsRaw)

	sourceImageReferenceRaw := d.Get("source_image_reference").([]interface{})
	sourceImageId := d.Get("source_image_id").(string)
	sourceImageReference, err := expandSourceImageReference(sourceImageReferenceRaw, sourceImageId)
	if err != nil {
		return err
	}

	sshKeysRaw := d.Get("admin_ssh_key").(*schema.Set).List()
	sshKeys := ExpandSSHKeys(sshKeysRaw)

	params := compute.VirtualMachine{
		Name:     utils.String(name),
		Location: utils.String(location),
		Identity: identity,
		Plan:     plan,
		VirtualMachineProperties: &compute.VirtualMachineProperties{
			HardwareProfile: &compute.HardwareProfile{
				VMSize: compute.VirtualMachineSizeTypes(size),
			},
			OsProfile: &compute.OSProfile{
				AdminUsername:            utils.String(adminUsername),
				ComputerName:             utils.String(computerName),
				AllowExtensionOperations: utils.Bool(allowExtensionOperations),
				LinuxConfiguration: &compute.LinuxConfiguration{
					DisablePasswordAuthentication: utils.Bool(disablePasswordAuthentication),
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

	if encryptionAtHostEnabled, ok := d.GetOk("encryption_at_host_enabled"); ok {
		params.VirtualMachineProperties.SecurityProfile = &compute.SecurityProfile{
			EncryptionAtHost: utils.Bool(encryptionAtHostEnabled.(bool)),
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

	if v, ok := d.GetOk("custom_data"); ok {
		params.OsProfile.CustomData = utils.String(v.(string))
	}

	if v, ok := d.GetOk("dedicated_host_id"); ok {
		params.Host = &compute.SubResource{
			ID: utils.String(v.(string)),
		}
	}

	if evictionPolicyRaw, ok := d.GetOk("eviction_policy"); ok {
		if params.Priority != compute.Spot {
			return fmt.Errorf("An `eviction_policy` can only be specified when `priority` is set to `Spot`")
		}

		params.EvictionPolicy = compute.VirtualMachineEvictionPolicyTypes(evictionPolicyRaw.(string))
	} else if priority == compute.Spot {
		return fmt.Errorf("An `eviction_policy` must be specified when `priority` is set to `Spot`")
	}

	if v, ok := d.Get("max_bid_price").(float64); ok && v > 0 {
		if priority != compute.Spot {
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

	if v, ok := d.GetOk("zone"); ok {
		params.Zones = &[]string{
			v.(string),
		}
	}

	// "Authentication using either SSH or by user name and password must be enabled in Linux profile." Target="linuxConfiguration"
	adminPassword := d.Get("admin_password").(string)
	if disablePasswordAuthentication && len(sshKeys) == 0 {
		return fmt.Errorf("At least one `admin_ssh_key` must be specified when `disable_password_authentication` is set to `true`")
	} else if !disablePasswordAuthentication {
		if adminPassword == "" {
			return fmt.Errorf("An `admin_password` must be specified if `disable_password_authentication` is set to `false`")
		}

		params.OsProfile.AdminPassword = utils.String(adminPassword)
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, params)
	if err != nil {
		return fmt.Errorf("creating Linux Virtual Machine %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Linux Virtual Machine %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		return fmt.Errorf("retrieving Linux Virtual Machine %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("retrieving Linux Virtual Machine %q (Resource Group %q): `id` was nil", name, resourceGroup)
	}

	d.SetId(*read.ID)
	return resourceLinuxVirtualMachineRead(d, meta)
}

func resourceLinuxVirtualMachineRead(d *schema.ResourceData, meta interface{}) error {
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

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Linux Virtual Machine %q was not found in Resource Group %q - removing from state!", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Linux Virtual Machine %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if err := d.Set("identity", flattenVirtualMachineIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	if err := d.Set("plan", flattenPlan(resp.Plan)); err != nil {
		return fmt.Errorf("setting `plan`: %+v", err)
	}

	if resp.VirtualMachineProperties == nil {
		return fmt.Errorf("retrieving Linux Virtual Machine %q (Resource Group %q): `properties` was nil", id.Name, id.ResourceGroup)
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

	virtualMachineScaleSetId := ""
	if props.VirtualMachineScaleSet != nil && props.VirtualMachineScaleSet.ID != nil {
		virtualMachineScaleSetId = *props.VirtualMachineScaleSet.ID
	}
	d.Set("virtual_machine_scale_set_id", virtualMachineScaleSetId)

	if profile := props.OsProfile; profile != nil {
		d.Set("admin_username", profile.AdminUsername)
		d.Set("allow_extension_operations", profile.AllowExtensionOperations)
		d.Set("computer_name", profile.ComputerName)

		if config := profile.LinuxConfiguration; config != nil {
			d.Set("disable_password_authentication", config.DisablePasswordAuthentication)
			d.Set("provision_vm_agent", config.ProvisionVMAgent)

			flattenedSSHKeys, err := FlattenSSHKeys(config.SSH)
			if err != nil {
				return fmt.Errorf("flattening `admin_ssh_key`: %+v", err)
			}
			if err := d.Set("admin_ssh_key", flattenedSSHKeys); err != nil {
				return fmt.Errorf("setting `admin_ssh_key`: %+v", err)
			}
		}

		if err := d.Set("secret", flattenLinuxSecrets(profile.Secrets)); err != nil {
			return fmt.Errorf("setting `secret`: %+v", err)
		}
	}
	// Resources created with azurerm_virtual_machine have priority set to ""
	// We need to treat "" as equal to "Regular" to allow migration azurerm_virtual_machine -> azurerm_linux_virtual_machine
	priority := string(compute.Regular)
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
	if props.SecurityProfile != nil && props.SecurityProfile.EncryptionAtHost != nil {
		encryptionAtHostEnabled = *props.SecurityProfile.EncryptionAtHost
	}
	d.Set("encryption_at_host_enabled", encryptionAtHostEnabled)

	d.Set("virtual_machine_id", props.VMID)

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

func resourceLinuxVirtualMachineUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VMClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualMachineID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.Name, virtualMachineResourceName)
	defer locks.UnlockByName(id.Name, virtualMachineResourceName)

	log.Printf("[DEBUG] Retrieving Linux Virtual Machine %q (Resource Group %q)..", id.Name, id.ResourceGroup)
	existing, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(existing.Response) {
			return nil
		}

		return fmt.Errorf("retrieving Linux Virtual Machine %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	log.Printf("[DEBUG] Retrieving InstanceView for Linux Virtual Machine %q (Resource Group %q)..", id.Name, id.ResourceGroup)
	instanceView, err := client.InstanceView(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving InstanceView for Linux Virtual Machine %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	shouldTurnBackOn := virtualMachineShouldBeStarted(instanceView)
	hasEphemeralOSDisk := false
	if props := existing.VirtualMachineProperties; props != nil {
		if storage := props.StorageProfile; storage != nil {
			if disk := storage.OsDisk; disk != nil {
				if settings := disk.DiffDiskSettings; settings != nil {
					hasEphemeralOSDisk = settings.Option == compute.Local
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
		osDisk := expandVirtualMachineOSDisk(osDiskRaw, compute.Linux)
		update.VirtualMachineProperties.StorageProfile = &compute.StorageProfile{
			OsDisk: osDisk,
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
			return fmt.Errorf("retrieving available sizes for Linux Virtual Machine %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
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

	if d.HasChange("allow_extension_operations") {
		allowExtensionOperations := d.Get("allow_extension_operations").(bool)

		shouldUpdate = true

		if update.OsProfile == nil {
			update.OsProfile = &compute.OSProfile{}
		}

		update.OsProfile.AllowExtensionOperations = utils.Bool(allowExtensionOperations)
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

		update.VirtualMachineProperties.SecurityProfile = &compute.SecurityProfile{
			EncryptionAtHost: utils.Bool(d.Get("encryption_at_host_enabled").(bool)),
		}
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
		log.Printf("[DEBUG] Shutting Down Linux Virtual Machine %q (Resource Group %q)..", id.Name, id.ResourceGroup)
		forceShutdown := false
		future, err := client.PowerOff(ctx, id.ResourceGroup, id.Name, utils.Bool(forceShutdown))
		if err != nil {
			return fmt.Errorf("sending Power Off to Linux Virtual Machine %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}

		if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for Power Off of Linux Virtual Machine %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}

		log.Printf("[DEBUG] Shut Down Linux Virtual Machine %q (Resource Group %q)..", id.Name, id.ResourceGroup)
	}

	if shouldDeallocate {
		if !hasEphemeralOSDisk {
			log.Printf("[DEBUG] Deallocating Linux Virtual Machine %q (Resource Group %q)..", id.Name, id.ResourceGroup)
			future, err := client.Deallocate(ctx, id.ResourceGroup, id.Name)
			if err != nil {
				return fmt.Errorf("Deallocating Linux Virtual Machine %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for Deallocation of Linux Virtual Machine %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
			}

			log.Printf("[DEBUG] Deallocated Linux Virtual Machine %q (Resource Group %q)..", id.Name, id.ResourceGroup)
		} else {
			// Code="OperationNotAllowed" Message="Operation 'deallocate' is not supported for VMs or VM Scale Set instances using an ephemeral OS disk."
			log.Printf("[DEBUG] Skipping deallocation for Linux Virtual Machine %q (Resource Group %q) since cannot deallocate a Virtual Machine with an Ephemeral OS Disk", id.Name, id.ResourceGroup)
		}
	}

	// now the VM's shutdown/deallocated we can update the disk which can't be done via the VM API:
	// Code="ResizeDiskError" Message="Managed disk resize via Virtual Machine [name] is not allowed. Please resize disk resource at [id]."
	// Portal: "Disks can be resized or account type changed only when they are unattached or the owner VM is deallocated."
	if d.HasChange("os_disk.0.disk_size_gb") {
		diskName := d.Get("os_disk.0.name").(string)
		newSize := d.Get("os_disk.0.disk_size_gb").(int)
		log.Printf("[DEBUG] Resizing OS Disk %q for Linux Virtual Machine %q (Resource Group %q) to %dGB..", diskName, id.Name, id.ResourceGroup, newSize)

		disksClient := meta.(*clients.Client).Compute.DisksClient

		update := compute.DiskUpdate{
			DiskUpdateProperties: &compute.DiskUpdateProperties{
				DiskSizeGB: utils.Int32(int32(newSize)),
			},
		}

		future, err := disksClient.Update(ctx, id.ResourceGroup, diskName, update)
		if err != nil {
			return fmt.Errorf("resizing OS Disk %q for Linux Virtual Machine %q (Resource Group %q): %+v", diskName, id.Name, id.ResourceGroup, err)
		}

		if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for resize of OS Disk %q for Linux Virtual Machine %q (Resource Group %q): %+v", diskName, id.Name, id.ResourceGroup, err)
		}

		log.Printf("[DEBUG] Resized OS Disk %q for Linux Virtual Machine %q (Resource Group %q) to %dGB.", diskName, id.Name, id.ResourceGroup, newSize)
	}

	if d.HasChange("os_disk.0.disk_encryption_set_id") {
		if diskEncryptionSetId := d.Get("os_disk.0.disk_encryption_set_id").(string); diskEncryptionSetId != "" {
			diskName := d.Get("os_disk.0.name").(string)
			log.Printf("[DEBUG] Updating encryption settings of OS Disk %q for Linux Virtual Machine %q (Resource Group %q) to %q..", diskName, id.Name, id.ResourceGroup, diskEncryptionSetId)

			disksClient := meta.(*clients.Client).Compute.DisksClient

			update := compute.DiskUpdate{
				DiskUpdateProperties: &compute.DiskUpdateProperties{
					Encryption: &compute.Encryption{
						Type:                compute.EncryptionAtRestWithCustomerKey,
						DiskEncryptionSetID: utils.String(diskEncryptionSetId),
					},
				},
			}

			future, err := disksClient.Update(ctx, id.ResourceGroup, diskName, update)
			if err != nil {
				return fmt.Errorf("updating encryption settings of OS Disk %q for Linux Virtual Machine %q (Resource Group %q): %+v", diskName, id.Name, id.ResourceGroup, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting to update encryption settings of OS Disk %q for Linux Virtual Machine %q (Resource Group %q): %+v", diskName, id.Name, id.ResourceGroup, err)
			}

			log.Printf("[DEBUG] Updating encryption settings of OS Disk %q for Linux Virtual Machine %q (Resource Group %q) to %q.", diskName, id.Name, id.ResourceGroup, diskEncryptionSetId)
		} else {
			return fmt.Errorf("Once a customer-managed key is used, you can’t change the selection back to a platform-managed key")
		}
	}

	if shouldUpdate {
		log.Printf("[DEBUG] Updating Linux Virtual Machine %q (Resource Group %q)..", id.Name, id.ResourceGroup)
		future, err := client.Update(ctx, id.ResourceGroup, id.Name, update)
		if err != nil {
			return fmt.Errorf("updating Linux Virtual Machine %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}

		if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for update of Linux Virtual Machine %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}

		log.Printf("[DEBUG] Updated Linux Virtual Machine %q (Resource Group %q).", id.Name, id.ResourceGroup)
	}

	// if we've shut it down and it was turned off, let's boot it back up
	if shouldTurnBackOn && shouldShutDown {
		log.Printf("[DEBUG] Starting Linux Virtual Machine %q (Resource Group %q)..", id.Name, id.ResourceGroup)
		future, err := client.Start(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return fmt.Errorf("starting Linux Virtual Machine %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}

		if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for start of Linux Virtual Machine %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}

		log.Printf("[DEBUG] Started Linux Virtual Machine %q (Resource Group %q)..", id.Name, id.ResourceGroup)
	}

	return resourceLinuxVirtualMachineRead(d, meta)
}

func resourceLinuxVirtualMachineDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VMClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualMachineID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.Name, virtualMachineResourceName)
	defer locks.UnlockByName(id.Name, virtualMachineResourceName)

	log.Printf("[DEBUG] Retrieving Linux Virtual Machine %q (Resource Group %q)..", id.Name, id.ResourceGroup)
	existing, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(existing.Response) {
			return nil
		}

		return fmt.Errorf("retrieving Linux Virtual Machine %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	// ISSUE: XXX
	// shutting down the Virtual Machine prior to removing it means users are no longer charged for the compute
	// thus this can be a large cost-saving when deleting larger instances
	log.Printf("[DEBUG] Powering Off Linux Virtual Machine %q (Resource Group %q)..", id.Name, id.ResourceGroup)
	skipShutdown := !meta.(*clients.Client).Features.VirtualMachine.GracefulShutdown
	powerOffFuture, err := client.PowerOff(ctx, id.ResourceGroup, id.Name, utils.Bool(skipShutdown))
	if err != nil {
		return fmt.Errorf("powering off Linux Virtual Machine %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	if err := powerOffFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for power off of Linux Virtual Machine %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	log.Printf("[DEBUG] Powered Off Linux Virtual Machine %q (Resource Group %q).", id.Name, id.ResourceGroup)

	log.Printf("[DEBUG] Deleting Linux Virtual Machine %q (Resource Group %q)..", id.Name, id.ResourceGroup)
	deleteFuture, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Linux Virtual Machine %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	if err := deleteFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of Linux Virtual Machine %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	log.Printf("[DEBUG] Deleted Linux Virtual Machine %q (Resource Group %q).", id.Name, id.ResourceGroup)

	deleteOSDisk := meta.(*clients.Client).Features.VirtualMachine.DeleteOSDiskOnDeletion
	if deleteOSDisk {
		log.Printf("[DEBUG] Deleting OS Disk from Linux Virtual Machine %q (Resource Group %q)..", id.Name, id.ResourceGroup)
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
					return fmt.Errorf("deleting OS Disk %q (Resource Group %q) for Linux Virtual Machine %q (Resource Group %q): %+v", diskId.DiskName, diskId.ResourceGroup, id.Name, id.ResourceGroup, err)
				}
			}
			if !response.WasNotFound(diskDeleteFuture.Response()) {
				if err := diskDeleteFuture.WaitForCompletionRef(ctx, disksClient.Client); err != nil {
					return fmt.Errorf("OS Disk %q (Resource Group %q) for Linux Virtual Machine %q (Resource Group %q): %+v", diskId.DiskName, diskId.ResourceGroup, id.Name, id.ResourceGroup, err)
				}
			}

			log.Printf("[DEBUG] Deleted OS Disk from Linux Virtual Machine %q (Resource Group %q).", diskId.DiskName, diskId.ResourceGroup)
		} else {
			log.Printf("[DEBUG] Skipping Deleting OS Disk from Linux Virtual Machine %q (Resource Group %q) - cannot determine OS Disk ID.", id.Name, id.ResourceGroup)
		}
	} else {
		log.Printf("[DEBUG] Skipping Deleting OS Disk from Linux Virtual Machine %q (Resource Group %q)..", id.Name, id.ResourceGroup)
	}

	// Need to add a get and a state wait to avoid bug in network API where the attached disk(s) are not actually deleted
	// Service team indicated that we need to do a get after VM delete call returns to verify that the VM and all attached
	// disks have actually been deleted.

	log.Printf("[INFO] verifying Linux Virtual Machine %q has been deleted", id.Name)
	virtualMachine, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil && !utils.ResponseWasNotFound(virtualMachine.Response) {
		return fmt.Errorf("verifying Linux Virtual Machine %q (Resource Group %q) has been deleted: %+v", id.Name, id.ResourceGroup, err)
	}

	if !utils.ResponseWasNotFound(virtualMachine.Response) {
		log.Printf("[INFO] Linux Virtual Machine still exists, waiting on Linux Virtual Machine %q to be deleted", id.Name)

		deleteWait := &resource.StateChangeConf{
			Pending:    []string{"200"},
			Target:     []string{"404"},
			MinTimeout: 30 * time.Second,
			Timeout:    d.Timeout(schema.TimeoutDelete),
			Refresh: func() (interface{}, string, error) {
				log.Printf("[INFO] checking on state of Linux Virtual Machine %q", id.Name)
				resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
				if err != nil {
					if utils.ResponseWasNotFound(resp.Response) {
						return resp, strconv.Itoa(resp.StatusCode), nil
					}
					return nil, "nil", fmt.Errorf("polling for the status of Linux Virtual Machine %q (Resource Group %q): %v", id.Name, id.ResourceGroup, err)
				}
				return resp, strconv.Itoa(resp.StatusCode), nil
			},
		}

		if _, err := deleteWait.WaitForState(); err != nil {
			return fmt.Errorf("waiting for the deletion of Linux Virtual Machine %q (Resource Group %q): %v", id.Name, id.ResourceGroup, err)
		}
	}

	return nil
}
