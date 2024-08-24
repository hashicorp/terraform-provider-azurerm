// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/transparentdataencryptions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
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

			"enclave_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"identity": commonschema.UserAssignedIdentityComputed(),

			"transparent_data_encryption_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"transparent_data_encryption_key_vault_key_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"transparent_data_encryption_key_automatic_rotation_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceMsSqlDatabaseRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.DatabasesClient
	transparentEncryptionClient := meta.(*clients.Client).MSSQL.TransparentDataEncryptionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	mssqlServerId := d.Get("server_id").(string)
	serverId, err := commonids.ParseSqlServerID(mssqlServerId)
	if err != nil {
		return err
	}

	databaseId := commonids.NewSqlDatabaseID(serverId.SubscriptionId, serverId.ResourceGroupName, serverId.ServerName, name)

	resp, err := client.Get(ctx, databaseId, databases.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", databaseId)
		}

		return fmt.Errorf("making Read request on AzureRM %s: %+v", databaseId, err)
	}

	d.SetId(databaseId.ID())
	d.Set("name", name)
	d.Set("server_id", mssqlServerId)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("collation", props.Collation)
			d.Set("elastic_pool_id", props.ElasticPoolId)
			d.Set("license_type", string(pointer.From(props.LicenseType)))
			d.Set("read_replica_count", props.HighAvailabilityReplicaCount)
			d.Set("sku_name", props.CurrentServiceObjectiveName)
			d.Set("zone_redundant", props.ZoneRedundant)
			d.Set("transparent_data_encryption_key_vault_key_id", props.EncryptionProtector)
			d.Set("transparent_data_encryption_key_automatic_rotation_enabled", props.EncryptionProtectorAutoRotation)

			maxSizeGb := int64(0)
			if props.MaxSizeBytes != nil {
				maxSizeGb = (pointer.From(props.MaxSizeBytes)) / int64(1073741824)
			}
			d.Set("max_size_gb", maxSizeGb)

			readScale := databases.DatabaseReadScaleDisabled
			if props.ReadScale != nil {
				readScale = pointer.From(props.ReadScale)
			}
			d.Set("read_scale", readScale == databases.DatabaseReadScaleEnabled)

			enclaveType := ""
			if props.PreferredEnclaveType != nil && *props.PreferredEnclaveType != databases.AlwaysEncryptedEnclaveTypeDefault {
				enclaveType = string(*props.PreferredEnclaveType)
			}
			d.Set("enclave_type", enclaveType)

			storageAccountType := string(databases.BackupStorageRedundancyGeo)
			if props.CurrentBackupStorageRedundancy != nil {
				storageAccountType = string(pointer.From(props.CurrentBackupStorageRedundancy))
			}
			d.Set("storage_account_type", storageAccountType)
		}

		identity, err := identity.FlattenUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if err := d.Set("identity", identity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		tde, err := transparentEncryptionClient.Get(ctx, databaseId)
		if err != nil {
			return fmt.Errorf("while retrieving Transparent Data Encryption state for %s: %+v", databaseId, err)
		}
		if model := tde.Model; model != nil {
			if props := model.Properties; props != nil {
				d.Set("transparent_data_encryption_enabled", props.State == transparentdataencryptions.TransparentDataEncryptionStateEnabled)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}
