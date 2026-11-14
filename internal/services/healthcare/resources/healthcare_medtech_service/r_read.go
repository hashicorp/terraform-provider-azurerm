// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package healthcare_medtech_service

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/iotconnectors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2024-03-31/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceHealthcareApisMedTechServiceRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceIotConnectorsClient
	domainSuffix, ok := meta.(*clients.Client).Account.Environment.ServiceBus.DomainSuffix()
	if !ok {
		return fmt.Errorf("unable to retrieve the Domain Suffix for ServiceBus, this is not configured for this Cloud Environment")
	}
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := iotconnectors.ParseIotConnectorID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[WARN] Healthcare Apis MedTech Service %s was not found", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.IotConnectorName)
	d.Set("workspace_id", workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName).ID())

	if m := resp.Model; m != nil {
		d.Set("location", location.NormalizeNilable(m.Location))

		i, err := identity.FlattenLegacySystemAndUserAssignedMap(m.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", i); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if props := m.Properties; props != nil {
			eventHubNamespaceName := ""
			if config := props.IngestionEndpointConfiguration; config != nil {
				d.Set("eventhub_consumer_group_name", pointer.From(config.ConsumerGroup))
				d.Set("eventhub_name", pointer.From(config.EventHubName))

				if props.IngestionEndpointConfiguration.FullyQualifiedEventHubNamespace != nil {
					suffixToTrim := "." + *domainSuffix
					eventHubNamespaceName = strings.TrimSuffix(*props.IngestionEndpointConfiguration.FullyQualifiedEventHubNamespace, suffixToTrim)
				}
			}

			d.Set("eventhub_namespace_name", eventHubNamespaceName)

			mapContent := ""
			if props.DeviceMapping != nil {
				deviceMapData, err := json.Marshal(props.DeviceMapping)
				if err != nil {
					return err
				}

				var m map[string]*json.RawMessage
				if err = json.Unmarshal(deviceMapData, &m); err != nil {
					return err
				}
				if v, ok := m["content"]; ok {
					contents, err := json.Marshal(v)
					if err != nil {
						return err
					}
					mapContent = string(contents)
				}
			}
			d.Set("device_mapping_json", mapContent)
		}

		if err := tags.FlattenAndSet(d, m.Tags); err != nil {
			return err
		}
	}
	return nil
}
