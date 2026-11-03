// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kusto_eventgrid_data_connection

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

func resourceKustoEventGridDataConnectionCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Kusto.DataConnectionsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := dataconnections.NewDataConnectionID(subscriptionId, d.Get("resource_group_name").(string), d.Get("cluster_name").(string), d.Get("database_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			resp, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %s", id, err)
				}
			}

			if !response.WasNotFound(resp.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_kusto_eventgrid_data_connection", id.ID())
			}
		}
	}

	dataConnection := dataconnections.EventGridDataConnection{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &dataconnections.EventGridConnectionProperties{
			StorageAccountResourceId: d.Get("storage_account_id").(string),
			EventHubResourceId:       d.Get("eventhub_id").(string),
			ConsumerGroup:            d.Get("eventhub_consumer_group_name").(string),
			IgnoreFirstRecord:        pointer.To(d.Get("skip_first_record").(bool)),
		},
	}

	blobStorageEventType := dataconnections.BlobStorageEventType(d.Get("blob_storage_event_type").(string))
	dataConnection.Properties.BlobStorageEventType = &blobStorageEventType

	if tableName, ok := d.GetOk("table_name"); ok {
		dataConnection.Properties.TableName = pointer.To(tableName.(string))
	}

	if mappingRuleName, ok := d.GetOk("mapping_rule_name"); ok {
		dataConnection.Properties.MappingRuleName = pointer.To(mappingRuleName.(string))
	}

	if df, ok := d.GetOk("data_format"); ok {
		dataConnection.Properties.DataFormat = pointer.ToEnum[dataconnections.EventGridDataFormat](df.(string))
	}

	if databaseRouting, ok := d.GetOk("database_routing_type"); ok {
		dataConnection.Properties.DatabaseRouting = pointer.ToEnum[dataconnections.DatabaseRouting](databaseRouting.(string))
	}

	if eventGridRID, ok := d.GetOk("eventgrid_event_subscription_id"); ok {
		dataConnection.Properties.EventGridResourceId = pointer.To(eventGridRID.(string))
	}

	if managedIdentityRID, ok := d.GetOk("managed_identity_id"); ok {
		dataConnection.Properties.ManagedIdentityResourceId = pointer.To(managedIdentityRID.(string))
	}

	if d.IsNewResource() {
		if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, dataConnection, sdk.SetIDCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
		d.SetId(id.ID())
	} else {
		if err := client.CreateOrUpdateThenPoll(ctx, id, dataConnection); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	return resourceKustoEventGridDataConnectionRead(d, meta)
}
