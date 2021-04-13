package postgres

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/postgresql/mgmt/2020-02-14-preview/postgresqlflexibleservers"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourcePostgresqlFlexibleServer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmPostgresqlFlexibleServerRead,

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

			"administrator_login": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"sku_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"storage_mb": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"delegated_subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"identity": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"principal_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"backup_retention_days": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"cmk_enabled": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"high_availiblity_state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"public_network_access_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"standby_availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmPostgresqlFlexibleServerRead(d *schema.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Postgres.FlexibleServersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := parse.NewFlexibleServerID(subscriptionId, resourceGroup, name).ID()

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Postgresqlflexibleservers Server %q does not exist", name)
		}
		return fmt.Errorf("retrieving Postgresqlflexibleservers Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(id)
	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	if props := resp.ServerProperties; props != nil {
		d.Set("administrator_login", props.AdministratorLogin)
		d.Set("storage_mb", props.StorageProfile.StorageMB)
		d.Set("version", props.Version)
		d.Set("cmk_enabled", props.ByokEnforcement)
		d.Set("fqdn", props.FullyQualifiedDomainName)
		d.Set("public_network_access_enabled", props.PublicNetworkAccess == postgresqlflexibleservers.ServerPublicNetworkAccessStateEnabled)
		d.Set("high_availiblity_state", string(props.HaState))
		d.Set("standby_availability_zone", props.StandbyAvailabilityZone)

		if props.DelegatedSubnetArguments != nil {
			d.Set("delegated_subnet_id", props.DelegatedSubnetArguments.SubnetArmResourceID)
		}

		if storage := props.StorageProfile; storage != nil {
			d.Set("storage_mb", storage.StorageMB)
			d.Set("backup_retention_days", storage.BackupRetentionDays)
		}
	}

	sku, err := flattenFlexibleServerSku(resp.Sku)
	if err != nil {
		return fmt.Errorf("flattening `sku_name` for PostgreSQL Flexible Server %s (Resource Group %q): %v", name, resourceGroup, err)
	}

	d.Set("sku_name", sku)
	if err := d.Set("identity", flattenArmServerIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}
	return tags.FlattenAndSet(d, resp.Tags)
}
