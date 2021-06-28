package powerbi

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/powerbidedicated/mgmt/2021-01-01/powerbidedicated"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/powerbi/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/powerbi/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourcePowerBIEmbedded() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePowerBIEmbeddedCreate,
		Read:   resourcePowerBIEmbeddedRead,
		Update: resourcePowerBIEmbeddedUpdate,
		Delete: resourcePowerBIEmbeddedDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.EmbeddedID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.EmbeddedName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"sku_name": {
				Type:     pluginsdk.TypeString,
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
				Type:     pluginsdk.TypeSet,
				Required: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validate.EmbeddedAdministratorName,
				},
			},

			"mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(powerbidedicated.ModeGen1),
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(powerbidedicated.ModeGen1),
					string(powerbidedicated.ModeGen2),
				}, false),
			},

			"tags": tags.Schema(),
		},
	}
}

func resourcePowerBIEmbeddedCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PowerBI.CapacityClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	existing, err := client.GetDetails(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("Error checking for present of existing PowerBI Embedded %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}
	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_powerbi_embedded", *existing.ID)
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	administrators := d.Get("administrators").(*pluginsdk.Set).List()
	skuName := d.Get("sku_name").(string)
	mode := d.Get("mode").(string)
	t := d.Get("tags").(map[string]interface{})

	parameters := powerbidedicated.DedicatedCapacity{
		Location: utils.String(location),
		DedicatedCapacityProperties: &powerbidedicated.DedicatedCapacityProperties{
			Administration: &powerbidedicated.DedicatedCapacityAdministrators{
				Members: utils.ExpandStringSlice(administrators),
			},
			Mode: powerbidedicated.Mode(mode),
		},
		Sku: &powerbidedicated.CapacitySku{
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

	return resourcePowerBIEmbeddedRead(d, meta)
}

func resourcePowerBIEmbeddedRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PowerBI.CapacityClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EmbeddedID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetDetails(ctx, id.ResourceGroup, id.CapacityName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] PowerBI Embedded %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading PowerBI Embedded %q (Resource Group %q): %+v", id.CapacityName, id.ResourceGroup, err)
	}

	d.Set("name", id.CapacityName)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if props := resp.DedicatedCapacityProperties; props != nil {
		if err := d.Set("administrators", utils.FlattenStringSlice(props.Administration.Members)); err != nil {
			return fmt.Errorf("Error setting `administration`: %+v", err)
		}

		d.Set("mode", props.Mode)
	}

	skuName := ""
	if resp.Sku != nil && resp.Sku.Name != nil {
		skuName = *resp.Sku.Name
	}
	d.Set("sku_name", skuName)

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourcePowerBIEmbeddedUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PowerBI.CapacityClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	administrators := d.Get("administrators").(*pluginsdk.Set).List()
	skuName := d.Get("sku_name").(string)
	mode := d.Get("mode").(string)
	t := d.Get("tags").(map[string]interface{})

	parameters := powerbidedicated.DedicatedCapacityUpdateParameters{
		DedicatedCapacityMutableProperties: &powerbidedicated.DedicatedCapacityMutableProperties{
			Administration: &powerbidedicated.DedicatedCapacityAdministrators{
				Members: utils.ExpandStringSlice(administrators),
			},
			Mode: powerbidedicated.Mode(mode),
		},
		Sku: &powerbidedicated.CapacitySku{
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

	return resourcePowerBIEmbeddedRead(d, meta)
}

func resourcePowerBIEmbeddedDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PowerBI.CapacityClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EmbeddedID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.CapacityName)
	if err != nil {
		return fmt.Errorf("Error deleting PowerBI Embedded %q (Resource Group %q): %+v", id.CapacityName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deleting PowerBI Embedded %q (Resource Group %q): %+v", id.CapacityName, id.ResourceGroup, err)
		}
	}

	return nil
}
