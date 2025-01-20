// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package healthcare

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/iotconnectors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2024-03-31/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceHealthcareIotConnector() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceHealthcareIotConnectorRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.MedTechServiceName(),
			},

			"workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: workspaces.ValidateWorkspaceID,
			},

			"identity": commonschema.SystemOrUserAssignedIdentityComputed(),

			"eventhub_namespace_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"eventhub_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"eventhub_consumer_group_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"device_mapping_json": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceHealthcareIotConnectorRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceIotConnectorsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	domainSuffix, ok := meta.(*clients.Client).Account.Environment.ServiceBus.DomainSuffix()
	if !ok {
		return fmt.Errorf("unable to retrieve the Domain Suffix for ServiceBus, this is not configured for this Cloud Environment")
	}

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspaceId, err := workspaces.ParseWorkspaceID(d.Get("workspace_id").(string))
	if err != nil {
		return fmt.Errorf("parsing workspace id error: %+v", err)
	}

	id := iotconnectors.NewIotConnectorID(subscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.IotConnectorName)

	d.Set("workspace_id", workspaceId.ID())

	if m := resp.Model; m != nil {
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
	}
	return nil
}
