package powerbidedicated

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/powerbidedicated/mgmt/2017-10-01/powerbidedicated"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmPowerBIDedicatedCapacity() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPowerBIDedicatedCapacityCreate,
		Read:   resourceArmPowerBIDedicatedCapacityRead,
		Update: resourceArmPowerBIDedicatedCapacityUpdate,
		Delete: resourceArmPowerBIDedicatedCapacityDelete,

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
				ValidateFunc: ValidateCapacityName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"sku_name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"A1",
					"A2",
					"A3",
					"A4",
					"A5",
					"A6",
				}, false),
			},

			"administrators": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: ValidateCapacityAdministratorName,
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmPowerBIDedicatedCapacityCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PowerBIDedicated.CapacityClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() {
		existing, err := client.GetDetails(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for present of existing Capacity %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_powerbi_dedicated_capacity", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	administrators := d.Get("administrators").(*schema.Set).List()
	skuName := d.Get("sku_name").(string)
	t := d.Get("tags").(map[string]interface{})

	parameters := powerbidedicated.DedicatedCapacity{
		Location: utils.String(location),
		DedicatedCapacityProperties: &powerbidedicated.DedicatedCapacityProperties{
			Administration: &powerbidedicated.DedicatedCapacityAdministrators{
				Members: utils.ExpandStringSlice(administrators),
			},
		},
		Sku: &powerbidedicated.ResourceSku{
			Name: utils.String(skuName),
		},
		Tags: tags.Expand(t),
	}

	future, err := client.Create(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating Capacity %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Capacity %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.GetDetails(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Capacity %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("Cannot read Capacity %q (Resource Group %q) ID", name, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmPowerBIDedicatedCapacityRead(d, meta)
}

func resourceArmPowerBIDedicatedCapacityRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PowerBIDedicated.CapacityClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["capacities"]

	resp, err := client.GetDetails(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Capacity %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Capacity %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if props := resp.DedicatedCapacityProperties; props != nil {
		if err := d.Set("administrators", utils.FlattenStringSlice(props.Administration.Members)); err != nil {
			return fmt.Errorf("Error setting `administration`: %+v", err)
		}
	}

	skuName := ""
	if resp.Sku != nil && resp.Sku.Name != nil {
		skuName = *resp.Sku.Name
	}
	d.Set("sku_name", skuName)

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmPowerBIDedicatedCapacityUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PowerBIDedicated.CapacityClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	administrators := d.Get("administrators").(*schema.Set).List()
	skuName := d.Get("sku_name").(string)
	t := d.Get("tags").(map[string]interface{})

	parameters := powerbidedicated.DedicatedCapacityUpdateParameters{
		DedicatedCapacityMutableProperties: &powerbidedicated.DedicatedCapacityMutableProperties{
			Administration: &powerbidedicated.DedicatedCapacityAdministrators{
				Members: utils.ExpandStringSlice(administrators),
			},
		},
		Sku: &powerbidedicated.ResourceSku{
			Name: utils.String(skuName),
		},
		Tags: tags.Expand(t),
	}

	future, err := client.Update(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error updating Capacity %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for update of Capacity %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return resourceArmPowerBIDedicatedCapacityRead(d, meta)
}

func resourceArmPowerBIDedicatedCapacityDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PowerBIDedicated.CapacityClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["capacities"]

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting Capacity %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deleting Capacity %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}
