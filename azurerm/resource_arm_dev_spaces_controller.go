package azurerm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/validation"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/devspaces/mgmt/2018-06-01-preview/devspaces"
	"github.com/hashicorp/terraform/helper/schema"
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"sku": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string {
								"S1",
							}, false),
						},
						"tier": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string {
								string(devspaces.Standard),
							}, false),
						},
					},
				},
			},

			"host_suffix": {
				Type: schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"data_plane_fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"target_container_host_resource_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"target_container_host_credentials_base64": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Sensitive: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmDevSpaceControllerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).devSpaceControllerClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for DevSpace Controller creation")

	name := d.Get("name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	resGroupName := d.Get("resource_group_name").(string)
	tags := d.Get("tags").(map[string]interface{})

	sku := expandDevSpaceControllerSku(d)

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
		return fmt.Errorf("Error creating DevSpace Controller %q (Resource Group %q): %+v", name, resGroupName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of DevSpace Controller %q (Resource Group %q): %+v", name, resGroupName, err)
	}

	var result devspaces.Controller
	if result, err = future.Result(client); err != nil {
		return fmt.Errorf("Error retrieving result of DevSpace Controller %q (Resource Group %q): %+v", name, resGroupName, err)
	}
	if result.ID == nil {
		return fmt.Errorf("Cannot read DevSpace Controller %q (Resource Group %q) ID", name, resGroupName)
	}
	d.SetId(*result.ID)

	return resourceArmDevSpaceControllerRead(d, meta)
}

func resourceArmDevSpaceControllerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).devSpaceControllerClient
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
			log.Printf("[DEBUG] DevSpace Controller %q was not found in Resource Group %q - removing from state!", name, resGroupName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on DevSpace Controller Lab %q (Resource Group %q): %+v", name, resGroupName, err)
	}

	d.Set("name", result.Name)
	d.Set("resource_group_name", resGroupName)
	if location := result.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	d.Set("sku", flattenDevSpaceControllerSku(result.Sku))

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

func resourceArmDevSpaceControllerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).devSpaceControllerClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for DevSpace Controller updating")

	name := d.Get("name").(string)
	resGroupName := d.Get("resource_group_name").(string)
	tags := d.Get("tags").(map[string]interface{})

	params := devspaces.ControllerUpdateParameters{
		Tags: expandTags(tags),
	}

	_, err := client.Update(ctx, resGroupName, name, params)
	if err != nil {
		return fmt.Errorf("Error updating DevSpace Controller %q (Resource Group %q): %+v", name, resGroupName, err)
	}

	return nil
}

func resourceArmDevSpaceControllerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).devSpaceControllerClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroupName := id.ResourceGroup
	name := id.Path["controllers"]

	future, err := client.Delete(ctx, resGroupName, name)
	if err != nil {
		return fmt.Errorf("Error deleting DevSpace Controller %q (Resource Group %q): %+v", name, resGroupName, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for the deletion of DevSpace Controller %q (Resource Group %q): %+v", name, resGroupName, err)
	}

	return nil
}

func expandDevSpaceControllerSku(d *schema.ResourceData) (sku devspaces.Sku) {
	skuConfigs := d.Get("sku").([]interface{})
	skuConfig := skuConfigs[0].(map[string]interface{})

	skuName := skuConfig["name"].(string)
	sku.Name = &skuName
	sku.Tier = (devspaces.SkuTier)(skuConfig["tier"].(string))

	return
}

func flattenDevSpaceControllerSku(skuObj *devspaces.Sku) (skuConfigs []interface{}) {
	if skuObj == nil {
		return
	}

	skuConfig := make(map[string]interface{}, 0)
	skuConfig["name"] = *skuObj.Name
	skuConfig["tier"] = skuObj.Tier

	skuConfigs = append(skuConfigs, skuConfig)

	return
}
