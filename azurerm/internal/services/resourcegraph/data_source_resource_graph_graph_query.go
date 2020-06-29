package resourcegraph

import (
	"fmt"
	"log"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

func dataSourceResourceGraphGraphQuery() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmResourceGraphGraphQueryRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),
			"resource_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"query": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"result_kind": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"time_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmResourceGraphGraphQueryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ResourceGraph.GraphQueryClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	resourceName := d.Get("resource_name").(string)

	resp, err := client.Get(ctx, resourceGroup, resourceName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] ResourceGraph %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failure in retrieving ResourceGraph GraphQuery (Resource Group %q / resourceName %q): %+v", resourceGroup, resourceName, err)
	}
	if id := resp.ID; id != nil {
		d.SetId(*resp.ID)
	}
	d.Set("resource_group_name", resourceGroup)
	d.Set("resource_name", resourceName)
	if name := resp.Name; name != nil {
		d.Set("name", name)
	}
	if props := resp.GraphQueryProperties; props != nil {
		d.Set("query", props.Query)
		d.Set("description", props.Description)
		d.Set("result_kind", props.ResultKind)
		d.Set("time_modified", props.TimeModified.Format(time.RFC3339))
	}
	return nil
}
