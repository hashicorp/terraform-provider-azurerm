// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package healthcare

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/iotconnectors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2024-03-31/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	eventhubValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceHealthcareApisMedTechService() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceHealthcareApisMedTechServiceCreate,
		Read:   resourceHealthcareApisMedTechServiceRead,
		Update: resourceHealthcareApisMedTechServiceUpdate,
		Delete: resourceHealthcareApisMedTechServiceDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(90 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(90 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(90 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.HealthCareIoTConnectorV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := iotconnectors.ParseIotConnectorID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.MedTechServiceName(),
			},

			"workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: workspaces.ValidateWorkspaceID,
			},

			"location": commonschema.Location(),

			"identity": commonschema.SystemOrUserAssignedIdentityOptional(),

			"eventhub_namespace_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: eventhubValidate.ValidateEventHubNamespaceName(),
			},

			"eventhub_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: eventhubValidate.ValidateEventHubName(),
			},

			"eventhub_consumer_group_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: eventhubValidate.ValidateEventHubConsumerName(),
			},

			"device_mapping_json": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				StateFunc:        utils.NormalizeJson,
				DiffSuppressFunc: suppressJsonOrderingDifference,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceHealthcareApisMedTechServiceCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceIotConnectorsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspace, err := workspaces.ParseWorkspaceID(d.Get("workspace_id").(string))
	if err != nil {
		return fmt.Errorf("parsing healthcare api workspace error: %+v", err)
	}
	id := iotconnectors.NewIotConnectorID(workspace.SubscriptionId, workspace.ResourceGroupName, workspace.WorkspaceName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_healthcare_medtech_service", id.ID())
		}
	}

	namespaceName := d.Get("eventhub_namespace_name").(string) + ".servicebus.windows.net"

	i, err := identity.ExpandLegacySystemAndUserAssignedMap(d.Get("identity").([]interface{}))
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
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	deviceContentMap := iotconnectors.IotMappingProperties{}
	deviceMappingJson := fmt.Sprintf(`{ "content": %s }`, d.Get("device_mapping_json").(string))
	if err := json.Unmarshal([]byte(deviceMappingJson), &deviceContentMap); err != nil {
		return err
	}
	parameters.Properties.DeviceMapping = &deviceContentMap

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceHealthcareApisMedTechServiceRead(d, meta)
}

func resourceHealthcareApisMedTechServiceRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

func resourceHealthcareApisMedTechServiceUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceIotConnectorsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspace, err := workspaces.ParseWorkspaceID(d.Get("workspace_id").(string))
	if err != nil {
		return fmt.Errorf("parsing healthcare api workspace error: %+v", err)
	}
	id := iotconnectors.NewIotConnectorID(workspace.SubscriptionId, workspace.ResourceGroupName, workspace.WorkspaceName, d.Get("name").(string))

	namespaceName := d.Get("eventhub_namespace_name").(string) + ".servicebus.windows.net"
	i, err := identity.ExpandLegacySystemAndUserAssignedMap(d.Get("identity").([]interface{}))
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
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
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

func resourceHealthcareApisMedTechServiceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceIotConnectorsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := iotconnectors.ParseIotConnectorID(d.Id())
	if err != nil {
		return err
	}

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func suppressJsonOrderingDifference(_, old, new string, _ *pluginsdk.ResourceData) bool {
	return utils.NormalizeJson(old) == utils.NormalizeJson(new)
}
