package compute

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmVirtualMachineExtension() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmVirtualMachineExtensionsCreateUpdate,
		Read:   resourceArmVirtualMachineExtensionsRead,
		Update: resourceArmVirtualMachineExtensionsCreateUpdate,
		Delete: resourceArmVirtualMachineExtensionsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			// TODO: Remove in 2.0
			"virtual_machine_name": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				ForceNew:   true,
				Deprecated: "This field has been deprecated in favor of `virtual_machine_id` - will be removed in 2.0 of the Azure Provider",
			},

			"virtual_machine_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: ValidateVirtualMachineID,
			},

			"publisher": {
				Type:     schema.TypeString,
				Required: true,
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
			},

			"type_handler_version": {
				Type:     schema.TypeString,
				Required: true,
			},

			"auto_upgrade_minor_version": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"settings": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: structure.SuppressJsonDiff,
			},

			// due to the sensitive nature, these are not returned by the API
			"protected_settings": {
				Type:             schema.TypeString,
				Optional:         true,
				Sensitive:        true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: structure.SuppressJsonDiff,
			},

			// TODO: Remove location and resource_group_name in 2.0
			"location": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Optional:         true,
				Computed:         true,
				StateFunc:        azure.NormalizeLocation,
				DiffSuppressFunc: azure.SuppressLocationDiff,
				Deprecated:       "location is no longer used",
			},

			"resource_group_name": azure.SchemaResourceGroupNameDeprecated(),

			"tags": tags.Schema(),
		},
	}
}

func resourceArmVirtualMachineExtensionsCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	vmExtensionClient := meta.(*clients.Client).Compute.VMExtensionClient
	vmClient := meta.(*clients.Client).Compute.VMClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	var virtualMachineName, resourceGroup, location string
	if virtualMachineId, ok := d.GetOk("virtual_machine_id"); ok {
		v, err := ParseVirtualMachineID(virtualMachineId.(string))
		if err != nil {
			return fmt.Errorf("Error parsing Virtual Machine ID %q: %+v", virtualMachineId, err)
		}

		virtualMachineName = v.Name
		resourceGroup = v.ResourceGroup

		virtualMachine, err := vmClient.Get(ctx, resourceGroup, virtualMachineName, "")
		if err != nil {
			return fmt.Errorf("Error getting Virtual Machine %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		if location = *virtualMachine.Location; location == "" {
			return fmt.Errorf("Error reading location of Virtual Machine %q", virtualMachineName)
		}
	} else {
		if vm, ok := d.GetOk("virtual_machine_name"); !ok {
			return fmt.Errorf("Error, one of `virtual_machine_name` (deprecated) or `virtual_machine_id` (preferred) required.")
		} else {
			virtualMachineName = vm.(string)
			resourceGroup = d.Get("resource_group_name").(string)
			if resourceGroup == "" {
				return fmt.Errorf("`resource_group_name` must be specified if `virtual_machine_id` is not used")
			}

			location = azure.NormalizeLocation(d.Get("location").(string))
			if location == "" {
				return fmt.Errorf("`location` must be specified if `virtual_machine_id` is not used")
			}
		}
	}

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := vmExtensionClient.Get(ctx, resourceGroup, virtualMachineName, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Extension %q (Virtual Machine %q / Resource Group %q): %s", name, virtualMachineName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_virtual_machine_extension", *existing.ID)
		}
	}

	publisher := d.Get("publisher").(string)
	extensionType := d.Get("type").(string)
	typeHandlerVersion := d.Get("type_handler_version").(string)
	autoUpgradeMinor := d.Get("auto_upgrade_minor_version").(bool)
	t := d.Get("tags").(map[string]interface{})

	extension := compute.VirtualMachineExtension{
		Location: &location,
		VirtualMachineExtensionProperties: &compute.VirtualMachineExtensionProperties{
			Publisher:               &publisher,
			Type:                    &extensionType,
			TypeHandlerVersion:      &typeHandlerVersion,
			AutoUpgradeMinorVersion: &autoUpgradeMinor,
		},
		Tags: tags.Expand(t),
	}

	if settingsString := d.Get("settings").(string); settingsString != "" {
		settings, err := structure.ExpandJsonFromString(settingsString)
		if err != nil {
			return fmt.Errorf("unable to parse settings: %s", err)
		}
		extension.VirtualMachineExtensionProperties.Settings = &settings
	}

	if protectedSettingsString := d.Get("protected_settings").(string); protectedSettingsString != "" {
		protectedSettings, err := structure.ExpandJsonFromString(protectedSettingsString)
		if err != nil {
			return fmt.Errorf("unable to parse protected_settings: %s", err)
		}
		extension.VirtualMachineExtensionProperties.ProtectedSettings = &protectedSettings
	}

	future, err := vmExtensionClient.CreateOrUpdate(ctx, resourceGroup, virtualMachineName, name, extension)
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, vmExtensionClient.Client); err != nil {
		return err
	}

	read, err := vmExtensionClient.Get(ctx, resourceGroup, virtualMachineName, name, "")
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read  Virtual Machine Extension %s (resource group %s) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmVirtualMachineExtensionsRead(d, meta)
}

func resourceArmVirtualMachineExtensionsRead(d *schema.ResourceData, meta interface{}) error {
	vmExtensionClient := meta.(*clients.Client).Compute.VMExtensionClient
	vmClient := meta.(*clients.Client).Compute.VMClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := ParseVirtualMachineExtensionID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	vmName := id.VirtualMachine
	name := id.Name

	virtualMachine, err := vmClient.Get(ctx, resourceGroup, vmName, "")
	if err != nil {
		if utils.ResponseWasNotFound(virtualMachine.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Virtual Machine %s: %s", name, err)
	}

	d.Set("virtual_machine_id", virtualMachine.ID)
	if location := virtualMachine.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	resp, err := vmExtensionClient.Get(ctx, resourceGroup, vmName, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Virtual Machine Extension %s: %s", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("virtual_machine_name", vmName)
	d.Set("resource_group_name", resourceGroup)

	if props := resp.VirtualMachineExtensionProperties; props != nil {
		d.Set("publisher", props.Publisher)
		d.Set("type", props.Type)
		d.Set("type_handler_version", props.TypeHandlerVersion)
		d.Set("auto_upgrade_minor_version", props.AutoUpgradeMinorVersion)

		if settings := props.Settings; settings != nil {
			settingsVal := settings.(map[string]interface{})
			settingsJson, err := structure.FlattenJsonToString(settingsVal)
			if err != nil {
				return fmt.Errorf("unable to parse settings from response: %s", err)
			}
			d.Set("settings", settingsJson)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmVirtualMachineExtensionsDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VMExtensionClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["extensions"]
	vmName := id.Path["virtualMachines"]

	future, err := client.Delete(ctx, resGroup, vmName, name)
	if err != nil {
		return err
	}

	return future.WaitForCompletionRef(ctx, client.Client)
}
