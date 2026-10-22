// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_database

import (
	"context"
	"fmt"
	"math"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/replicationlinks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlDatabaseRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MSSQL.DatabasesClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseSqlDatabaseID(d.Id())
	if err != nil {
		return err
	}

	serverId := commonids.NewSqlServerID(id.SubscriptionId, id.ResourceGroupName, id.ServerName)
	d.Set("server_id", serverId.ID())

	resp, err := client.Get(ctx, pointer.From(id), databases.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return resourceMssqlDatabaseSetFlatten(d, id, resp.Model, meta.(*clients.Client))
}

func resourceMsSqlDatabaseImporter(ctx context.Context, d *pluginsdk.ResourceData, meta any) ([]*pluginsdk.ResourceData, error) {
	replicationLinksClient := meta.(*clients.Client).MSSQL.ReplicationLinksClient
	client := meta.(*clients.Client).MSSQL.DatabasesClient
	resourcesClient := meta.(*clients.Client).Resource.ResourcesClient

	id, err := commonids.ParseSqlDatabaseID(d.Id())
	if err != nil {
		return nil, err
	}

	// NOTE: The service default is actually nil/empty which indicates enclave is disabled. the value `Default` is NOT the default.
	var enclaveType databases.AlwaysEncryptedEnclaveType
	if v, ok := d.GetOk("enclave_type"); ok && v.(string) != "" {
		enclaveType = databases.AlwaysEncryptedEnclaveType(v.(string))
	}
	d.Set("enclave_type", enclaveType)

	partnerDatabases, err := helper.FindDatabaseReplicationPartners(ctx, client, replicationLinksClient, resourcesClient, *id, enclaveType, []replicationlinks.ReplicationRole{replicationlinks.ReplicationRolePrimary})
	if err != nil {
		return nil, err
	}

	if len(partnerDatabases) > 0 {
		partnerDatabase := partnerDatabases[0]

		partnerDatabaseId, err := commonids.ParseSqlDatabaseIDInsensitively(*partnerDatabase.Id)
		if err != nil {
			return nil, fmt.Errorf("parsing ID for Replication Partner Database %q: %+v", *partnerDatabase.Id, err)
		}

		d.Set("create_mode", string(databases.CreateModeSecondary))
		d.Set("creation_source_database_id", partnerDatabaseId.ID())

		return []*pluginsdk.ResourceData{d}, nil
	}

	d.Set("create_mode", string(databases.CreateModeDefault))

	return []*pluginsdk.ResourceData{d}, nil
}

func calculateMaxSizeBytes(v float64) (*int64, error) {
	var maxSizeBytes int64
	if isPositiveInteger(v) {
		maxSizeBytes = int64(math.Round(v)) * 1073741824
	} else {
		// When `max_size_gb` is below 1GB, the API only accepts 100MB and 500MB. 100MB is 104857600, and 500MB is 524288000.
		// Since 100MB != 0.1 * 1073741824 and 500MB != 0.5 * 1073741824, directly specify the Bytes when `max_size_gb` is specified `0.1` or `0.5`.
		switch int64(math.Round(v * 10)) {
		case 1:
			maxSizeBytes = 104857600
		case 5:
			maxSizeBytes = 524288000
		default:
			return nil, fmt.Errorf("`max_size_gb` allows `0.1`, `0.5` and positive integers greater than or equal to 1")
		}
	}
	return pointer.To(maxSizeBytes), nil
}

func isPositiveInteger(num float64) bool {
	// Determine whether `num` is a positive integer and >= 1
	return num >= 1 && num == math.Floor(num)
}
