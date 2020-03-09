package mssql

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmMsSqlDatabase() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmMsSqlDatabaseRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"mssql_server_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmMsSqlDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.DatabasesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	mssqlServerId := d.Get("mssql_server_id").(string)
	serverId, err := parse.MsSqlServerID(mssqlServerId)
	if err != nil {
		return err
	}

	serverClient := meta.(*clients.Client).MSSQL.ServersClient
	serverResp, err := serverClient.Get(ctx, serverId.ResourceGroup, serverId.Name)
	if err != nil {
		return fmt.Errorf("Error making Read request on MsSql Server  %q (Resource Group %q): %s", serverId.Name, serverId.ResourceGroup, err)
	}

	location := *serverResp.Location
	if location == "" {
		return fmt.Errorf("Error location is empty from making Read request on MsSql Server %q", serverId.Name)
	}

	resp, err := client.Get(ctx, serverId.ResourceGroup, serverId.Name, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Database %q (Resource Group %q, SQL Server %q) was not found", name, serverId.ResourceGroup, serverId.Name)
		}

		return fmt.Errorf("Error making Read request on AzureRM Database %s (Resource Group %q, SQL Server %q): %+v", name, serverId.ResourceGroup, serverId.Name, err)
	}

	if id := resp.ID; id != nil {
		d.SetId(*resp.ID)
	}
	d.Set("name", name)
	d.Set("mssql_server_id", mssqlServerId)

	return tags.FlattenAndSet(d, resp.Tags)
}
