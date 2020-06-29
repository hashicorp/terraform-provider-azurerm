package compute

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-12-01/compute"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// NOTE (also in the docs): this is not intended to be used with the `azurerm_virtual_machine_scale_set` resource

func resourceArmVirtualMachineScaleSetExtension() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmVirtualMachineScaleSetExtensionCreate,
		Read:   resourceArmVirtualMachineScaleSetExtensionRead,
		Update: resourceArmVirtualMachineScaleSetExtensionUpdate,
		Delete: resourceArmVirtualMachineScaleSetExtensionDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.VirtualMachineScaleSetExtensionID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"virtual_machine_scale_set_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.VirtualMachineScaleSetID,
			},

			"publisher": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"type_handler_version": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"auto_upgrade_minor_version": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"force_update_tag": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"protected_settings": {
				Type:             schema.TypeString,
				Optional:         true,
				Sensitive:        true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: structure.SuppressJsonDiff,
			},

			"provision_after_extensions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"settings": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: structure.SuppressJsonDiff,
			},
		},
	}
}

func resourceArmVirtualMachineScaleSetExtensionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VMScaleSetExtensionsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	virtualMachineScaleSetId, err := parse.VirtualMachineScaleSetID(d.Get("virtual_machine_scale_set_id").(string))
	if err != nil {
		return err
	}
	resourceGroup := virtualMachineScaleSetId.ResourceGroup
	vmssName := virtualMachineScaleSetId.Name

	if features.ShouldResourcesBeImported() {
		resp, err := client.Get(ctx, resourceGroup, vmssName, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error checking for existing Extension %q (Virtual Machine Scale Set %q / Resource Group %q): %+v", name, vmssName, resourceGroup, err)
			}
		}

		if !utils.ResponseWasNotFound(resp.Response) {
			return tf.ImportAsExistsError("azurerm_linux_virtual_machine_scale_set", *resp.ID)
		}
	}

	settings := map[string]interface{}{}
	if settingsString := d.Get("settings").(string); settingsString != "" {
		s, err := structure.ExpandJsonFromString(settingsString)
		if err != nil {
			return fmt.Errorf("unable to parse `settings`: %s", err)
		}
		settings = s
	}

	provisionAfterExtensionsRaw := d.Get("provision_after_extensions").([]interface{})
	provisionAfterExtensions := utils.ExpandStringSlice(provisionAfterExtensionsRaw)

	protectedSettings := map[string]interface{}{}
	if protectedSettingsString := d.Get("protected_settings").(string); protectedSettingsString != "" {
		ps, err := structure.ExpandJsonFromString(protectedSettingsString)
		if err != nil {
			return fmt.Errorf("unable to parse `protected_settings`: %s", err)
		}
		protectedSettings = ps
	}

	props := compute.VirtualMachineScaleSetExtension{
		Name: utils.String(name),
		VirtualMachineScaleSetExtensionProperties: &compute.VirtualMachineScaleSetExtensionProperties{
			Publisher:                utils.String(d.Get("publisher").(string)),
			Type:                     utils.String(d.Get("type").(string)),
			TypeHandlerVersion:       utils.String(d.Get("type_handler_version").(string)),
			AutoUpgradeMinorVersion:  utils.Bool(d.Get("auto_upgrade_minor_version").(bool)),
			ProtectedSettings:        protectedSettings,
			ProvisionAfterExtensions: provisionAfterExtensions,
			Settings:                 settings,
		},
	}
	if v, ok := d.GetOk("force_update_tag"); ok {
		props.VirtualMachineScaleSetExtensionProperties.ForceUpdateTag = utils.String(v.(string))
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, vmssName, name, props)
	if err != nil {
		return fmt.Errorf("Error creating Extension %q (Virtual Machine Scale Set %q / Resource Group %q): %+v", name, vmssName, resourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Extension %q (Virtual Machine Scale Set %q / Resource Group %q): %+v", name, vmssName, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, vmssName, name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Extension %q (Virtual Machine Scale Set %q / Resource Group %q): %+v", name, vmssName, resourceGroup, err)
	}
	d.SetId(*resp.ID)

	return resourceArmVirtualMachineScaleSetExtensionRead(d, meta)
}

func resourceArmVirtualMachineScaleSetExtensionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VMScaleSetExtensionsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualMachineScaleSetExtensionID(d.Id())
	if err != nil {
		return err
	}

	props := compute.VirtualMachineScaleSetExtensionProperties{
		// if this isn't specified it defaults to false
		AutoUpgradeMinorVersion: utils.Bool(d.Get("auto_upgrade_minor_version").(bool)),
	}

	if d.HasChange("force_update_tag") {
		props.ForceUpdateTag = utils.String(d.Get("force_update_tag").(string))
	}

	if d.HasChange("protected_settings") {
		protectedSettings := map[string]interface{}{}
		if protectedSettingsString := d.Get("protected_settings").(string); protectedSettingsString != "" {
			ps, err := structure.ExpandJsonFromString(protectedSettingsString)
			if err != nil {
				return fmt.Errorf("unable to parse `protected_settings`: %s", err)
			}
			protectedSettings = ps
		}

		props.ProtectedSettings = protectedSettings
	}

	if d.HasChange("provision_after_extensions") {
		provisionAfterExtensionsRaw := d.Get("provision_after_extensions").([]interface{})
		props.ProvisionAfterExtensions = utils.ExpandStringSlice(provisionAfterExtensionsRaw)
	}

	if d.HasChange("publisher") {
		props.Publisher = utils.String(d.Get("publisher").(string))
	}

	if d.HasChange("settings") {
		settings := map[string]interface{}{}

		if settingsString := d.Get("settings").(string); settingsString != "" {
			s, err := structure.ExpandJsonFromString(settingsString)
			if err != nil {
				return fmt.Errorf("unable to parse `settings`: %s", err)
			}
			settings = s
		}

		props.Settings = settings
	}

	if d.HasChange("type") {
		props.Type = utils.String(d.Get("type").(string))
	}

	if d.HasChange("type_handler_version") {
		props.TypeHandlerVersion = utils.String(d.Get("type_handler_version").(string))
	}

	extension := compute.VirtualMachineScaleSetExtension{
		Name: utils.String(id.Name),
		VirtualMachineScaleSetExtensionProperties: &props,
	}
	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.VirtualMachineScaleSetName, id.Name, extension)
	if err != nil {
		return fmt.Errorf("Error updating Extension %q (Virtual Machine Scale Set %q / Resource Group %q): %+v", id.Name, id.VirtualMachineScaleSetName, id.ResourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for update of Extension %q (Virtual Machine Scale Set %q / Resource Group %q): %+v", id.Name, id.VirtualMachineScaleSetName, id.ResourceGroup, err)
	}

	return resourceArmVirtualMachineScaleSetExtensionRead(d, meta)
}

func resourceArmVirtualMachineScaleSetExtensionRead(d *schema.ResourceData, meta interface{}) error {
	vmssClient := meta.(*clients.Client).Compute.VMScaleSetClient
	client := meta.(*clients.Client).Compute.VMScaleSetExtensionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualMachineScaleSetExtensionID(d.Id())
	if err != nil {
		return err
	}

	vmss, err := vmssClient.Get(ctx, id.ResourceGroup, id.VirtualMachineScaleSetName)
	if err != nil {
		if utils.ResponseWasNotFound(vmss.Response) {
			log.Printf("Virtual Machine Scale Set %q was not found in Resource Group %q - removing Extension from state!", id.VirtualMachineScaleSetName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Virtual Machine Scale Set %q (Resource Group %q): %+v", id.VirtualMachineScaleSetName, id.ResourceGroup, err)
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.VirtualMachineScaleSetName, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("Extension %q (Virtual Machine Scale Set %q / Resource Group %q) was not found - removing from state!", id.Name, id.VirtualMachineScaleSetName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Extension %q (Virtual Machine Scale Set %q / Resource Group %q): %+v", id.Name, id.VirtualMachineScaleSetName, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("virtual_machine_scale_set_id", vmss.ID)

	if props := resp.VirtualMachineScaleSetExtensionProperties; props != nil {
		d.Set("auto_upgrade_minor_version", props.AutoUpgradeMinorVersion)
		d.Set("force_update_tag", props.ForceUpdateTag)
		d.Set("provision_after_extensions", utils.FlattenStringSlice(props.ProvisionAfterExtensions))
		d.Set("publisher", props.Publisher)
		d.Set("type", props.Type)
		d.Set("type_handler_version", props.TypeHandlerVersion)

		settings := ""
		if props.Settings != nil {
			settingsVal, ok := props.Settings.(map[string]interface{})
			if ok {
				settingsJson, err := structure.FlattenJsonToString(settingsVal)
				if err != nil {
					return fmt.Errorf("unable to parse settings from response: %s", err)
				}
				settings = settingsJson
			}
		}
		d.Set("settings", settings)
	}

	return nil
}

func resourceArmVirtualMachineScaleSetExtensionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VMScaleSetExtensionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualMachineScaleSetExtensionID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.VirtualMachineScaleSetName, id.Name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return fmt.Errorf("Error deleting Extension %q (Virtual Machine Scale Set %q / Resource Group %q): %+v", id.Name, id.VirtualMachineScaleSetName, id.ResourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of Extension %q (Virtual Machine Scale Set %q / Resource Group %q): %+v", id.Name, id.VirtualMachineScaleSetName, id.ResourceGroup, err)
	}

	return nil
}
