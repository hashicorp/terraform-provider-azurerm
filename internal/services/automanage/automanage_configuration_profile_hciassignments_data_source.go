package automanage

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"time"
)

func dataSourceAutomanageConfigurationProfileHCIAssignment() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAutomanageConfigurationProfileHCIAssignmentRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceAutomanageConfigurationProfileHCIAssignmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automanage.ConfigurationProfileHCIAssignmentClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	clusterName := d.Get("cluster_name").(string)

	resp, err := client.Get(ctx, resourceGroup, clusterName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Automanage ConfigurationProfileHCIAssignment %q (Resource Group %q / clusterName %q) does not exist", name, resourceGroup, clusterName)
		}
		return fmt.Errorf("retrieving Automanage ConfigurationProfileHCIAssignment %q (Resource Group %q / clusterName %q): %+v", name, resourceGroup, clusterName, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Automanage ConfigurationProfileHCIAssignment %q (Resource Group %q / clusterName %q) ID", name, resourceGroup, clusterName)
	}

	d.SetId(*resp.ID)
	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("cluster_name", clusterName)
	return nil
}
