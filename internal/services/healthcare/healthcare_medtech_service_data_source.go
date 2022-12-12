package healthcare

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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
				ForceNew:     true,
				ValidateFunc: validate.WorkspaceID,
			},

			"identity": commonschema.SystemAssignedIdentityComputed(),

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
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceMedTechServiceClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspaceId, err := parse.WorkspaceID(d.Get("workspace_id").(string))
	if err != nil {
		return fmt.Errorf("parsing workspace id error: %+v", err)
	}

	id := parse.NewMedTechServiceID(subscriptionId, workspaceId.ResourceGroup, workspaceId.Name, d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.IotConnectorName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.IotConnectorName)

	d.Set("workspace_id", workspaceId.ID())

	if err := d.Set("identity", flattenMedTechServiceIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}
	if props := resp.IotConnectorProperties; props != nil {
		if props.IngestionEndpointConfiguration.EventHubName != nil {
			d.Set("eventhub_name", props.IngestionEndpointConfiguration.EventHubName)
		}

		if props.IngestionEndpointConfiguration.ConsumerGroup != nil {
			d.Set("eventhub_consumer_group_name", props.IngestionEndpointConfiguration.ConsumerGroup)
		}

		if props.DeviceMapping != nil {
			deviceMapData, err := json.Marshal(props.DeviceMapping)
			if err != nil {
				return err
			}

			var m map[string]*json.RawMessage
			if err = json.Unmarshal(deviceMapData, &m); err != nil {
				return err
			}
			mapContent := ""
			if v, ok := m["content"]; ok {
				contents, err := json.Marshal(v)
				if err != nil {
					return err
				}
				mapContent = string(contents)
			}
			d.Set("device_mapping_json", mapContent)
		}

		if props.IngestionEndpointConfiguration.FullyQualifiedEventHubNamespace != nil {
			d.Set("eventhub_namespace_name", strings.TrimSuffix(*props.IngestionEndpointConfiguration.FullyQualifiedEventHubNamespace, ".servicebus.windows.net"))
		}
	}
	return nil
}
