// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loganalytics

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2021-11-01/eventhubs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2022-01-01-preview/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/dataexport"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLogAnalyticsDataExport() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceOperationalinsightsDataExportCreateUpdate,
		Read:   resourceOperationalinsightsDataExportRead,
		Update: resourceOperationalinsightsDataExportCreateUpdate,
		Delete: resourceOperationalinsightsDataExportDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := dataexport.ParseDataExportID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.DataExportV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     validate.LogAnalyticsDataExportName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"workspace_resource_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: dataexport.ValidateWorkspaceID,
			},

			"destination_resource_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"table_names": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
			},

			"export_rule_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOperationalinsightsDataExportCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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

	// Tracked on https://github.com/Azure/azure-rest-api-specs/issues/31399
	log.Printf("[DEBUG] Waiting for Log Analytics Workspace Data Export Rule %q to become ready", id.ID())
	pollerType := custompollers.NewLogAnalyticsDataExportPoller(client, id)
	poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return err
	}

	d.SetId(id.ID())
	return resourceOperationalinsightsDataExportRead(d, meta)
}

func resourceOperationalinsightsDataExportRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
			dataExportId := ""
			if props.DataExportId != nil {
				dataExportId = *props.DataExportId
			}
			d.Set("export_rule_id", dataExportId)

			destinationId, err := flattenDataExportDestination(props.Destination)
			if err != nil {
				return fmt.Errorf("flattening destination ID error: %+v", err)
			}
			d.Set("destination_resource_id", destinationId)

			enabled := false
			if props.Enable != nil {
				enabled = *props.Enable
			}
			d.Set("enabled", enabled)
			d.Set("table_names", props.TableNames)
		}
	}
	return nil
}

func resourceOperationalinsightsDataExportDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.DataExportClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dataexport.ParseDataExportID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
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
