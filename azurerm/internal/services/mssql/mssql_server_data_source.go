package mssql

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceMsSqlServer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMsSqlServerRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"administrator_login": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"fully_qualified_domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"identity": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"principal_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"restorable_dropped_database_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceMsSqlServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ServersClient
	restorableDroppedDatabasesClient := meta.(*clients.Client).MSSQL.RestorableDroppedDatabasesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("sql Server %q was not found in Resource Group %q", name, resourceGroup)
		}

		return fmt.Errorf("retrieving Sql Server %q (Resource Group %q): %s", name, resourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("reading Ms Sql Server %q (Resource Group %q) ID is empty or nil", name, resourceGroup)
	}
	d.SetId(*resp.ID)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.ServerProperties; props != nil {
		d.Set("version", props.Version)
		d.Set("administrator_login", props.AdministratorLogin)
		d.Set("fully_qualified_domain_name", props.FullyQualifiedDomainName)
	}

	if err := d.Set("identity", flattenSqlServerIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	restorableResp, err := restorableDroppedDatabasesClient.ListByServer(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("listing SQL Server %s Restorable Dropped Databases: %v", name, err)
	}
	if err := d.Set("restorable_dropped_database_ids", flattenSqlServerRestorableDatabases(restorableResp)); err != nil {
		return fmt.Errorf("setting `restorable_dropped_database_ids`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
