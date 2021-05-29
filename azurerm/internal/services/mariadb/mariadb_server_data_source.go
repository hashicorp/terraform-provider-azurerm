package mariadb

import (
	"fmt"
	"regexp"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceMariaDbServer() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceMariaDbServerRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9]{3,50}$"),
					"MariaDB server name must be 3 - 50 characters long, contain only letters, numbers and hyphens.",
				),
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"administrator_login": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"version": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"storage_profile": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"storage_mb": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"backup_retention_days": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"geo_redundant_backup": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"auto_grow": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"ssl_enforcement": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"fqdn": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceMariaDbServerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MariaDB.ServersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("MariaDB Server %q (Resource Group %q) was not found", name, resourceGroup)
		}

		return fmt.Errorf("retrieving MariaDB Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("retrieving MariaDB Server %q (Resource Group %q): `id` was nil", name, resourceGroup)
	}

	d.SetId(*resp.ID)

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku_name", sku.Name)
	}

	if properties := resp.ServerProperties; properties != nil {
		d.Set("administrator_login", properties.AdministratorLogin)
		d.Set("version", string(properties.Version))
		d.Set("ssl_enforcement", string(properties.SslEnforcement))
		d.Set("fqdn", properties.FullyQualifiedDomainName)

		if err := d.Set("storage_profile", flattenMariaDbStorageProfile(properties.StorageProfile)); err != nil {
			return fmt.Errorf("Error setting `storage_profile`: %+v", err)
		}
	}
	return tags.FlattenAndSet(d, resp.Tags)
}
