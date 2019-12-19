package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	computeSvc "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/base64"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmWindowsVirtualMachineScaleSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmWindowsVirtualMachineScaleSetCreate,
		Read:   resourceArmWindowsVirtualMachineScaleSetRead,
		Update: resourceArmWindowsVirtualMachineScaleSetUpdate,
		Delete: resourceArmWindowsVirtualMachineScaleSetDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := computeSvc.ParseVirtualMachineScaleSetID(id)
			// TODO: (prior to Beta) look up the VM & confirm this is a Windows VMSS
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		// TODO: exposing requireGuestProvisionSignal once it's available
		// https://github.com/Azure/azure-rest-api-specs/pull/7246

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: computeSvc.ValidateWindowsName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			// Required
			"admin_username": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"admin_password": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Sensitive:    true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"network_interface": computeSvc.VirtualMachineScaleSetNetworkInterfaceSchema(),

			"os_disk": computeSvc.VirtualMachineScaleSetOSDiskSchema(),

			"instances": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},

			"sku": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			// Optional
			"additional_capabilities": computeSvc.VirtualMachineScaleSetAdditionalCapabilitiesSchema(),

			"additional_unattend_config": {
				Type:     schema.TypeList,
				Optional: true,
				// whilst the SDK supports updating, the API doesn't:
				//   Code="PropertyChangeNotAllowed"
				//   Message="Changing property 'windowsConfiguration.additionalUnattendContent' is not allowed."
				//   Target="windowsConfiguration.additionalUnattendContent
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content": {
							Type:      schema.TypeString,
							Required:  true,
							ForceNew:  true,
							Sensitive: true,
						},
						"setting": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(compute.AutoLogon),
								string(compute.FirstLogonCommands),
							}, false),
						},
					},
				},
			},

			"automatic_os_upgrade_policy": computeSvc.VirtualMachineScaleSetAutomatedOSUpgradePolicySchema(),

			"boot_diagnostics": computeSvc.VirtualMachineScaleSetBootDiagnosticsSchema(),

			"computer_name_prefix": {
				Type:     schema.TypeString,
				Optional: true,

				// Computed since we reuse the VM name if one's not specified
				Computed: true,
				ForceNew: true,
				// note: whilst the portal says 1-15 characters it seems to mirror the rules for the vm name
				// (e.g. 1-15 for Windows, 1-63 for Windows)
				ValidateFunc: computeSvc.ValidateWindowsName,
			},

			"custom_data": base64.OptionalSchema(),

			"data_disk": computeSvc.VirtualMachineScaleSetDataDiskSchema(),

			"do_not_run_extensions_on_overprovisioned_machines": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"enable_automatic_updates": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"eviction_policy": {
				// only applicable when `priority` is set to `Low`
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.Deallocate),
					string(compute.Delete),
				}, false),
			},

			"health_probe_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"identity": computeSvc.VirtualMachineScaleSetIdentitySchema(),

			"license_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Windows_Client",
					"Windows_Server",
				}, false),
			},

			"max_bid_price": {
				Type:     schema.TypeFloat,
				Optional: true,
				Default:  -1,
			},

			"overprovision": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"plan": computeSvc.PlanSchema(),

			"priority": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(compute.Regular),
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.Low),
					string(compute.Regular),
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
				ValidateFunc: azure.ValidateResourceID,
				// the Compute API is broken and returns the Resource Group name in UPPERCASE :shrug:
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"rolling_upgrade_policy": computeSvc.VirtualMachineScaleSetRollingUpgradePolicySchema(),

			"secret": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// whilst this isn't present in the nested object it's required when this is specified
						"key_vault_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},

						"certificate": {
							Type:     schema.TypeSet,
							Required: true,
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"store": {
										Type:     schema.TypeString,
										Required: true,
									},
									"url": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: azure.ValidateKeyVaultChildId,
									},
								},
							},
						},
					},
				},
			},

			"single_placement_group": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"source_image_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"source_image_reference": computeSvc.VirtualMachineScaleSetSourceImageReferenceSchema(),

			"tags": tags.Schema(),

			"terraform_should_roll_instances_when_required": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"timezone": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.VirtualMachineTimeZone(),
			},

			"upgrade_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(compute.Manual),
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.Automatic),
					string(compute.Manual),
					string(compute.Rolling),
				}, false),
			},

			"winrm_listener": {
				Type:     schema.TypeSet,
				Optional: true,
				// Whilst the SDK allows you to modify this, the API does not:
				//   Code="PropertyChangeNotAllowed"
				//   Message="Changing property 'windowsConfiguration.winRM.listeners' is not allowed."
				//   Target="windowsConfiguration.winRM.listeners"
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(compute.HTTP),
								string(compute.HTTPS),
							}, false),
						},

						"certificate_url": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: azure.ValidateKeyVaultChildId,
						},
					},
				},
			},

			"zone_balance": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},

			"zones": azure.SchemaZones(),

			// Computed
			"unique_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmWindowsVirtualMachineScaleSetCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VMScaleSetClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	if features.ShouldResourcesBeImported() {
		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error checking for existing Windows Virtual Machine Scale Set %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if !utils.ResponseWasNotFound(resp.Response) {
			return tf.ImportAsExistsError("azurerm_windows_virtual_machine_scale_set", *resp.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	additionalCapabilitiesRaw := d.Get("additional_capabilities").([]interface{})
	additionalCapabilities := computeSvc.ExpandVirtualMachineScaleSetAdditionalCapabilities(additionalCapabilitiesRaw)

	additionalUnattendConfigRaw := d.Get("additional_unattend_config").([]interface{})
	additionalUnattendConfig := expandWindowsVirtualMachineScaleSetAdditionalUnattendConfig(additionalUnattendConfigRaw)

	bootDiagnosticsRaw := d.Get("boot_diagnostics").([]interface{})
	bootDiagnostics := computeSvc.ExpandVirtualMachineScaleSetBootDiagnostics(bootDiagnosticsRaw)

	dataDisksRaw := d.Get("data_disk").([]interface{})
	dataDisks := computeSvc.ExpandVirtualMachineScaleSetDataDisk(dataDisksRaw)

	identityRaw := d.Get("identity").([]interface{})
	identity, err := computeSvc.ExpandVirtualMachineScaleSetIdentity(identityRaw)
	if err != nil {
		return fmt.Errorf("Error expanding `identity`: %+v", err)
	}

	networkInterfacesRaw := d.Get("network_interface").([]interface{})
	networkInterfaces, err := computeSvc.ExpandVirtualMachineScaleSetNetworkInterface(networkInterfacesRaw)
	if err != nil {
		return fmt.Errorf("Error expanding `network_interface`: %+v", err)
	}

	osDiskRaw := d.Get("os_disk").([]interface{})
	osDisk := computeSvc.ExpandVirtualMachineScaleSetOSDisk(osDiskRaw, compute.Windows)

	planRaw := d.Get("plan").([]interface{})
	plan := computeSvc.ExpandPlan(planRaw)

	sourceImageReferenceRaw := d.Get("source_image_reference").([]interface{})
	sourceImageId := d.Get("source_image_id").(string)
	sourceImageReference, err := computeSvc.ExpandVirtualMachineScaleSetSourceImageReference(sourceImageReferenceRaw, sourceImageId)
	if err != nil {
		return err
	}

	healthProbeId := d.Get("health_probe_id").(string)
	upgradeMode := compute.UpgradeMode(d.Get("upgrade_mode").(string))
	automaticOSUpgradePolicyRaw := d.Get("automatic_os_upgrade_policy").([]interface{})
	automaticOSUpgradePolicy := computeSvc.ExpandVirtualMachineScaleSetAutomaticUpgradePolicy(automaticOSUpgradePolicyRaw)
	rollingUpgradePolicyRaw := d.Get("rolling_upgrade_policy").([]interface{})
	rollingUpgradePolicy := computeSvc.ExpandVirtualMachineScaleSetRollingUpgradePolicy(rollingUpgradePolicyRaw)

	if upgradeMode != compute.Manual && healthProbeId == "" {
		return fmt.Errorf("`healthProbeId` must be set when `upgrade_mode` is set to %q", string(upgradeMode))
	}

	if upgradeMode != compute.Automatic && len(automaticOSUpgradePolicyRaw) > 0 {
		return fmt.Errorf("An `automatic_os_upgrade_policy` block cannot be specified when `upgrade_mode` is not set to `Automatic`")
	}
	if upgradeMode == compute.Automatic && len(automaticOSUpgradePolicyRaw) == 0 {
		return fmt.Errorf("An `automatic_os_upgrade_policy` block must be specified when `upgrade_mode` is set to `Automatic`")
	}

	shouldHaveRollingUpgradePolicy := upgradeMode == compute.Automatic || upgradeMode == compute.Rolling
	if !shouldHaveRollingUpgradePolicy && len(rollingUpgradePolicyRaw) > 0 {
		return fmt.Errorf("A `rolling_upgrade_policy` block cannot be specified when `upgrade_mode` is set to %q", string(upgradeMode))
	}
	if shouldHaveRollingUpgradePolicy && len(rollingUpgradePolicyRaw) == 0 {
		return fmt.Errorf("A `rolling_upgrade_policy` block must be specified when `upgrade_mode` is set to %q", string(upgradeMode))
	}

	winRmListenersRaw := d.Get("winrm_listener").(*schema.Set).List()
	winRmListeners := expandWindowsVirtualMachineScaleSetWinRMListeners(winRmListenersRaw)

	secretsRaw := d.Get("secret").([]interface{})
	secrets := expandWindowsVirtualMachineScaleSetSecrets(secretsRaw)

	zonesRaw := d.Get("zones").([]interface{})
	zones := azure.ExpandZones(zonesRaw)

	var computerNamePrefix string
	if v, ok := d.GetOk("computer_name_prefix"); ok && len(v.(string)) > 0 {
		computerNamePrefix = v.(string)
	} else {
		computerNamePrefix = name
	}

	networkProfile := &compute.VirtualMachineScaleSetNetworkProfile{
		NetworkInterfaceConfigurations: networkInterfaces,
	}
	if healthProbeId != "" {
		networkProfile.HealthProbe = &compute.APIEntityReference{
			ID: utils.String(healthProbeId),
		}
	}

	priority := compute.VirtualMachinePriorityTypes(d.Get("priority").(string))
	upgradePolicy := compute.UpgradePolicy{
		Mode:                     upgradeMode,
		AutomaticOSUpgradePolicy: automaticOSUpgradePolicy,
		RollingUpgradePolicy:     rollingUpgradePolicy,
	}

	virtualMachineProfile := compute.VirtualMachineScaleSetVMProfile{
		Priority: priority,
		OsProfile: &compute.VirtualMachineScaleSetOSProfile{
			AdminPassword:      utils.String(d.Get("admin_password").(string)),
			AdminUsername:      utils.String(d.Get("admin_username").(string)),
			ComputerNamePrefix: utils.String(computerNamePrefix),
			WindowsConfiguration: &compute.WindowsConfiguration{
				ProvisionVMAgent: utils.Bool(d.Get("provision_vm_agent").(bool)),
				WinRM:            winRmListeners,
			},
			Secrets: secrets,
		},
		DiagnosticsProfile: bootDiagnostics,
		NetworkProfile:     networkProfile,
		StorageProfile: &compute.VirtualMachineScaleSetStorageProfile{
			ImageReference: sourceImageReference,
			OsDisk:         osDisk,
			DataDisks:      dataDisks,
		},
	}

	enableAutomaticUpdates := d.Get("enable_automatic_updates").(bool)
	if upgradeMode != compute.Automatic {
		virtualMachineProfile.OsProfile.WindowsConfiguration.EnableAutomaticUpdates = utils.Bool(enableAutomaticUpdates)
	} else if !enableAutomaticUpdates {
		return fmt.Errorf("`enable_automatic_updates` must be set to `true` when `upgrade_mode` is set to `Automatic`")
	}

	if v, ok := d.Get("max_bid_price").(float64); ok && v > 0 {
		if priority != compute.Low {
			return fmt.Errorf("`max_bid_price` can only be configured when `priority` is set to `Low`")
		}

		virtualMachineProfile.BillingProfile = &compute.BillingProfile{
			MaxPrice: utils.Float(v),
		}
	}

	if v, ok := d.GetOk("custom_data"); ok {
		virtualMachineProfile.OsProfile.CustomData = utils.String(v.(string))
	}

	if evictionPolicyRaw, ok := d.GetOk("eviction_policy"); ok {
		if virtualMachineProfile.Priority != compute.Low {
			return fmt.Errorf("An `eviction_policy` can only be specified when `priority` is set to `low`")
		}
		virtualMachineProfile.EvictionPolicy = compute.VirtualMachineEvictionPolicyTypes(evictionPolicyRaw.(string))
	} else if priority == compute.Low {
		return fmt.Errorf("An `eviction_policy` must be specified when `priority` is set to `low`")
	}

	if len(additionalUnattendConfigRaw) > 0 {
		virtualMachineProfile.OsProfile.WindowsConfiguration.AdditionalUnattendContent = additionalUnattendConfig
	}

	if v, ok := d.GetOk("license_type"); ok {
		virtualMachineProfile.LicenseType = utils.String(v.(string))
	}

	if v, ok := d.GetOk("timezone"); ok {
		virtualMachineProfile.OsProfile.WindowsConfiguration.TimeZone = utils.String(v.(string))
	}

	props := compute.VirtualMachineScaleSet{
		Location: utils.String(location),
		Sku: &compute.Sku{
			Name:     utils.String(d.Get("sku").(string)),
			Capacity: utils.Int64(int64(d.Get("instances").(int))),

			// doesn't appear this can be set to anything else, even Promo machines are Standard
			Tier: utils.String("Standard"),
		},
		Identity: identity,
		Plan:     plan,
		Tags:     tags.Expand(t),
		VirtualMachineScaleSetProperties: &compute.VirtualMachineScaleSetProperties{
			AdditionalCapabilities:                 additionalCapabilities,
			DoNotRunExtensionsOnOverprovisionedVMs: utils.Bool(d.Get("do_not_run_extensions_on_overprovisioned_machines").(bool)),
			Overprovision:                          utils.Bool(d.Get("overprovision").(bool)),
			SinglePlacementGroup:                   utils.Bool(d.Get("single_placement_group").(bool)),
			VirtualMachineProfile:                  &virtualMachineProfile,
			UpgradePolicy:                          &upgradePolicy,
		},
		Zones: zones,
	}

	if v, ok := d.GetOk("proximity_placement_group_id"); ok {
		props.VirtualMachineScaleSetProperties.ProximityPlacementGroup = &compute.SubResource{
			ID: utils.String(v.(string)),
		}
	}

	if v, ok := d.GetOk("zone_balance"); ok && v.(bool) {
		if len(zonesRaw) == 0 {
			return fmt.Errorf("`zone_balance` can only be set to `true` when zones are specified")
		}

		props.VirtualMachineScaleSetProperties.ZoneBalance = utils.Bool(v.(bool))
	}

	log.Printf("[DEBUG] Creating Windows Virtual Machine Scale Set %q (Resource Group %q)..", name, resourceGroup)
	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, props)
	if err != nil {
		return fmt.Errorf("Error creating Windows Virtual Machine Scale Set %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	log.Printf("[DEBUG] Waiting for Windows Virtual Machine Scale Set %q (Resource Group %q) to be created..", name, resourceGroup)
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Windows Virtual Machine Scale Set %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	log.Printf("[DEBUG] Virtual Machine Scale Set %q (Resource Group %q) was created", name, resourceGroup)

	log.Printf("[DEBUG] Retrieving Virtual Machine Scale Set %q (Resource Group %q)..", name, resourceGroup)
	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Windows Virtual Machine Scale Set %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Error retrieving Windows Virtual Machine Scale Set %q (Resource Group %q): ID was nil", name, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmWindowsVirtualMachineScaleSetRead(d, meta)
}

func resourceArmWindowsVirtualMachineScaleSetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VMScaleSetClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := computeSvc.ParseVirtualMachineScaleSetID(d.Id())
	if err != nil {
		return err
	}

	updateInstances := false

	// retrieve
	existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("Error retrieving Windows Virtual Machine Scale Set %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	if existing.VirtualMachineScaleSetProperties == nil {
		return fmt.Errorf("Error retrieving Windows Virtual Machine Scale Set %q (Resource Group %q): `properties` was nil", id.Name, id.ResourceGroup)
	}

	updateProps := compute.VirtualMachineScaleSetUpdateProperties{
		VirtualMachineProfile: &compute.VirtualMachineScaleSetUpdateVMProfile{},
		// if an upgrade policy's been configured previously (which it will have) it must be threaded through
		// this doesn't matter for Manual - but breaks when updating anything on a Automatic and Rolling Mode Scale Set
		UpgradePolicy: existing.VirtualMachineScaleSetProperties.UpgradePolicy,
	}
	update := compute.VirtualMachineScaleSetUpdate{}

	upgradeMode := compute.UpgradeMode(d.Get("upgrade_mode").(string))
	if d.HasChange("automatic_os_upgrade_policy") || d.HasChange("rolling_upgrade_policy") {
		upgradePolicy := compute.UpgradePolicy{
			Mode: upgradeMode,
		}

		if d.HasChange("automatic_os_upgrade_policy") {
			automaticRaw := d.Get("automatic_os_upgrade_policy").([]interface{})
			upgradePolicy.AutomaticOSUpgradePolicy = computeSvc.ExpandVirtualMachineScaleSetAutomaticUpgradePolicy(automaticRaw)
		}

		if d.HasChange("rolling_upgrade_policy") {
			rollingRaw := d.Get("rolling_upgrade_policy").([]interface{})
			upgradePolicy.RollingUpgradePolicy = computeSvc.ExpandVirtualMachineScaleSetRollingUpgradePolicy(rollingRaw)
		}

		updateProps.UpgradePolicy = &upgradePolicy
	}

	priority := compute.VirtualMachinePriorityTypes(d.Get("priority").(string))
	if d.HasChange("max_bid_price") {
		if priority != compute.Low {
			return fmt.Errorf("`max_bid_price` can only be configured when `priority` is set to `Low`")
		}

		updateProps.VirtualMachineProfile.BillingProfile = &compute.BillingProfile{
			MaxPrice: utils.Float(d.Get("max_bid_price").(float64)),
		}
	}

	if d.HasChange("single_placement_group") {
		updateProps.SinglePlacementGroup = utils.Bool(d.Get("single_placement_group").(bool))
	}

	if d.HasChange("enable_automatic_updates") ||
		d.HasChange("custom_data") ||
		d.HasChange("provision_vm_agent") ||
		d.HasChange("secret") ||
		d.HasChange("timezone") {
		osProfile := compute.VirtualMachineScaleSetUpdateOSProfile{}

		if d.HasChange("enable_automatic_updates") || d.HasChange("provision_vm_agent") || d.HasChange("timezone") {
			windowsConfig := compute.WindowsConfiguration{}

			if d.HasChange("enable_automatic_updates") {
				if upgradeMode == compute.Automatic {
					return fmt.Errorf("`enable_automatic_updates` cannot be changed for when `upgrade_mode` is `Automatic`")
				}

				windowsConfig.EnableAutomaticUpdates = utils.Bool(d.Get("enable_automatic_updates").(bool))
			}

			if d.HasChange("provision_vm_agent") {
				windowsConfig.ProvisionVMAgent = utils.Bool(d.Get("provision_vm_agent").(bool))
			}

			if d.HasChange("timezone") {
				windowsConfig.TimeZone = utils.String(d.Get("timezone").(string))
			}

			osProfile.WindowsConfiguration = &windowsConfig
		}

		if d.HasChange("custom_data") {
			updateInstances = true

			// customData can only be sent if it's a base64 encoded string,
			// so it's not possible to remove this without tainting the resource
			if v, ok := d.GetOk("custom_data"); ok {
				osProfile.CustomData = utils.String(v.(string))
			}
		}

		if d.HasChange("secret") {
			secretsRaw := d.Get("secret").([]interface{})
			osProfile.Secrets = expandWindowsVirtualMachineScaleSetSecrets(secretsRaw)
		}

		updateProps.VirtualMachineProfile.OsProfile = &osProfile
	}

	if d.HasChange("data_disk") || d.HasChange("os_disk") || d.HasChange("source_image_id") || d.HasChange("source_image_reference") {
		updateInstances = true

		storageProfile := &compute.VirtualMachineScaleSetUpdateStorageProfile{}

		if d.HasChange("data_disk") {
			dataDisksRaw := d.Get("data_disk").([]interface{})
			storageProfile.DataDisks = computeSvc.ExpandVirtualMachineScaleSetDataDisk(dataDisksRaw)
		}

		if d.HasChange("os_disk") {
			osDiskRaw := d.Get("os_disk").([]interface{})
			storageProfile.OsDisk = computeSvc.ExpandVirtualMachineScaleSetOSDiskUpdate(osDiskRaw)
		}

		if d.HasChange("source_image_id") || d.HasChange("source_image_reference") {
			sourceImageReferenceRaw := d.Get("source_image_reference").([]interface{})
			sourceImageId := d.Get("source_image_id").(string)
			sourceImageReference, err := computeSvc.ExpandVirtualMachineScaleSetSourceImageReference(sourceImageReferenceRaw, sourceImageId)
			if err != nil {
				return err
			}

			storageProfile.ImageReference = sourceImageReference
		}

		updateProps.VirtualMachineProfile.StorageProfile = storageProfile
	}

	if d.HasChange("network_interface") {
		networkInterfacesRaw := d.Get("network_interface").([]interface{})
		networkInterfaces, err := computeSvc.ExpandVirtualMachineScaleSetNetworkInterfaceUpdate(networkInterfacesRaw)
		if err != nil {
			return fmt.Errorf("Error expanding `network_interface`: %+v", err)
		}

		// TODO: setting the health probe on update once https://github.com/Azure/azure-rest-api-specs/pull/7355 has been fixed
		//healthProbeId := d.Get("health_probe_id").(string)

		updateProps.VirtualMachineProfile.NetworkProfile = &compute.VirtualMachineScaleSetUpdateNetworkProfile{
			NetworkInterfaceConfigurations: networkInterfaces,
		}
	}

	if d.HasChange("boot_diagnostics") {
		updateInstances = true

		bootDiagnosticsRaw := d.Get("boot_diagnostics").([]interface{})
		updateProps.VirtualMachineProfile.DiagnosticsProfile = computeSvc.ExpandVirtualMachineScaleSetBootDiagnostics(bootDiagnosticsRaw)
	}

	if d.HasChange("identity") {
		identityRaw := d.Get("identity").([]interface{})
		identity, err := computeSvc.ExpandVirtualMachineScaleSetIdentity(identityRaw)
		if err != nil {
			return fmt.Errorf("Error expanding `identity`: %+v", err)
		}

		update.Identity = identity
	}

	if d.HasChange("plan") {
		planRaw := d.Get("plan").([]interface{})
		update.Plan = computeSvc.ExpandPlan(planRaw)
	}

	if d.HasChange("sku") || d.HasChange("instances") {
		// in-case ignore_changes is being used, since both fields are required
		// look up the current values and override them as needed
		sku := existing.Sku

		if d.HasChange("sku") {
			updateInstances = true

			sku.Name = utils.String(d.Get("sku").(string))
		}

		if d.HasChange("instances") {
			sku.Capacity = utils.Int64(int64(d.Get("instances").(int)))
		}

		update.Sku = sku
	}

	if d.HasChange("tags") {
		update.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	update.VirtualMachineScaleSetUpdateProperties = &updateProps

	log.Printf("[DEBUG] Updating Windows Virtual Machine Scale Set %q (Resource Group %q)..", id.Name, id.ResourceGroup)
	future, err := client.Update(ctx, id.ResourceGroup, id.Name, update)
	if err != nil {
		return fmt.Errorf("Error updating Windows Virtual Machine Scale Set %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	log.Printf("[DEBUG] Waiting for update of Windows Virtual Machine Scale Set %q (Resource Group %q)..", id.Name, id.ResourceGroup)
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for update of Windows Virtual Machine Scale Set %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	log.Printf("[DEBUG] Updated Windows Virtual Machine Scale Set %q (Resource Group %q).", id.Name, id.ResourceGroup)

	// if we update the SKU, we also need to subsequently roll the instances using the `UpdateInstances` API
	if updateInstances {
		if userWantsToRollInstances := d.Get("terraform_should_roll_instances_when_required").(bool); userWantsToRollInstances {
			log.Printf("[DEBUG] Rolling the VM Instances for Windows Virtual Machine Scale Set %q (Resource Group %q)..", id.Name, id.ResourceGroup)
			instancesClient := meta.(*clients.Client).Compute.VMScaleSetVMsClient
			instances, err := instancesClient.ListComplete(ctx, id.ResourceGroup, id.Name, "", "", "")
			if err != nil {
				return fmt.Errorf("Error listing VM Instances for Windows Virtual Machine Scale Set %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
			}

			log.Printf("[DEBUG] Determining instances to roll..")
			instanceIdsToRoll := make([]string, 0)
			for instances.NotDone() {
				instance := instances.Value()
				props := instance.VirtualMachineScaleSetVMProperties
				if props != nil && instance.InstanceID != nil {
					latestModel := props.LatestModelApplied
					if latestModel != nil || !*latestModel {
						instanceIdsToRoll = append(instanceIdsToRoll, *instance.InstanceID)
					}
				}

				if err := instances.NextWithContext(ctx); err != nil {
					return fmt.Errorf("Error enumerating instances: %s", err)
				}
			}

			// there's a performance enhancement to do batches here, but this is fine for a first pass
			for _, instanceId := range instanceIdsToRoll {
				instanceIds := []string{instanceId}

				log.Printf("[DEBUG] Updating Instance %q to the Latest Configuration..", instanceId)
				ids := compute.VirtualMachineScaleSetVMInstanceRequiredIDs{
					InstanceIds: &instanceIds,
				}
				future, err := client.UpdateInstances(ctx, id.ResourceGroup, id.Name, ids)
				if err != nil {
					return fmt.Errorf("Error updating Instance %q (Windows VM Scale Set %q / Resource Group %q) to the Latest Configuration: %+v", instanceId, id.Name, id.ResourceGroup, err)
				}

				if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
					return fmt.Errorf("Error waiting for update of Instance %q (Windows VM Scale Set %q / Resource Group %q) to the Latest Configuration: %+v", instanceId, id.Name, id.ResourceGroup, err)
				}
				log.Printf("[DEBUG] Updated Instance %q to the Latest Configuration.", instanceId)

				// TODO: does this want to be a separate, user-configurable toggle?
				log.Printf("[DEBUG] Reimaging Instance %q..", instanceId)
				reimageInput := &compute.VirtualMachineScaleSetReimageParameters{
					InstanceIds: &instanceIds,
				}
				reimageFuture, err := client.Reimage(ctx, id.ResourceGroup, id.Name, reimageInput)
				if err != nil {
					return fmt.Errorf("Error reimaging Instance %q (Windows VM Scale Set %q / Resource Group %q): %+v", instanceId, id.Name, id.ResourceGroup, err)
				}

				if err = reimageFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
					return fmt.Errorf("Error waiting for reimage of Instance %q (Windows VM Scale Set %q / Resource Group %q): %+v", instanceId, id.Name, id.ResourceGroup, err)
				}
				log.Printf("[DEBUG] Reimaged Instance %q..", instanceId)
			}

			log.Printf("[DEBUG] Rolled the VM Instances for Windows Virtual Machine Scale Set %q (Resource Group %q).", id.Name, id.ResourceGroup)
		} else {
			log.Printf("[DEBUG] Terraform wants to roll the VM Instances for Windows Virtual Machine Scale Set %q (Resource Group %q) - but user has opted out - skipping..", id.Name, id.ResourceGroup)
		}
	}

	return resourceArmWindowsVirtualMachineScaleSetRead(d, meta)
}

func resourceArmWindowsVirtualMachineScaleSetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VMScaleSetClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := computeSvc.ParseVirtualMachineScaleSetID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Windows Virtual Machine Scale Set %q was not found in Resource Group %q - removing from state!", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Windows Virtual Machine Scale Set %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	var skuName *string
	var instances int
	if resp.Sku != nil {
		skuName = resp.Sku.Name
		if resp.Sku.Capacity != nil {
			instances = int(*resp.Sku.Capacity)
		}
	}
	d.Set("instances", instances)
	d.Set("sku", skuName)

	if err := d.Set("identity", computeSvc.FlattenVirtualMachineScaleSetIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("Error setting `identity`: %+v", err)
	}

	if err := d.Set("plan", computeSvc.FlattenPlan(resp.Plan)); err != nil {
		return fmt.Errorf("Error setting `plan`: %+v", err)
	}

	if resp.VirtualMachineScaleSetProperties == nil {
		return fmt.Errorf("Error retrieving Windows Virtual Machine Scale Set %q (Resource Group %q): `properties` was nil", id.Name, id.ResourceGroup)
	}
	props := *resp.VirtualMachineScaleSetProperties

	if err := d.Set("additional_capabilities", computeSvc.FlattenVirtualMachineScaleSetAdditionalCapabilities(props.AdditionalCapabilities)); err != nil {
		return fmt.Errorf("Error setting `additional_capabilities`: %+v", props.AdditionalCapabilities)
	}

	d.Set("do_not_run_extensions_on_overprovisioned_machines", props.DoNotRunExtensionsOnOverprovisionedVMs)
	d.Set("overprovision", props.Overprovision)
	proximityPlacementGroupId := ""
	if props.ProximityPlacementGroup != nil && props.ProximityPlacementGroup.ID != nil {
		proximityPlacementGroupId = *props.ProximityPlacementGroup.ID
	}
	d.Set("proximity_placement_group_id", proximityPlacementGroupId)
	d.Set("single_placement_group", props.SinglePlacementGroup)
	d.Set("unique_id", props.UniqueID)
	d.Set("zone_balance", props.ZoneBalance)

	var upgradeMode compute.UpgradeMode
	if policy := props.UpgradePolicy; policy != nil {
		upgradeMode = policy.Mode
		d.Set("upgrade_mode", string(policy.Mode))

		flattenedAutomatic := computeSvc.FlattenVirtualMachineScaleSetAutomaticOSUpgradePolicy(policy.AutomaticOSUpgradePolicy)
		if err := d.Set("automatic_os_upgrade_policy", flattenedAutomatic); err != nil {
			return fmt.Errorf("Error setting `automatic_os_upgrade_policy`: %+v", err)
		}

		flattenedRolling := computeSvc.FlattenVirtualMachineScaleSetRollingUpgradePolicy(policy.RollingUpgradePolicy)
		if err := d.Set("rolling_upgrade_policy", flattenedRolling); err != nil {
			return fmt.Errorf("Error setting `rolling_upgrade_policy`: %+v", err)
		}
	}

	if profile := props.VirtualMachineProfile; profile != nil {
		if err := d.Set("boot_diagnostics", computeSvc.FlattenVirtualMachineScaleSetBootDiagnostics(profile.DiagnosticsProfile)); err != nil {
			return fmt.Errorf("Error setting `boot_diagnostics`: %+v", err)
		}

		// defaulted since BillingProfile isn't returned if it's unset
		maxBidPrice := float64(-1.0)
		if profile.BillingProfile != nil && profile.BillingProfile.MaxPrice != nil {
			maxBidPrice = *profile.BillingProfile.MaxPrice
		}
		d.Set("max_bid_price", maxBidPrice)

		d.Set("eviction_policy", string(profile.EvictionPolicy))
		d.Set("license_type", profile.LicenseType)
		d.Set("priority", string(profile.Priority))

		if storageProfile := profile.StorageProfile; storageProfile != nil {
			if err := d.Set("os_disk", computeSvc.FlattenVirtualMachineScaleSetOSDisk(storageProfile.OsDisk)); err != nil {
				return fmt.Errorf("Error setting `os_disk`: %+v", err)
			}

			if err := d.Set("data_disk", computeSvc.FlattenVirtualMachineScaleSetDataDisk(storageProfile.DataDisks)); err != nil {
				return fmt.Errorf("Error setting `data_disk`: %+v", err)
			}

			if err := d.Set("source_image_reference", computeSvc.FlattenVirtualMachineScaleSetSourceImageReference(storageProfile.ImageReference)); err != nil {
				return fmt.Errorf("Error setting `source_image_reference`: %+v", err)
			}

			var storageImageId string
			if storageProfile.ImageReference != nil && storageProfile.ImageReference.ID != nil {
				storageImageId = *storageProfile.ImageReference.ID
			}
			d.Set("source_image_id", storageImageId)
		}

		if osProfile := profile.OsProfile; osProfile != nil {
			// admin_password isn't returned, but it's a top level field so we can ignore it without consequence
			d.Set("admin_username", osProfile.AdminUsername)
			d.Set("computer_name_prefix", osProfile.ComputerNamePrefix)

			if err := d.Set("secret", flattenWindowsVirtualMachineScaleSetSecrets(osProfile.Secrets)); err != nil {
				return fmt.Errorf("Error setting `secret`: %+v", err)
			}

			if windows := osProfile.WindowsConfiguration; windows != nil {
				if err := d.Set("additional_unattend_config", flattenWindowsVirtualMachineScaleSetAdditionalUnattendConfig(windows.AdditionalUnattendContent, d)); err != nil {
					return fmt.Errorf("Error setting `additional_unattend_config`: %+v", err)
				}

				enableAutomaticUpdates := false
				if windows.EnableAutomaticUpdates != nil {
					enableAutomaticUpdates = *windows.EnableAutomaticUpdates
				}

				// the API requires this is set to 'true' on submission (since it's now required for Windows VMSS's with
				// an Automatic Upgrade Mode configured) however it actually returns false from the API..
				// after a bunch of testing the least bad option appears to be not to set this if it's an Automatic Upgrade Mode
				if upgradeMode != compute.Automatic {
					d.Set("enable_automatic_updates", enableAutomaticUpdates)
				}

				d.Set("provision_vm_agent", windows.ProvisionVMAgent)
				d.Set("timezone", windows.TimeZone)

				if err := d.Set("winrm_listener", flattenWindowsVirtualMachineScaleSetWinRMListener(windows.WinRM)); err != nil {
					return fmt.Errorf("Error setting `winrm_listener`: %+v", err)
				}
			}
		}

		if nwProfile := profile.NetworkProfile; nwProfile != nil {
			flattenedNics := computeSvc.FlattenVirtualMachineScaleSetNetworkInterface(nwProfile.NetworkInterfaceConfigurations)
			if err := d.Set("network_interface", flattenedNics); err != nil {
				return fmt.Errorf("Error setting `network_interface`: %+v", err)
			}

			healthProbeId := ""
			if nwProfile.HealthProbe != nil && nwProfile.HealthProbe.ID != nil {
				healthProbeId = *nwProfile.HealthProbe.ID
			}
			d.Set("health_probe_id", healthProbeId)
		}
	}

	if err := d.Set("zones", resp.Zones); err != nil {
		return fmt.Errorf("Error setting `zones`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmWindowsVirtualMachineScaleSetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VMScaleSetClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := computeSvc.ParseVirtualMachineScaleSetID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return nil
		}

		return fmt.Errorf("Error retrieving Windows Virtual Machine Scale Set %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	// Sometimes VMSS's aren't fully deleted when the `Delete` call returns - as such we'll try to scale the cluster
	// to 0 nodes first, then delete the cluster - which should ensure there's no Network Interfaces kicking around
	// and work around this Azure API bug:
	// Original Error: Code="InUseSubnetCannotBeDeleted" Message="Subnet internal is in use by
	// /{nicResourceID}/|providers|Microsoft.Compute|virtualMachineScaleSets|acctestvmss-190923101253410278|virtualMachines|0|networkInterfaces|example/ipConfigurations/internal and cannot be deleted.
	// In order to delete the subnet, delete all the resources within the subnet. See aka.ms/deletesubnet.
	if resp.Sku != nil {
		resp.Sku.Capacity = utils.Int64(int64(0))

		log.Printf("[DEBUG] Scaling instances to 0 prior to deletion - this helps avoids networking issues within Azure")
		update := compute.VirtualMachineScaleSetUpdate{
			Sku: resp.Sku,
		}
		future, err := client.Update(ctx, id.ResourceGroup, id.Name, update)
		if err != nil {
			return fmt.Errorf("Error updating number of instances in Windows Virtual Machine Scale Set %q (Resource Group %q) to scale to 0: %+v", id.Name, id.ResourceGroup, err)
		}

		log.Printf("[DEBUG] Waiting for scaling of instances to 0 prior to deletion - this helps avoids networking issues within Azure")
		err = future.WaitForCompletionRef(ctx, client.Client)
		if err != nil {
			return fmt.Errorf("Error waiting for number of instances in Windows Virtual Machine Scale Set %q (Resource Group %q) to scale to 0: %+v", id.Name, id.ResourceGroup, err)
		}
		log.Printf("[DEBUG] Scaled instances to 0 prior to deletion - this helps avoids networking issues within Azure")
	} else {
		log.Printf("[DEBUG] Unable to scale instances to `0` since the `sku` block is nil - trying to delete anyway")
	}

	log.Printf("[DEBUG] Deleting Windows Virtual Machine Scale Set %q (Resource Group %q)..", id.Name, id.ResourceGroup)
	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("Error deleting Windows Virtual Machine Scale Set %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	log.Printf("[DEBUG] Waiting for deletion of Windows Virtual Machine Scale Set %q (Resource Group %q)..", id.Name, id.ResourceGroup)
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of Windows Virtual Machine Scale Set %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	log.Printf("[DEBUG] Deleted Windows Virtual Machine Scale Set %q (Resource Group %q).", id.Name, id.ResourceGroup)

	return nil
}

func expandWindowsVirtualMachineScaleSetAdditionalUnattendConfig(input []interface{}) *[]compute.AdditionalUnattendContent {
	output := make([]compute.AdditionalUnattendContent, 0)

	for _, v := range input {
		raw := v.(map[string]interface{})

		output = append(output, compute.AdditionalUnattendContent{
			SettingName: compute.SettingNames(raw["setting"].(string)),
			Content:     utils.String(raw["content"].(string)),

			// no other possible values
			PassName:      compute.OobeSystem,
			ComponentName: compute.MicrosoftWindowsShellSetup,
		})
	}

	return &output
}

func flattenWindowsVirtualMachineScaleSetAdditionalUnattendConfig(input *[]compute.AdditionalUnattendContent, d *schema.ResourceData) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	existing := make([]interface{}, 0)
	if v, ok := d.GetOk("additional_unattend_config"); ok {
		existing = v.([]interface{})
	}

	output := make([]interface{}, 0)
	for i, v := range *input {
		// content isn't returned from the API as it's sensitive so we need to look it up
		content := ""
		if len(existing) > i {
			existingVal := existing[i]
			existingRaw, ok := existingVal.(map[string]interface{})
			if ok {
				contentRaw, ok := existingRaw["content"]
				if ok {
					content = contentRaw.(string)
				}
			}
		}

		output = append(output, map[string]interface{}{
			"content": content,
			"setting": string(v.SettingName),
		})
	}

	return output
}

func expandWindowsVirtualMachineScaleSetSecrets(input []interface{}) *[]compute.VaultSecretGroup {
	output := make([]compute.VaultSecretGroup, 0)

	for _, raw := range input {
		v := raw.(map[string]interface{})

		keyVaultId := v["key_vault_id"].(string)
		certificatesRaw := v["certificate"].(*schema.Set).List()
		certificates := make([]compute.VaultCertificate, 0)
		for _, certificateRaw := range certificatesRaw {
			certificateV := certificateRaw.(map[string]interface{})

			store := certificateV["store"].(string)
			url := certificateV["url"].(string)
			certificates = append(certificates, compute.VaultCertificate{
				CertificateStore: utils.String(store),
				CertificateURL:   utils.String(url),
			})
		}

		output = append(output, compute.VaultSecretGroup{
			SourceVault: &compute.SubResource{
				ID: utils.String(keyVaultId),
			},
			VaultCertificates: &certificates,
		})
	}

	return &output
}

func flattenWindowsVirtualMachineScaleSetSecrets(input *[]compute.VaultSecretGroup) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, v := range *input {
		keyVaultId := ""
		if v.SourceVault != nil && v.SourceVault.ID != nil {
			keyVaultId = *v.SourceVault.ID
		}

		certificates := make([]interface{}, 0)

		if v.VaultCertificates != nil {
			for _, c := range *v.VaultCertificates {
				store := ""
				if c.CertificateStore != nil {
					store = *c.CertificateStore
				}

				url := ""
				if c.CertificateURL != nil {
					url = *c.CertificateURL
				}

				certificates = append(certificates, map[string]interface{}{
					"store": store,
					"url":   url,
				})
			}
		}

		output = append(output, map[string]interface{}{
			"key_vault_id": keyVaultId,
			"certificate":  certificates,
		})
	}

	return output
}

func expandWindowsVirtualMachineScaleSetWinRMListeners(input []interface{}) *compute.WinRMConfiguration {
	listeners := make([]compute.WinRMListener, 0)

	for _, v := range input {
		raw := v.(map[string]interface{})

		listener := compute.WinRMListener{
			Protocol: compute.ProtocolTypes(raw["protocol"].(string)),
		}

		certificateUrl := raw["certificate_url"].(string)
		if certificateUrl != "" {
			listener.CertificateURL = utils.String(certificateUrl)
		}

		listeners = append(listeners, listener)
	}

	return &compute.WinRMConfiguration{
		Listeners: &listeners,
	}
}

func flattenWindowsVirtualMachineScaleSetWinRMListener(input *compute.WinRMConfiguration) []interface{} {
	if input == nil || input.Listeners == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, v := range *input.Listeners {
		certificateUrl := ""
		if v.CertificateURL != nil {
			certificateUrl = *v.CertificateURL
		}

		output = append(output, map[string]interface{}{
			"certificate_url": certificateUrl,
			"protocol":        string(v.Protocol),
		})
	}

	return output
}
