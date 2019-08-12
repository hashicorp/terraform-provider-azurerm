package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-06-01/compute"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmProximityPlacementGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmProximityPlacementGroupCreateUpdate,
		Read:   resourceArmProximityPlacementGroupRead,
		Update: resourceArmProximityPlacementGroupCreateUpdate,
		Delete: resourceArmProximityPlacementGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"tags": tagsSchema(),
		},
	}
}

func resourceArmProximityPlacementGroupCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).ppgClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM Proximity Placement Group creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Proximity Placement Group %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_proximity_placement_group", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	tags := d.Get("tags").(map[string]interface{})

	ppg := compute.ProximityPlacementGroup{
		Name:     &name,
		Location: &location,
		Tags:     expandTags(tags),
	}

	resp, err := client.CreateOrUpdate(ctx, resGroup, name, ppg)
	if err != nil {
		return err
	}

	d.SetId(*resp.ID)

	return resourceArmProximityPlacementGroupRead(d, meta)
}

func resourceArmProximityPlacementGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).ppgClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["proximityPlacementGroups"]

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Proximity Placement Group %q (Resource Group %q): %+v", name, resGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmProximityPlacementGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).ppgClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["proximityPlacementGroups"]

	_, err = client.Delete(ctx, resGroup, name)

	return err
}
