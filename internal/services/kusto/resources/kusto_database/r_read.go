// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kusto_database

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/databases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKustoDatabaseRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Kusto.DatabasesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseKustoDatabaseID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if resp.Model == nil {
		return fmt.Errorf("retrieving %s: response was nil", *id)
	}

	database, ok := resp.Model.(databases.ReadWriteDatabase)
	if !ok {
		return fmt.Errorf("%s was not a Read/Write Database", *id)
	}

	d.Set("name", id.KustoDatabaseName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("cluster_name", id.KustoClusterName)
	d.Set("location", location.NormalizeNilable(database.Location))

	if props := database.Properties; props != nil {
		d.Set("hot_cache_period", props.HotCachePeriod)
		d.Set("soft_delete_period", props.SoftDeletePeriod)

		if statistics := props.Statistics; statistics != nil {
			d.Set("size", statistics.Size)
		}
	}

	return nil
}
