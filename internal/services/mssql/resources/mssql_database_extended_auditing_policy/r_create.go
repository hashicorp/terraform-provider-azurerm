// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_database_extended_auditing_policy

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/blobauditing"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlDatabaseExtendedAuditingPolicyCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
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
