// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sql

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceArmSqlMiServer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmSqlMiServerRead,

		DeprecationMessage: "The `azurerm_sql_managed_instance` data source is deprecated and will be removed in version 4.0 of the AzureRM provider. Please use the `azurerm_mssql_managed_instance` data source instead.",

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.ValidateMsSqlServerName,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"sku_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"administrator_login": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"vcores": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"storage_size_in_gb": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"license_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"collation": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"public_data_endpoint_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"minimum_tls_version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"proxy_override": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"timezone_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"dns_zone_partner_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"identity": commonschema.SystemAssignedIdentityComputed(),

			"storage_account_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func dataSourceArmSqlMiServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.ManagedInstancesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewManagedInstanceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("reading SQL Managed Instance %q: %v", id.ID(), err)
	}

	d.SetId(id.ID())

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	d.Set("location", location.NormalizeNilable(resp.Location))

	if sku := resp.Sku; sku != nil {
		d.Set("sku_name", sku.Name)
	}

	if err := d.Set("identity", flattenManagedInstanceIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	if props := resp.ManagedInstanceProperties; props != nil {
		d.Set("license_type", string(props.LicenseType))
		d.Set("administrator_login", props.AdministratorLogin)
		d.Set("subnet_id", props.SubnetID)
		d.Set("storage_size_in_gb", props.StorageSizeInGB)
		d.Set("vcores", props.VCores)
		d.Set("fqdn", props.FullyQualifiedDomainName)
		d.Set("collation", props.Collation)
		d.Set("public_data_endpoint_enabled", props.PublicDataEndpointEnabled)
		d.Set("minimum_tls_version", props.MinimalTLSVersion)
		d.Set("proxy_override", props.ProxyOverride)
		d.Set("timezone_id", props.TimezoneID)
		d.Set("storage_account_type", props.StorageAccountType)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
