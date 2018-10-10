package azurerm

import (
	"bytes"
	"fmt"
	"log"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/preview/devspaces/mgmt/2018-06-01-preview/devspaces"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
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

		Schema: map[string] *schema.Schema {
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"tags": tagsSchema(),

			"sku": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource {
					Schema: map[string] *schema.Schema {
						"name": {
							Type: schema.TypeString,
							Required: true,
						},
						"tier": {
							Type: schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"host_suffix": {
				Type: schema.TypeString,
			},

			"data_plane_fqdn": {
				Type: schema.TypeString,
				Computed: true,
			},

			"target_container_host_resource_id": {
				Type: schema.TypeString,
				Required: true,
			},

			"target_container_host_credentials_base64": {
				Type: schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceArmDevSpacesControllerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).devSpacesControllersClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for DevSpaces Controller creation")

	name := d.Get("name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	resourceGroupName := d.Get("resource_group_name").(string)
	tags := d.Get("tags").(map[string]interface{})

	sku := expandDevSpacesControllerSku(d)

	hostSuffix := d.Get("host_suffix").(string)
	tarCHResId := d.Get("target_container_host_resource_id").(string)
	tarCHCredBase64 := d.Get("target_container_host_credentials_base64").(string)

	controller := devspaces.Controller {
		Location: &location,
		Tags: expandTags(tags),
		Sku: &sku,
		ControllerProperties: &devspaces.ControllerProperties {
			HostSuffix: &hostSuffix,
			TargetContainerHostResourceID: &tarCHResId,
			TargetContainerHostCredentialsBase64: &tarCHCredBase64,
		},
	}

	future, err := client.Create(ctx, resourceGroupName, name, controller)
	if err != nil {
		return fmt.Errorf("Error creating DevSpaces Controller %q (Resource Group %q): %+v", name, resourceGroupName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of DevSpaces Controller %q (Resource Group %q): %+v", name, resourceGroupName, err)
	}

	var result devspaces.Controller
	if result, err = future.Result(client); err != nil {
		return fmt.Errorf("Error retrieving result of DevSpaces Controller %q (Resource Group %q): %+v", name, resourceGroupName, err)
	}
	if result.ID == nil {
		return fmt.Errorf("Cannot read DevSpaces Controller %q (Resource Group %q) ID", name, resourceGroupName)
	}
	d.SetId(*result.ID)

	return resourceArmDevSpacesControllerRead(d, meta)
}

func resourceArmDevSpacesControllerRead(d *schema.ResourceData, meta interface{}) error {

}

func resourceArmDevSpacesControllerUpdate(d *schema.ResourceData, meta interface{}) error {
}

func resourceArmDevSpacesControllerDelete(d *schema.ResourceData, meta interface{}) error {
}

func expandDevSpacesControllerSku(d *schema.ResourceData) devspaces.Sku {
	skus := d.Get("sku").([]interface)
	skuData := skus[0].(map[string]interface{})

	name := skuData["name"].(string)
	tier := skuData["tier"].(devspaces.SkuTier)

	return devspaces.Sku {
		Name: &name,
		Tier: tier,
	}
}

