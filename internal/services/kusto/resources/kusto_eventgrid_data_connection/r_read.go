// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kusto_eventgrid_data_connection

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/dataconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKustoEventGridDataConnectionRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Kusto.DataConnectionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dataconnections.ParseDataConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.DataConnectionName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("cluster_name", id.ClusterName)
	d.Set("database_name", id.DatabaseName)

	if resp.Model != nil {
		if dataConnection, ok := resp.Model.(dataconnections.EventGridDataConnection); ok {
			d.Set("location", location.NormalizeNilable(dataConnection.Location))
			if props := dataConnection.Properties; props != nil {
				d.Set("storage_account_id", props.StorageAccountResourceId)
				d.Set("eventhub_id", props.EventHubResourceId)
				d.Set("eventhub_consumer_group_name", props.ConsumerGroup)
				d.Set("skip_first_record", props.IgnoreFirstRecord)
				d.Set("blob_storage_event_type", pointer.FromEnum(props.BlobStorageEventType))
				d.Set("table_name", props.TableName)
				d.Set("mapping_rule_name", props.MappingRuleName)
				d.Set("data_format", pointer.FromEnum(props.DataFormat))
				d.Set("database_routing_type", pointer.FromEnum(props.DatabaseRouting))
				d.Set("eventgrid_event_subscription_id", props.EventGridResourceId)

				managedIdentityResourceId := ""
				if props.ManagedIdentityResourceId != nil && *props.ManagedIdentityResourceId != "" {
					managedIdentityResourceId = *props.ManagedIdentityResourceId
					clusterId, clusterIdErr := commonids.ParseKustoClusterIDInsensitively(managedIdentityResourceId)
					if clusterIdErr == nil {
						managedIdentityResourceId = clusterId.ID()
					} else {
						userAssignedIdentityId, userAssignedIdentityIdErr := commonids.ParseUserAssignedIdentityIDInsensitively(managedIdentityResourceId)
						if userAssignedIdentityIdErr == nil {
							managedIdentityResourceId = userAssignedIdentityId.ID()
						} else {
							return fmt.Errorf("parsing `managed_identity_id`: %+v; %+v", clusterIdErr, userAssignedIdentityIdErr)
						}
					}
				}

				d.Set("managed_identity_id", managedIdentityResourceId)
			}
		}
	}

	return nil
}
