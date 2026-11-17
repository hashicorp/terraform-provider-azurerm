// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package healthcare_medtech_service

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/iotconnectors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2024-03-31/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceHealthcareApisMedTechServiceUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceIotConnectorsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspace, err := workspaces.ParseWorkspaceID(d.Get("workspace_id").(string))
	if err != nil {
		return fmt.Errorf("parsing healthcare api workspace error: %+v", err)
	}
	id := iotconnectors.NewIotConnectorID(workspace.SubscriptionId, workspace.ResourceGroupName, workspace.WorkspaceName, d.Get("name").(string))

	namespaceName := d.Get("eventhub_namespace_name").(string) + ".servicebus.windows.net"
	i, err := identity.ExpandLegacySystemAndUserAssignedMap(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	parameters := iotconnectors.IotConnector{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Identity: i,
		Properties: &iotconnectors.IotConnectorProperties{
			IngestionEndpointConfiguration: &iotconnectors.IotEventHubIngestionEndpointConfiguration{
				EventHubName:                    pointer.To(d.Get("eventhub_name").(string)),
				ConsumerGroup:                   pointer.To(d.Get("eventhub_consumer_group_name").(string)),
				FullyQualifiedEventHubNamespace: &namespaceName,
			},
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	deviceContentMap := iotconnectors.IotMappingProperties{}
	deviceMappingJson := fmt.Sprintf(`{ "content": %s }`, d.Get("device_mapping_json").(string))
	if err := json.Unmarshal([]byte(deviceMappingJson), &deviceContentMap); err != nil {
		return err
	}
	parameters.Properties.DeviceMapping = &deviceContentMap

	if err = client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceHealthcareApisMedTechServiceRead(d, meta)
}
