// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kusto_eventhub_data_connection

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

func resourceKustoEventHubDataConnectionCreate(d *pluginsdk.ResourceData, meta any) error {
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
			return tf.ImportAsExistsError("azurerm_kusto_eventhub_data_connection", id.ID())
		}
	}

	location := location.Normalize(d.Get("location").(string))

	eventHubDataConnectionProperties := expandKustoEventHubDataConnectionProperties(d)

	dataConnection1 := dataconnections.EventHubDataConnection{
		Name:       &id.DataConnectionName,
		Location:   &location,
		Properties: eventHubDataConnectionProperties,
	}

	if databaseRouting, ok := d.GetOk("database_routing_type"); ok {
		dataConnection1.Properties.DatabaseRouting = pointer.ToEnum[dataconnections.DatabaseRouting](databaseRouting.(string))
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, dataConnection1, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	d.SetId(id.ID())

	return resourceKustoEventHubDataConnectionRead(d, meta)
}

func expandKustoEventHubDataConnectionProperties(d *pluginsdk.ResourceData) *dataconnections.EventHubConnectionProperties {
	eventHubConnectionProperties := &dataconnections.EventHubConnectionProperties{}

	if eventhubResourceID, ok := d.GetOk("eventhub_id"); ok {
		eventHubConnectionProperties.EventHubResourceId = eventhubResourceID.(string)
	}

	if consumerGroup, ok := d.GetOk("consumer_group"); ok {
		eventHubConnectionProperties.ConsumerGroup = consumerGroup.(string)
	}

	if tableName, ok := d.GetOk("table_name"); ok {
		eventHubConnectionProperties.TableName = pointer.To(tableName.(string))
	}

	if mappingRuleName, ok := d.GetOk("mapping_rule_name"); ok {
		eventHubConnectionProperties.MappingRuleName = pointer.To(mappingRuleName.(string))
	}

	if df, ok := d.GetOk("data_format"); ok {
		eventHubConnectionProperties.DataFormat = pointer.ToEnum[dataconnections.EventHubDataFormat](df.(string))
	}

	if compression, ok := d.GetOk("compression"); ok {
		eventHubConnectionProperties.Compression = pointer.ToEnum[dataconnections.Compression](compression.(string))
	}

	if eventSystemProperties, ok := d.GetOk("event_system_properties"); ok {
		props := make([]string, 0)
		for _, prop := range eventSystemProperties.([]any) {
			props = append(props, prop.(string))
		}
		eventHubConnectionProperties.EventSystemProperties = &props
	}

	if identityId, ok := d.GetOk("identity_id"); ok {
		eventHubConnectionProperties.ManagedIdentityResourceId = pointer.To(identityId.(string))
	}

	if retrievalStartDate, ok := d.GetOk("retrieval_start_date"); ok {
		eventHubConnectionProperties.RetrievalStartDate = pointer.To(retrievalStartDate.(string))
	}

	return eventHubConnectionProperties
}
