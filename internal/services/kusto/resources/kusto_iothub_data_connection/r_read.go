// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kusto_iothub_data_connection

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/dataconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKustoIotHubDataConnectionRead(d *pluginsdk.ResourceData, meta any) error {
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
		if dataConnection, ok := resp.Model.(dataconnections.IotHubDataConnection); ok {
			d.Set("location", location.NormalizeNilable(dataConnection.Location))
			if props := dataConnection.Properties; props != nil {
				iotHubId := ""
				if parsedIoTHubId, err := parse.IotHubIDInsensitively(props.IotHubResourceId); err == nil {
					iotHubId = parsedIoTHubId.ID()
				}
				d.Set("iothub_id", iotHubId)
				d.Set("consumer_group", props.ConsumerGroup)
				d.Set("table_name", props.TableName)
				d.Set("mapping_rule_name", props.MappingRuleName)
				d.Set("data_format", pointer.FromEnum(props.DataFormat))
				d.Set("database_routing_type", pointer.FromEnum(props.DatabaseRouting))
				d.Set("shared_access_policy_name", props.SharedAccessPolicyName)
				d.Set("event_system_properties", pluginsdk.FlattenSlice(props.EventSystemProperties))
				d.Set("retrieval_start_date", pointer.From(props.RetrievalStartDate))
			}
		}
	}

	return nil
}
