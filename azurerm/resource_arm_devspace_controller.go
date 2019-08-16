package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/devspaces/mgmt/2018-06-01-preview/devspaces"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDevSpaceController() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDevSpaceControllerCreate,
		Read:   resourceArmDevSpaceControllerRead,
		Update: resourceArmDevSpaceControllerUpdate,
		Delete: resourceArmDevSpaceControllerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DevSpaceName(),
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"sku": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"S1",
							}, false),
						},
						"tier": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(devspaces.Standard),
							}, false),
						},
					},
				},
			},

			"host_suffix": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"target_container_host_resource_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"target_container_host_credentials_base64": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Sensitive:    true,
				ValidateFunc: validate.Base64String(),
			},

			"tags": tagsSchema(),

			"data_plane_fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmDevSpaceControllerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).devSpace.ControllersClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for DevSpace Controller creation")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing DevSpace Controller %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_devspace_controller", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	tags := d.Get("tags").(map[string]interface{})

	sku := expandDevSpaceControllerSku(d)

	hostSuffix := d.Get("host_suffix").(string)
	tarCHResId := d.Get("target_container_host_resource_id").(string)
	tarCHCredBase64 := d.Get("target_container_host_credentials_base64").(string)

	controller := devspaces.Controller{
		Location: &location,
		Tags:     expandTags(tags),
		Sku:      sku,
		ControllerProperties: &devspaces.ControllerProperties{
			HostSuffix:                           &hostSuffix,
			TargetContainerHostResourceID:        &tarCHResId,
			TargetContainerHostCredentialsBase64: &tarCHCredBase64,
		},
	}

	future, err := client.Create(ctx, resGroup, name, controller)
	if err != nil {
		return fmt.Errorf("Error creating DevSpace Controller %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of DevSpace Controller %q (Resource Group %q): %+v", name, resGroup, err)
	}

	result, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving DevSpace %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if result.ID == nil {
		return fmt.Errorf("Cannot read DevSpace Controller %q (Resource Group %q) ID", name, resGroup)
	}
	d.SetId(*result.ID)

	return resourceArmDevSpaceControllerRead(d, meta)
}

func resourceArmDevSpaceControllerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).devSpace.ControllersClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroupName := id.ResourceGroup
	name := id.Path["controllers"]

	result, err := client.Get(ctx, resGroupName, name)
	if err != nil {
		if utils.ResponseWasNotFound(result.Response) {
			log.Printf("[DEBUG] DevSpace Controller %q was not found in Resource Group %q - removing from state!", name, resGroupName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on DevSpace Controller %q (Resource Group %q): %+v", name, resGroupName, err)
	}

	d.Set("name", result.Name)
	d.Set("resource_group_name", resGroupName)
	if location := result.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if err := d.Set("sku", flattenDevSpaceControllerSku(result.Sku)); err != nil {
		return fmt.Errorf("Error flattenning `sku`: %+v", err)
	}

	if props := result.ControllerProperties; props != nil {
		d.Set("host_suffix", props.HostSuffix)
		d.Set("data_plane_fqdn", props.DataPlaneFqdn)
		d.Set("target_container_host_resource_id", props.TargetContainerHostResourceID)
	}

	flattenAndSetTags(d, result.Tags)

	return nil
}

func resourceArmDevSpaceControllerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).devSpace.ControllersClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for DevSpace Controller updating")

	name := d.Get("name").(string)
	resGroupName := d.Get("resource_group_name").(string)
	tags := d.Get("tags").(map[string]interface{})

	params := devspaces.ControllerUpdateParameters{
		Tags: expandTags(tags),
	}

	result, err := client.Update(ctx, resGroupName, name, params)
	if err != nil {
		return fmt.Errorf("Error updating DevSpace Controller %q (Resource Group %q): %+v", name, resGroupName, err)
	}

	if result.ID == nil {
		return fmt.Errorf("Cannot read DevSpace Controller %q (Resource Group %q) ID", name, resGroupName)
	}
	d.SetId(*result.ID)

	return resourceArmDevSpaceControllerRead(d, meta)
}

func resourceArmDevSpaceControllerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).devSpace.ControllersClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroupName := id.ResourceGroup
	name := id.Path["controllers"]

	future, err := client.Delete(ctx, resGroupName, name)
	if err != nil {
		return fmt.Errorf("Error deleting DevSpace Controller %q (Resource Group %q): %+v", name, resGroupName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the deletion of DevSpace Controller %q (Resource Group %q): %+v", name, resGroupName, err)
	}

	return nil
}

func expandDevSpaceControllerSku(d *schema.ResourceData) *devspaces.Sku {
	if _, ok := d.GetOk("sku"); !ok {
		return nil
	}

	skuConfigs := d.Get("sku").([]interface{})
	skuConfig := skuConfigs[0].(map[string]interface{})
	skuName := skuConfig["name"].(string)
	skuTier := devspaces.SkuTier(skuConfig["tier"].(string))

	return &devspaces.Sku{
		Name: &skuName,
		Tier: skuTier,
	}
}

func flattenDevSpaceControllerSku(skuObj *devspaces.Sku) []interface{} {
	if skuObj == nil {
		return []interface{}{}
	}

	skuConfig := make(map[string]interface{})
	if skuObj.Name != nil {
		skuConfig["name"] = *skuObj.Name
	}

	skuConfig["tier"] = skuObj.Tier

	return []interface{}{skuConfig}
}
