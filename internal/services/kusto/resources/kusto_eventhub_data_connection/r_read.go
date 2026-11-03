// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kusto_eventhub_data_connection

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

func resourceKustoEventHubDataConnectionRead(d *pluginsdk.ResourceData, meta any) error {
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
		return fmt.Errorf("retrieving Kusto Event Hub Data Connection %q (Resource Group %q, Cluster %q, Database %q): %+v", id.DataConnectionName, id.ResourceGroupName, id.ClusterName, id.DatabaseName, err)
	}

	d.Set("name", id.DataConnectionName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("cluster_name", id.ClusterName)
	d.Set("database_name", id.DatabaseName)

	if resp.Model != nil {
		if dataConnection, ok := resp.Model.(dataconnections.EventHubDataConnection); ok {
			if loc := dataConnection.Location; loc != nil {
				d.Set("location", location.Normalize(*loc))
			}

			if props := dataConnection.Properties; props != nil {
				d.Set("eventhub_id", props.EventHubResourceId)
				d.Set("consumer_group", props.ConsumerGroup)
				d.Set("table_name", props.TableName)
				d.Set("mapping_rule_name", props.MappingRuleName)
				d.Set("data_format", pointer.FromEnum(props.DataFormat))
				d.Set("database_routing_type", pointer.FromEnum(props.DatabaseRouting))
				d.Set("compression", pointer.FromEnum(props.Compression))
				d.Set("event_system_properties", props.EventSystemProperties)
				d.Set("retrieval_start_date", pointer.From(props.RetrievalStartDate))

				identityId := ""
				if props.ManagedIdentityResourceId != nil {
					identityId = *props.ManagedIdentityResourceId
					clusterId, clusterIdErr := commonids.ParseKustoClusterIDInsensitively(identityId)
					if clusterIdErr == nil {
						identityId = clusterId.ID()
					} else {
						userAssignedIdentityId, userAssignedIdentityIdErr := commonids.ParseUserAssignedIdentityIDInsensitively(identityId)
						if userAssignedIdentityIdErr == nil {
							identityId = userAssignedIdentityId.ID()
						} else {
							return fmt.Errorf("parsing `identity_id`: %+v; %+v", clusterIdErr, userAssignedIdentityIdErr)
						}
					}
				}
				d.Set("identity_id", identityId)
			}
		}
	}

	return nil
}
