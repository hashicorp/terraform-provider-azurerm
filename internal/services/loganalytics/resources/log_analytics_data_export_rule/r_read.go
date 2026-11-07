// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_data_export_rule

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2021-11-01/eventhubs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/dataexport"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceOperationalinsightsDataExportRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).LogAnalytics.DataExportClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dataexport.ParseDataExportID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Log Analytics %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	d.Set("name", id.DataExportName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("workspace_resource_id", dataexport.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("export_rule_id", pointer.From(props.DataExportId))

			destinationId, err := flattenDataExportDestination(props.Destination)
			if err != nil {
				return fmt.Errorf("flattening destination ID error: %+v", err)
			}
			d.Set("destination_resource_id", destinationId)

			d.Set("enabled", pointer.From(props.Enable))
			d.Set("table_names", props.TableNames)
		}
	}
	return nil
}

func flattenDataExportDestination(input *dataexport.Destination) (string, error) {
	if input == nil {
		return "", nil
	}

	var resourceID string
	if input.ResourceId != "" {
		resourceID = input.ResourceId
		if *input.Type == dataexport.TypeEventHub {
			if input.MetaData != nil && input.MetaData.EventHubName != nil {
				eventhubName := *input.MetaData.EventHubName
				eventhubNamespaceId, err := eventhubs.ParseNamespaceIDInsensitively(resourceID)
				eventhubId := eventhubs.NewEventhubID(eventhubNamespaceId.SubscriptionId, eventhubNamespaceId.ResourceGroupName, eventhubNamespaceId.NamespaceName, eventhubName)
				if err != nil {
					return "", fmt.Errorf("parsing destination eventhub namespace ID error")
				}
				resourceID = eventhubId.ID()
			}
		}
	}

	return resourceID, nil
}
