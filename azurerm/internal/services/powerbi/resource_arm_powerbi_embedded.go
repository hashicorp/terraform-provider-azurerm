package powerbi

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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/powerbi/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmPowerBIEmbedded() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPowerBIEmbeddedCreate,
		Read:   resourceArmPowerBIEmbeddedRead,
		Update: resourceArmPowerBIEmbeddedUpdate,
		Delete: resourceArmPowerBIEmbeddedDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.PowerBIEmbeddedID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidatePowerBIEmbeddedName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

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
					ValidateFunc: ValidatePowerBIEmbeddedAdministratorName,
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmPowerBIEmbeddedCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PowerBI.CapacityClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() {
		existing, err := client.GetDetails(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for present of existing PowerBI Embedded %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_powerbi_embedded", *existing.ID)
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
		return fmt.Errorf("Error creating PowerBI Embedded %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of PowerBI Embedded %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.GetDetails(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving PowerBI Embedded %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("Cannot read PowerBI Embedded %q (Resource Group %q) ID", name, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmPowerBIEmbeddedRead(d, meta)
}

func resourceArmPowerBIEmbeddedRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PowerBI.CapacityClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PowerBIEmbeddedID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetDetails(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] PowerBI Embedded %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading PowerBI Embedded %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
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

func resourceArmPowerBIEmbeddedUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PowerBI.CapacityClient
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
		return fmt.Errorf("Error updating PowerBI Embedded %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for update of PowerBI Embedded %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return resourceArmPowerBIEmbeddedRead(d, meta)
}

func resourceArmPowerBIEmbeddedDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PowerBI.CapacityClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PowerBIEmbeddedID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("Error deleting PowerBI Embedded %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deleting PowerBI Embedded %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}

	return nil
}
