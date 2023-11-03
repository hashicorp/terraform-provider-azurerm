// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/databases" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceMsSqlDatabase() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceMsSqlDatabaseRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.ValidateMsSqlDatabaseName,
			},

			"server_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.ServerID,
			},

			"collation": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"elastic_pool_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"license_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"max_size_gb": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"read_replica_count": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"read_scale": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"storage_account_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"zone_redundant": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceMsSqlDatabaseRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.DatabasesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	mssqlServerId := d.Get("server_id").(string)
	serverId, err := parse.ServerID(mssqlServerId)
	if err != nil {
		return err
	}

	databaseId := databases.DatabaseId{
		SubscriptionId:    serverId.SubscriptionId,
		ResourceGroupName: serverId.ResourceGroup,
		ServerName:        serverId.Name,
		DatabaseName:      name,
	}

	resp, err := client.Get(ctx, databaseId, databases.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("database %q %s was not found", name, serverId.ID())
		}

		return fmt.Errorf("making Read request on AzureRM Database %q %s: %+v", name, serverId.ID(), err)
	}

	d.SetId(parse.NewDatabaseID(serverId.SubscriptionId, serverId.ResourceGroup, serverId.Name, name).ID())
	d.Set("name", name)
	d.Set("server_id", mssqlServerId)

	if props := resp.Model.Properties; props != nil {
		d.Set("collation", props.Collation)
		d.Set("elastic_pool_id", props.ElasticPoolId)
		d.Set("license_type", props.LicenseType)
		d.Set("read_replica_count", props.HighAvailabilityReplicaCount)
		d.Set("sku_name", props.CurrentServiceObjectiveName)
		d.Set("zone_redundant", props.ZoneRedundant)

		if props.MaxSizeBytes != nil {
			d.Set("max_size_gb", int32((pointer.From(props.MaxSizeBytes))/int64(1073741824)))
		}

		readScale := databases.DatabaseReadScaleDisabled
		if props.ReadScale != nil {
			readScale = pointer.From(props.ReadScale)
		}
		d.Set("read_scale", readScale == databases.DatabaseReadScaleEnabled)

		storageAccountType := databases.BackupStorageRedundancyGeo
		if props.CurrentBackupStorageRedundancy != nil {
			storageAccountType = pointer.From(props.CurrentBackupStorageRedundancy)
		}
		d.Set("storage_account_type", storageAccountType)
	} else {
		log.Print("Model Properties were nil")
	}

	return tags.FlattenAndSet(d, tags.ExpandFromPointer(resp.Model.Tags))
}
