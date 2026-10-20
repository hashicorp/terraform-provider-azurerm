// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kusto_iothub_data_connection

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/dataconnections"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKustoIotHubDataConnectionCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Kusto.DataConnectionsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := dataconnections.NewDataConnectionID(subscriptionId, d.Get("resource_group_name").(string), d.Get("cluster_name").(string), d.Get("database_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		resp, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(resp.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_kusto_iothub_data_connection", id.ID())
		}
	}

	iotHubDataConnectionProperties := expandKustoIotHubDataConnectionProperties(d)

	dataConnection := dataconnections.IotHubDataConnection{
		Location:   pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: iotHubDataConnectionProperties,
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, dataConnection, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceKustoIotHubDataConnectionRead(d, meta)
}

func expandKustoIotHubDataConnectionProperties(d *pluginsdk.ResourceData) *dataconnections.IotHubConnectionProperties {
	iotHubDataConnectionProperties := &dataconnections.IotHubConnectionProperties{
		IotHubResourceId:       d.Get("iothub_id").(string),
		ConsumerGroup:          d.Get("consumer_group").(string),
		SharedAccessPolicyName: d.Get("shared_access_policy_name").(string),
	}

	if tableName, ok := d.GetOk("table_name"); ok {
		iotHubDataConnectionProperties.TableName = pointer.To(tableName.(string))
	}

	if mappingRuleName, ok := d.GetOk("mapping_rule_name"); ok {
		iotHubDataConnectionProperties.MappingRuleName = pointer.To(mappingRuleName.(string))
	}

	if df, ok := d.GetOk("data_format"); ok {
		iotHubDataConnectionProperties.DataFormat = pointer.ToEnum[dataconnections.IotHubDataFormat](df.(string))
	}

	if databaseRouting, ok := d.GetOk("database_routing_type"); ok {
		iotHubDataConnectionProperties.DatabaseRouting = pointer.ToEnum[dataconnections.DatabaseRouting](databaseRouting.(string))
	}

	if eventSystemProperties, ok := d.GetOk("event_system_properties"); ok {
		iotHubDataConnectionProperties.EventSystemProperties = pluginsdk.ExpandStringSlice(eventSystemProperties.(*pluginsdk.Set).List())
	}

	if retrievalStartDate, ok := d.GetOk("retrieval_start_date"); ok {
		iotHubDataConnectionProperties.RetrievalStartDate = pointer.To(retrievalStartDate.(string))
	}

	return iotHubDataConnectionProperties
}
