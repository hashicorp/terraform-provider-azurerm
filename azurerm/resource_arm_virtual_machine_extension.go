package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2017-12-01/compute"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/structure"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmVirtualMachineExtensions() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmVirtualMachineExtensionsCreate,
		Read:   resourceArmVirtualMachineExtensionsRead,
		Update: resourceArmVirtualMachineExtensionsCreate,
		Delete: resourceArmVirtualMachineExtensionsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"virtual_machine_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
				ValidateFunc:     validation.ValidateJsonString,
				DiffSuppressFunc: structure.SuppressJsonDiff,
			},

			// due to the sensitive nature, these are not returned by the API
			"protected_settings": {
				Type:             schema.TypeString,
				Optional:         true,
				Sensitive:        true,
				ValidateFunc:     validation.ValidateJsonString,
				DiffSuppressFunc: structure.SuppressJsonDiff,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmVirtualMachineExtensionsCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).vmExtensionClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	vmName := d.Get("virtual_machine_name").(string)
	resGroup := d.Get("resource_group_name").(string)
	publisher := d.Get("publisher").(string)
	extensionType := d.Get("type").(string)
	typeHandlerVersion := d.Get("type_handler_version").(string)
	autoUpgradeMinor := d.Get("auto_upgrade_minor_version").(bool)
	tags := d.Get("tags").(map[string]interface{})

	extension := compute.VirtualMachineExtension{
		Location: &location,
		VirtualMachineExtensionProperties: &compute.VirtualMachineExtensionProperties{
			Publisher:               &publisher,
			Type:                    &extensionType,
			TypeHandlerVersion:      &typeHandlerVersion,
			AutoUpgradeMinorVersion: &autoUpgradeMinor,
		},
		Tags: expandTags(tags),
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

	future, err := client.CreateOrUpdate(ctx, resGroup, vmName, name, extension)
	if err != nil {
		return err
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, resGroup, vmName, name, "")
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read  Virtual Machine Extension %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmVirtualMachineExtensionsRead(d, meta)
}

func resourceArmVirtualMachineExtensionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).vmExtensionClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	vmName := id.Path["virtualMachines"]
	name := id.Path["extensions"]

	resp, err := client.Get(ctx, resGroup, vmName, name, "")

	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Virtual Machine Extension %s: %s", name, err)
	}

	d.Set("name", resp.Name)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}
	d.Set("virtual_machine_name", vmName)
	d.Set("resource_group_name", resGroup)

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

	if resp.VirtualMachineExtensionProperties.Settings != nil {
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmVirtualMachineExtensionsDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).vmExtensionClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
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

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return err
	}

	return nil
}
