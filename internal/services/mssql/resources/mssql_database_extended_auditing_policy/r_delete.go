// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_database_extended_auditing_policy

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/blobauditing"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlDatabaseExtendedAuditingPolicyDelete(d *pluginsdk.ResourceData, meta any) error {
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
