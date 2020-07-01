package resourcegraph

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/resourcegraph/mgmt/2018-09-01/resourcegraph"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resourcegraph/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmResourceGraphGraphQuery() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmResourceGraphGraphQueryCreate,
		Read:   resourceArmResourceGraphGraphQueryRead,
		Update: resourceArmResourceGraphGraphQueryUpdate,
		Delete: resourceArmResourceGraphGraphQueryDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ResourceGraphGraphQueryID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"resource_group_name": azure.SchemaResourceGroupName(),
			"resource_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"query": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"result_kind": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"time_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tags.Schema(),
		},
	}
}
func resourceArmResourceGraphGraphQueryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ResourceGraph.GraphQueryClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	resourceName := d.Get("resource_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, resourceName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("failure in checking for present of existing ResourceGraph GraphQuery (Resource Group %q / resourceName %q): %+v", resourceGroup, resourceName, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_resource_graph_graph_query", *existing.ID)
		}
	}

	props := resourcegraph.GraphQueryResource{
		Location: utils.String("global"),
		GraphQueryProperties: &resourcegraph.GraphQueryProperties{
			Description: utils.String(d.Get("description").(string)),
			Query:       utils.String(d.Get("query").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, resourceName, props); err != nil {
		return fmt.Errorf("failure in creating/updating ResourceGraph GraphQuery (Resource Group %q / resourceName %q / properties %v): %+v", resourceGroup, resourceName, props, err)
	}

	resp, err := client.Get(ctx, resourceGroup, resourceName)
	if err != nil {
		return fmt.Errorf("failure in retrieving ResourceGraph GraphQuery (Resource Group %q / resourceName %q): %+v", resourceGroup, resourceName, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("cannot read ResourceGraph GraphQuery (Resource Group %q / resourceName %q): %+v", resourceGroup, resourceName, err)
	}

	d.SetId(*resp.ID)
	return resourceArmResourceGraphGraphQueryRead(d, meta)
}

func resourceArmResourceGraphGraphQueryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ResourceGraph.GraphQueryClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ResourceGraphGraphQueryID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ResourceName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] ResourceGraph %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failure in retrieving ResourceGraph GraphQuery (Resource Group %q / resourceName %q): %+v", id.ResourceGroup, id.ResourceName, err)
	}
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("resource_name", id.ResourceName)
	if name := resp.Name; name != nil {
		d.Set("name", name)
	}

	if props := resp.GraphQueryProperties; props != nil {
		d.Set("query", props.Query)
		d.Set("description", props.Description)
		d.Set("result_kind", props.ResultKind)
		d.Set("time_modified", props.TimeModified.Format(time.RFC3339))
	}
	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmResourceGraphGraphQueryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ResourceGraph.GraphQueryClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	resourceName := d.Get("resource_name").(string)

	body := resourcegraph.GraphQueryUpdateParameters{
		GraphQueryPropertiesUpdateParameters: &resourcegraph.GraphQueryPropertiesUpdateParameters{
			Description: utils.String(d.Get("description").(string)),
			Query:       utils.String(d.Get("query").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if _, err := client.Update(ctx, resourceGroup, resourceName, body); err != nil {
		return fmt.Errorf("failure in creating/updating ResourceGraph GraphQuery (Resource Group %q / resourceName %q / body %v): %+v", resourceGroup, resourceName, body, err)
	}

	resp, err := client.Get(ctx, resourceGroup, resourceName)
	if err != nil {
		return fmt.Errorf("failure in ResourceGraph GraphQuery (Resource Group %q / resourceName %q): %+v", resourceGroup, resourceName, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("cannot read ResourceGraph GraphQuery (Resource Group %q / resourceName %q): %+v", resourceGroup, resourceName, err)
	}

	d.SetId(*resp.ID)
	return resourceArmResourceGraphGraphQueryRead(d, meta)
}

func resourceArmResourceGraphGraphQueryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ResourceGraph.GraphQueryClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ResourceGraphGraphQueryID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.ResourceName); err != nil {
		return fmt.Errorf("failure in deleting ResourceGraph GraphQuery (Resource Group %q / resourceName %q): %+v", id.ResourceGroup, id.ResourceName, err)
	}
	return nil
}
