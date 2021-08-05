package postgres

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2021-06-01/postgresqlflexibleservers"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourcePostgresqlFlexibleServer() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceArmPostgresqlFlexibleServerRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"administrator_login": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"storage_mb": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"version": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"delegated_subnet_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"backup_retention_days": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"cmk_enabled": {
				Type:       pluginsdk.TypeString,
				Computed:   true,
				Deprecated: "This attribute has been removed from the API and will be removed in version 3.0 of the provider.",
			},

			"fqdn": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmPostgresqlFlexibleServerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Postgres.FlexibleServersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := parse.NewFlexibleServerID(subscriptionId, resourceGroup, name)

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Postgresqlflexibleservers Server %q does not exist", id.Name)
		}
		return fmt.Errorf("retrieving Postgresqlflexibleservers Server %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	// `cmk_enabled` has been removed from API since 2021-06-01
	// and should be removed in version 3.0 of the provider.
	d.Set("cmk_enabled", "")

	if props := resp.ServerProperties; props != nil {
		d.Set("administrator_login", props.AdministratorLogin)
		d.Set("version", props.Version)
		d.Set("fqdn", props.FullyQualifiedDomainName)

		if storage := props.Storage; storage != nil && storage.StorageSizeGB != nil {
			d.Set("storage_mb", (*props.Storage.StorageSizeGB * 1024))
		}

		if backup := props.Backup; backup != nil {
			d.Set("backup_retention_days", props.Backup.BackupRetentionDays)
		}

		if network := props.Network; network != nil {
			d.Set("delegated_subnet_id", network.DelegatedSubnetResourceID)
			d.Set("public_network_access_enabled", network.PublicNetworkAccess == postgresqlflexibleservers.ServerPublicNetworkAccessStateEnabled)
		}
	}

	sku, err := flattenFlexibleServerSku(resp.Sku)
	if err != nil {
		return fmt.Errorf("flattening `sku_name` for PostgreSQL Flexible Server %s (Resource Group %q): %v", id.Name, id.ResourceGroup, err)
	}

	d.Set("sku_name", sku)
	return tags.FlattenAndSet(d, resp.Tags)
}
