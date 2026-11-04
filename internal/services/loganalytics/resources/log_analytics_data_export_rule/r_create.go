// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_data_export_rule

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2021-11-01/eventhubs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2022-01-01-preview/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/dataexport"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceOperationalinsightsDataExportCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).LogAnalytics.DataExportClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspace, err := dataexport.ParseWorkspaceID(d.Get("workspace_resource_id").(string))
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	id := dataexport.NewDataExportID(workspace.SubscriptionId, d.Get("resource_group_name").(string), workspace.WorkspaceName, d.Get("name").(string))
	id.Segments()

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_log_analytics_data_export_rule", id.ID())
			}
		}
	}

	destinationId := d.Get("destination_resource_id").(string)

	tableNamesGet := d.Get("table_names").(*pluginsdk.Set).List()
	tableNames := make([]string, 0, len(tableNamesGet))
	for _, v := range tableNamesGet {
		tableNames = append(tableNames, v.(string))
	}

	parameters := dataexport.DataExport{
		Properties: &dataexport.DataExportProperties{
			Destination: &dataexport.Destination{
				ResourceId: destinationId,
			},
			TableNames: tableNames,
			Enable:     pointer.To(d.Get("enabled").(bool)),
		},
	}

	if strings.Contains(destinationId, "Microsoft.EventHub") {
		if _, err := eventhubs.ValidateNamespaceID(destinationId, "destination_resource_id"); err == nil {
			eventhubNamespace, err := eventhubs.ParseNamespaceID(destinationId)
			if err != nil {
				return fmt.Errorf("parsing destination eventhub namespaces id error: %+v", err)
			}

			parameters.Properties.Destination.ResourceId = eventhubNamespace.ID()
		} else {
			eventhubId, err := eventhubs.ParseEventhubID(destinationId)
			if err != nil {
				return fmt.Errorf("parsing destination eventhub id error: %+v", err)
			}

			destinationId = namespaces.NewNamespaceID(eventhubId.SubscriptionId, eventhubId.ResourceGroupName, eventhubId.NamespaceName).ID()
			parameters.Properties.Destination.ResourceId = destinationId
			parameters.Properties.Destination.MetaData = &dataexport.DestinationMetaData{
				EventHubName: pointer.To(eventhubId.EventhubName),
			}
		}
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	// Tracked on https://github.com/Azure/azure-rest-api-specs/issues/31399
	log.Printf("[DEBUG] Waiting for Log Analytics Workspace Data Export Rule %q to become ready", id.ID())
	pollerType := custompollers.NewLogAnalyticsDataExportPoller(client, id)
	poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return err
	}

	return resourceOperationalinsightsDataExportRead(d, meta)
}
