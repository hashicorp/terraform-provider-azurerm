package compute

import (
	"fmt"
	"log"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceDedicatedHostGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceDedicatedHostGroupCreate,
		Read:   resourceDedicatedHostGroupRead,
		Update: resourceDedicatedHostGroupUpdate,
		Delete: resourceDedicatedHostGroupDelete,

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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateDedicatedHostGroupName(),
			},

			"location": azure.SchemaLocation(),

			// There's a bug in the Azure API where this is returned in upper-case
			// BUG: https://github.com/Azure/azure-rest-api-specs/issues/8068
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"platform_fault_domain_count": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 3),
			},

			// Currently only one endpoint is allowed.
			// we'll leave this open to enhancement when they add multiple zones support.
			"zones": azure.SchemaSingleZone(),

			"tags": tags.Schema(),
		},
	}
}

func resourceDedicatedHostGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DedicatedHostGroupsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroupName := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroupName, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for present of existing Dedicated Host Group %q (Resource Group %q): %+v", name, resourceGroupName, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_dedicated_host_group", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	platformFaultDomainCount := d.Get("platform_fault_domain_count").(int)
	t := d.Get("tags").(map[string]interface{})

	parameters := compute.DedicatedHostGroup{
		Location: utils.String(location),
		DedicatedHostGroupProperties: &compute.DedicatedHostGroupProperties{
			PlatformFaultDomainCount: utils.Int32(int32(platformFaultDomainCount)),
		},
		Tags: tags.Expand(t),
	}
	if zones, ok := d.GetOk("zones"); ok {
		parameters.Zones = utils.ExpandStringSlice(zones.([]interface{}))
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroupName, name, parameters); err != nil {
		return fmt.Errorf("Error creating Dedicated Host Group %q (Resource Group %q): %+v", name, resourceGroupName, err)
	}

	resp, err := client.Get(ctx, resourceGroupName, name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Dedicated Host Group %q (Resource Group %q): %+v", name, resourceGroupName, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Dedicated Host Group %q (Resource Group %q) ID", name, resourceGroupName)
	}
	d.SetId(*resp.ID)

	return resourceDedicatedHostGroupRead(d, meta)
}

func resourceDedicatedHostGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DedicatedHostGroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroupName := id.ResourceGroup
	name := id.Path["hostGroups"]

	resp, err := client.Get(ctx, resourceGroupName, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Dedicated Host Group %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Dedicated Host Group %q (Resource Group %q): %+v", name, resourceGroupName, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroupName)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if props := resp.DedicatedHostGroupProperties; props != nil {
		platformFaultDomainCount := 0
		if props.PlatformFaultDomainCount != nil {
			platformFaultDomainCount = int(*props.PlatformFaultDomainCount)
		}
		d.Set("platform_fault_domain_count", platformFaultDomainCount)
	}
	d.Set("zones", utils.FlattenStringSlice(resp.Zones))

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceDedicatedHostGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DedicatedHostGroupsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroupName := d.Get("resource_group_name").(string)
	t := d.Get("tags").(map[string]interface{})

	parameters := compute.DedicatedHostGroupUpdate{
		Tags: tags.Expand(t),
	}

	if _, err := client.Update(ctx, resourceGroupName, name, parameters); err != nil {
		return fmt.Errorf("Error updating Dedicated Host Group %q (Resource Group %q): %+v", name, resourceGroupName, err)
	}

	return resourceDedicatedHostGroupRead(d, meta)
}

func resourceDedicatedHostGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DedicatedHostGroupsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["hostGroups"]

	if _, err := client.Delete(ctx, resourceGroup, name); err != nil {
		return fmt.Errorf("Error deleting Dedicated Host Group %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return nil
}
