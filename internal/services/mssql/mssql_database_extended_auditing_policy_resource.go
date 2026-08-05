// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/blobauditing"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlDatabaseExtendedAuditingPolicy() *pluginsdk.Resource {
	r := &pluginsdk.Resource{
		Create: resourceMsSqlDatabaseExtendedAuditingPolicyCreateUpdate,
		Read:   resourceMsSqlDatabaseExtendedAuditingPolicyRead,
		Update: resourceMsSqlDatabaseExtendedAuditingPolicyCreateUpdate,
		Delete: resourceMsSqlDatabaseExtendedAuditingPolicyDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.DatabaseExtendedAuditingPolicyID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"database_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateSqlDatabaseID,
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"blob_storage_endpoint": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsURLWithHTTPS,
			},

			"storage_account_access_key": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"storage_account_access_key_is_secondary": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"retention_in_days": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntBetween(0, 3285),
			},

			"log_monitoring_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}

	if !features.FivePointOh() {
		r.Schema["storage_endpoint"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Computed:      true,
			ValidateFunc:  validation.IsURLWithHTTPS,
			ConflictsWith: []string{"blob_storage_endpoint"},
			Deprecated:    "`storage_endpoint` is deprecated in favour of `blob_storage_endpoint` and will be removed in version 5.0 of the AzureRM provider",
		}

		r.Schema["blob_storage_endpoint"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Computed:      true,
			ValidateFunc:  validation.IsURLWithHTTPS,
			ConflictsWith: []string{"storage_endpoint"},
		}
	}

	return r
}

func resourceMsSqlDatabaseExtendedAuditingPolicyCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.BlobAuditingPoliciesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dbId, err := commonids.ParseSqlDatabaseID(d.Get("database_id").(string))
	if err != nil {
		return err
	}

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.ExtendedDatabaseBlobAuditingPoliciesGet(ctx, *dbId)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of existing %s: %+v", dbId, err)
				}
			}

			// if state is not disabled, we should import it.
			if existing.Model != nil && existing.Model.Id != nil && *existing.Model.Id != "" && existing.Model.Properties != nil && existing.Model.Properties.State != blobauditing.BlobAuditingPolicyStateDisabled {
				return tf.ImportAsExistsError("azurerm_mssql_database_extended_auditing_policy", *existing.Model.Id)
			}
		}
	}

	params := blobauditing.ExtendedDatabaseBlobAuditingPolicy{
		Properties: &blobauditing.ExtendedDatabaseBlobAuditingPolicyProperties{
			StorageEndpoint:             pointer.To(d.Get("blob_storage_endpoint").(string)),
			IsStorageSecondaryKeyInUse:  pointer.To(d.Get("storage_account_access_key_is_secondary").(bool)),
			RetentionDays:               pointer.To(int64(d.Get("retention_in_days").(int))),
			IsAzureMonitorTargetEnabled: pointer.To(d.Get("log_monitoring_enabled").(bool)),
		},
	}

	if !features.FivePointOh() {
		if !pluginsdk.IsExplicitlyNullInConfig(d, "storage_endpoint") {
			params.Properties.StorageEndpoint = pointer.To(d.Get("storage_endpoint").(string))
		}
	}

	if d.Get("enabled").(bool) {
		params.Properties.State = blobauditing.BlobAuditingPolicyStateEnabled
	} else {
		params.Properties.State = blobauditing.BlobAuditingPolicyStateDisabled
	}

	if v, ok := d.GetOk("storage_account_access_key"); ok {
		params.Properties.StorageAccountAccessKey = pointer.To(v.(string))
	}

	if _, err = client.ExtendedDatabaseBlobAuditingPoliciesCreateOrUpdate(ctx, *dbId, params); err != nil {
		return fmt.Errorf("creating extended auditing policy for %s: %+v", dbId, err)
	}

	read, err := client.ExtendedDatabaseBlobAuditingPoliciesGet(ctx, *dbId)
	if err != nil {
		return fmt.Errorf("retrieving the extended auditing policy for %s: %+v", dbId, err)
	}

	if read.Model == nil || read.Model.Id == nil || pointer.From(read.Model.Id) == "" {
		return fmt.Errorf("the extended auditing policy ID for %s is 'nil' or 'empty'", dbId.String())
	}

	// TODO: update this to use the Database ID - requiring a State Migration
	readId, err := parse.DatabaseExtendedAuditingPolicyID(pointer.From(read.Model.Id))
	if err != nil {
		return err
	}

	d.SetId(readId.ID())

	return resourceMsSqlDatabaseExtendedAuditingPolicyRead(d, meta)
}

func resourceMsSqlDatabaseExtendedAuditingPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.BlobAuditingPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DatabaseExtendedAuditingPolicyID(d.Id())
	if err != nil {
		return err
	}

	dbId := commonids.NewSqlDatabaseID(id.SubscriptionId, id.ResourceGroup, id.ServerName, id.DatabaseName)

	resp, err := client.ExtendedDatabaseBlobAuditingPoliciesGet(ctx, dbId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading MsSql Database %s Extended Auditing Policy (MsSql Server Name %q / Resource Group %q): %s", id.DatabaseName, id.ServerName, id.ResourceGroup, err)
	}

	databaseId := commonids.NewSqlDatabaseID(id.SubscriptionId, id.ResourceGroup, id.ServerName, id.DatabaseName)
	d.Set("database_id", databaseId.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("blob_storage_endpoint", props.StorageEndpoint)
			if !features.FivePointOh() {
				d.Set("storage_endpoint", props.StorageEndpoint)
			}
			d.Set("storage_account_access_key_is_secondary", props.IsStorageSecondaryKeyInUse)
			d.Set("retention_in_days", props.RetentionDays)
			d.Set("log_monitoring_enabled", props.IsAzureMonitorTargetEnabled)
			d.Set("enabled", props.State == blobauditing.BlobAuditingPolicyStateEnabled)
		}
	}

	return nil
}

func resourceMsSqlDatabaseExtendedAuditingPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.BlobAuditingPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DatabaseExtendedAuditingPolicyID(d.Id())
	if err != nil {
		return err
	}

	dbId := commonids.NewSqlDatabaseID(id.SubscriptionId, id.ResourceGroup, id.ServerName, id.DatabaseName)

	params := blobauditing.ExtendedDatabaseBlobAuditingPolicy{
		Properties: &blobauditing.ExtendedDatabaseBlobAuditingPolicyProperties{
			State: blobauditing.BlobAuditingPolicyStateDisabled,
		},
	}

	if _, err = client.ExtendedDatabaseBlobAuditingPoliciesCreateOrUpdate(ctx, dbId, params); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
