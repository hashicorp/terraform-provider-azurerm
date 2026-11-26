// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_database

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/transparentdataencryptions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing SQL database.",
		Read:        dataSourceMsSqlDatabaseRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description:  "The name of the MS SQL Database.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.ValidateMsSqlDatabaseName,
			},

			"server_id": {
				Description:  "The ID of the MS SQL Server on which to read the database.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.AsGeneratedID(commonids.ParseSqlServerIDInsensitively),
			},

			"collation": {
				Description: "The collation of the database.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"elastic_pool_id": {
				Description: "The ID of the elastic pool containing this database.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"license_type": {
				Description: "The license type to apply for this database.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"max_size_gb": {
				Description: "The max size of the database in gigabytes.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"read_replica_count": {
				Description: "The number of readonly secondary replicas associated with the database to which readonly application intent connections may be routed.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"read_scale": {
				Description: "If enabled, connections that have application intent set to readonly in their connection string may be routed to a readonly secondary replica.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"sku_name": {
				Description: "The name of the SKU of the database.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"storage_account_type": {
				Description: "The Storage Account type used to store backups for this database.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"zone_redundant": {
				Description: "Whether this database is zone redundant, which means the replicas of this database will be spread across multiple Availability Zones.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"enclave_type": {
				Description: "The type of enclave being used by the database.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"identity": commonschema.UserAssignedIdentityComputed(),

			"transparent_data_encryption_enabled": {
				Description: "Whether Transparent Data Encryption is enabled.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"transparent_data_encryption_key_vault_key_id": {
				Description: "The Key Vault Key URI to be used as the `Customer Managed Key` (CMK/BYOK) for the `Transparent Data Encryption` (TDE) layer.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"transparent_data_encryption_key_automatic_rotation_enabled": {
				Description: "Whether TDE automatically rotates the encryption Key to latest version.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceMsSqlDatabaseRead(d *pluginsdk.ResourceData, meta any) error {
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
			d.Set("license_type", pointer.FromEnum(props.LicenseType))
			d.Set("read_replica_count", props.HighAvailabilityReplicaCount)
			d.Set("sku_name", props.CurrentServiceObjectiveName)
			d.Set("zone_redundant", props.ZoneRedundant)
			d.Set("transparent_data_encryption_key_vault_key_id", props.EncryptionProtector)
			d.Set("transparent_data_encryption_key_automatic_rotation_enabled", props.EncryptionProtectorAutoRotation)

			maxSizeGb := int64(0)
			if props.MaxSizeBytes != nil {
				maxSizeGb = pointer.From(props.MaxSizeBytes) / int64(1073741824)
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
				storageAccountType = pointer.FromEnum(props.CurrentBackupStorageRedundancy)
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
