package mssql

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func dataSourceArmMSSQLManagedDatabase() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmMSSQLManagedDatabaseRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"managed_instance_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group": {
				Type:     schema.TypeString,
				Required: true,
			},

			"managed_instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"collation": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"earliest_restore_point": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"default_secondary_location": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func dataSourceArmMSSQLManagedDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ManagedDatabasesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	managedInstanceName := d.Get("managed_instance_name").(string)
	resourceGroup := d.Get("resource_group").(string)

	resp, err := client.Get(ctx, resourceGroup, managedInstanceName, name)
	if err != nil {
		return fmt.Errorf("while reading managed SQL Database %s: %v", name, err)
	}

	if id := resp.ID; id != nil {
		d.SetId(*id)
	}

	managedInstanceId, _ := azure.GetSQLResourceParentId(d.Id())

	d.Set("name", name)
	d.Set("resource_group", resourceGroup)
	d.Set("type", resp.Type)
	d.Set("managed_instance_id", managedInstanceId)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.ManagedDatabaseProperties; props != nil {
		d.Set("collation", props.Collation)
		d.Set("status", props.Status)
		if props.CreationDate != nil && props.CreationDate.String() != "" {
			d.Set("creation_date", props.CreationDate.String())
		}
		if props.EarliestRestorePoint != nil && props.EarliestRestorePoint.String() != "" {
			d.Set("earliest_restore_point", props.EarliestRestorePoint.String())
		}
		d.Set("default_secondary_location", props.DefaultSecondaryLocation)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
