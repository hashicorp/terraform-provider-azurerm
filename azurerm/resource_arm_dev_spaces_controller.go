package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/devspaces/mgmt/2018-06-01-preview/devspaces"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDevSpacesController() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDevSpacesControllerCreate,
		Read:   resourceArmDevSpacesControllerRead,
		Update: resourceArmDevSpacesControllerUpdate,
		Delete: resourceArmDevSpacesControllerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"sku": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"tier": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"host_suffix": {
				Type: schema.TypeString,
			},

			"data_plane_fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"target_container_host_resource_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"target_container_host_credentials_base64": {
				Type:     schema.TypeString,
				Required: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmDevSpacesControllerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).devSpacesControllersClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for DevSpaces Controller creation")

	name := d.Get("name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	resGroupName := d.Get("resource_group_name").(string)
	tags := d.Get("tags").(map[string]interface{})

	sku := expandDevSpacesControllerSku(d)

	hostSuffix := d.Get("host_suffix").(string)
	tarCHResId := d.Get("target_container_host_resource_id").(string)
	tarCHCredBase64 := d.Get("target_container_host_credentials_base64").(string)

	controller := devspaces.Controller{
		Location: &location,
		Tags:     expandTags(tags),
		Sku:      &sku,
		ControllerProperties: &devspaces.ControllerProperties{
			HostSuffix:                           &hostSuffix,
			TargetContainerHostResourceID:        &tarCHResId,
			TargetContainerHostCredentialsBase64: &tarCHCredBase64,
		},
	}

	future, err := client.Create(ctx, resGroupName, name, controller)
	if err != nil {
		return fmt.Errorf("Error creating DevSpaces Controller %q (Resource Group %q): %+v", name, resGroupName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of DevSpaces Controller %q (Resource Group %q): %+v", name, resGroupName, err)
	}

	var result devspaces.Controller
	if result, err = future.Result(client); err != nil {
		return fmt.Errorf("Error retrieving result of DevSpaces Controller %q (Resource Group %q): %+v", name, resGroupName, err)
	}
	if result.ID == nil {
		return fmt.Errorf("Cannot read DevSpaces Controller %q (Resource Group %q) ID", name, resGroupName)
	}
	d.SetId(*result.ID)

	return resourceArmDevSpacesControllerRead(d, meta)
}

func resourceArmDevSpacesControllerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).devSpacesControllersClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroupName := id.ResourceGroup
	name := id.Path["controllers"]

	result, err := client.Get(ctx, resGroupName, name)
	if err != nil {
		if utils.ResponseWasNotFound(result.Response) {
			log.Printf("[DEBUG] DevSpaces Controller %q was not found in Resource Group %q - removing from state!", name, resGroupName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on DevSpaces Controller Lab %q (Resource Group %q): %+v", name, resGroupName, err)
	}

	d.Set("name", result.Name)
	d.Set("resource_group_name", resGroupName)
	if location := result.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	d.Set("sku", flattenDevSpacesControllerSku(result.Sku))

	if props := result.ControllerProperties; props != nil {
		if props.HostSuffix != nil {
			d.Set("host_suffix", props.HostSuffix)
		}

		if props.DataPlaneFqdn != nil {
			d.Set("data_plane_fqdn", props.DataPlaneFqdn)
		}

		if props.TargetContainerHostResourceID != nil {
			d.Set("target_container_host_resource_id", props.TargetContainerHostResourceID)
		}

		if props.TargetContainerHostCredentialsBase64 != nil {
			d.Set("target_container_host_credentials_base64", props.TargetContainerHostCredentialsBase64)
		}
	}

	flattenAndSetTags(d, result.Tags)

	return nil
}

func resourceArmDevSpacesControllerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).devSpacesControllersClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for DevSpaces Controller updating")

	name := d.Get("name").(string)
	resGroupName := d.Get("resource_group_name").(string)
	tags := d.Get("tags").(map[string]interface{})

	params := devspaces.ControllerUpdateParameters{
		Tags: expandTags(tags),
	}

	_, err := client.Update(ctx, resGroupName, name, params)
	if err != nil {
		return fmt.Errorf("Error updating DevSpaces Controller %q (Resource Group %q): %+v", name, resGroupName, err)
	}

	return nil
}

func resourceArmDevSpacesControllerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).devSpacesControllersClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroupName := id.ResourceGroup
	name := id.Path["controllers"]

	future, err := client.Delete(ctx, resGroupName, name)
	if err != nil {
		return fmt.Errorf("Error deleting DevSpaces Controller %q (Resource Group %q): %+v", name, resGroupName, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for the deletion of DevSpaces Controller %q (Resource Group %q): %+v", name, resGroupName, err)
	}

	return nil
}

func expandDevSpacesControllerSku(d *schema.ResourceData) (sku devspaces.Sku) {
	skuConfigs := d.Get("sku").([]interface{})
	skuConfig := skuConfigs[0].(map[string]interface{})

	skuName := skuConfig["name"].(string)
	sku.Name = &skuName
	sku.Tier = skuConfig["tier"].(devspaces.SkuTier)

	return
}

func flattenDevSpacesControllerSku(skuObj *devspaces.Sku) (skuConfigs []interface{}) {
	if skuObj == nil {
		return
	}

	skuConfig := make(map[string]interface{}, 0)
	skuConfig["name"] = *skuObj.Name
	skuConfig["tier"] = skuObj.Tier

	skuConfigs = append(skuConfigs, skuConfig)

	return
}
