package azurerm

import (
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-06-01/compute"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAvailabilitySet() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAvailabilitySetCreateUpdate,
		Read:   resourceArmAvailabilitySetRead,
		Update: resourceArmAvailabilitySetCreateUpdate,
		Delete: resourceArmAvailabilitySetDelete,
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

			"platform_update_domain_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      5,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 20),
			},

			"platform_fault_domain_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 3),
			},

			"managed": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmAvailabilitySetCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).compute.AvailabilitySetsClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM Availability Set creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Availability Set %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_availability_set", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	updateDomainCount := d.Get("platform_update_domain_count").(int)
	faultDomainCount := d.Get("platform_fault_domain_count").(int)
	managed := d.Get("managed").(bool)
	tags := d.Get("tags").(map[string]interface{})

	availSet := compute.AvailabilitySet{
		Name:     &name,
		Location: &location,
		AvailabilitySetProperties: &compute.AvailabilitySetProperties{
			PlatformFaultDomainCount:  utils.Int32(int32(faultDomainCount)),
			PlatformUpdateDomainCount: utils.Int32(int32(updateDomainCount)),
		},
		Tags: expandTags(tags),
	}

	if managed {
		n := "Aligned"
		availSet.Sku = &compute.Sku{
			Name: &n,
		}
	}

	resp, err := client.CreateOrUpdate(ctx, resGroup, name, availSet)
	if err != nil {
		return err
	}

	d.SetId(*resp.ID)

	return resourceArmAvailabilitySetRead(d, meta)
}

func resourceArmAvailabilitySetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).compute.AvailabilitySetsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["availabilitySets"]

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure Availability Set %q (Resource Group %q): %+v", name, resGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if resp.Sku != nil && resp.Sku.Name != nil {
		d.Set("managed", strings.EqualFold(*resp.Sku.Name, "Aligned"))
	}

	if props := resp.AvailabilitySetProperties; props != nil {
		d.Set("platform_update_domain_count", props.PlatformUpdateDomainCount)
		d.Set("platform_fault_domain_count", props.PlatformFaultDomainCount)
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmAvailabilitySetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).compute.AvailabilitySetsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["availabilitySets"]

	_, err = client.Delete(ctx, resGroup, name)

	return err
}
