// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/hybridconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type WebAppSlotHybridConnectionResource struct{}

type WebAppSlotHybridConnectionModel struct {
	Name                string `tfschema:"name"`
	WebAppId            string `tfschema:"web_app_id"`
	RelayId             string `tfschema:"relay_id"`
	HostName            string `tfschema:"hostname"`
	HostPort            int64  `tfschema:"port"`
	SendKeyName         string `tfschema:"send_key_name"`
	NamespaceName       string `tfschema:"namespace_name"`
	RelayName           string `tfschema:"relay_name"`
	ServiceBusNamespace string `tfschema:"service_bus_namespace"`
	ServiceBusSuffix    string `tfschema:"service_bus_suffix"`
	SendKeyValue        string `tfschema:"send_key_value"`
}

var _ sdk.ResourceWithUpdate = WebAppSlotHybridConnectionResource{}

var _ sdk.ResourceWithCustomImporter = WebAppSlotHybridConnectionResource{}

func (r WebAppSlotHybridConnectionResource) ModelObject() interface{} {
	return &WebAppHybridConnectionModel{}
}

func (r WebAppSlotHybridConnectionResource) ResourceType() string {
	return "azurerm_web_app_slot_hybrid_connection"
}

func (r WebAppSlotHybridConnectionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return webapps.ValidateSlotHybridConnectionNamespaceRelayID
}

func (r WebAppSlotHybridConnectionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.WebAppName,
			Description:  "Specifies the name of the Web App Slot for this Hybrid Connection.",
		},
		"web_app_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateWebAppID,
			Description:  "The ID of the Web App for this Hybrid Connection.",
		},

		"relay_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: hybridconnections.ValidateHybridConnectionID,
			Description:  "The ID of the Relay Hybrid Connection to use.",
		},

		"hostname": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "The hostname of the endpoint.",
		},

		"port": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: azValidate.PortNumberOrZero,
			Description:  "The port to use for the endpoint",
		},

		"send_key_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "RootManageSharedAccessKey",
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "The name of the Relay key with `Send` permission to use. Defaults to `RootManageSharedAccessKey`",
		},
	}
}

func (r WebAppSlotHybridConnectionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"namespace_name": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The name of the Relay Namespace.",
		},

		"relay_name": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The name of the Relay in use.",
		},

		"service_bus_namespace": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The Service Bus Namespace.",
		},

		"service_bus_suffix": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The suffix for the endpoint.",
		},

		"send_key_value": {
			Type:        pluginsdk.TypeString,
			Sensitive:   true,
			Computed:    true,
			Description: "The Primary Access Key for the `send_key_name`",
		},
	}
}

func (r WebAppSlotHybridConnectionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var appSlotHybridConn WebAppSlotHybridConnectionModel

			client := metadata.Client.AppService.WebAppsClient

			if err := metadata.Decode(&appSlotHybridConn); err != nil {
				return err
			}
			appId, err := commonids.ParseWebAppID(appSlotHybridConn.WebAppId)
			if err != nil {
				return err
			}
			relayId, err := hybridconnections.ParseHybridConnectionID(appSlotHybridConn.RelayId)
			if err != nil {
				return err
			}

			id := webapps.NewSlotHybridConnectionNamespaceRelayID(appId.SubscriptionId, appId.ResourceGroupName, appId.SiteName, appSlotHybridConn.Name, relayId.NamespaceName, relayId.HybridConnectionName)

			existing, err := client.GetHybridConnectionSlot(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %s", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError(r.ResourceType(), id.ID())
			}

			sendKeyValue, err := helpers.GetSendKeyValue(ctx, metadata, *relayId, appSlotHybridConn.SendKeyName)
			if err != nil {
				return err
			}

			envelope := webapps.HybridConnection{
				Properties: &webapps.HybridConnectionProperties{
					RelayArmUri:  pointer.To(relayId.ID()),
					Hostname:     pointer.To(appSlotHybridConn.HostName),
					Port:         pointer.To(appSlotHybridConn.HostPort),
					SendKeyName:  pointer.To(appSlotHybridConn.SendKeyName),
					SendKeyValue: sendKeyValue,
				},
			}

			_, err = client.CreateOrUpdateHybridConnectionSlot(ctx, id, envelope)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r WebAppSlotHybridConnectionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient

			id, err := webapps.ParseSlotHybridConnectionNamespaceRelayID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.GetHybridConnectionSlot(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("reading %s: %+v", id, err)
			}

			appSlotHybridConn := WebAppSlotHybridConnectionModel{
				Name:          id.SlotName,
				WebAppId:      commonids.NewAppServiceID(id.SubscriptionId, id.ResourceGroupName, id.SiteName).ID(),
				RelayName:     id.RelayName,
				NamespaceName: id.HybridConnectionNamespaceName,
			}

			if model := existing.Model; model != nil {
				if props := model.Properties; props != nil {
					appSlotHybridConn.RelayId = pointer.From(props.RelayArmUri)
					appSlotHybridConn.HostName = pointer.From(props.Hostname)
					appSlotHybridConn.HostPort = pointer.From(props.Port)
					appSlotHybridConn.SendKeyName = pointer.From(props.SendKeyName)
					appSlotHybridConn.ServiceBusNamespace = pointer.From(props.ServiceBusNamespace)
					appSlotHybridConn.ServiceBusSuffix = pointer.From(props.ServiceBusSuffix)
					appSlotHybridConn.SendKeyValue = pointer.From(props.SendKeyValue)
				}

				if appSlotHybridConn.ServiceBusNamespace != "" && appSlotHybridConn.SendKeyName != "" {
					relayNamespaceClient := metadata.Client.Relay.NamespacesClient
					relayId, err := hybridconnections.ParseHybridConnectionIDInsensitively(appSlotHybridConn.RelayId)
					if err != nil {
						return err
					}

					if keys, err := relayNamespaceClient.ListKeys(ctx, namespaces.NewAuthorizationRuleID(relayId.SubscriptionId, relayId.ResourceGroupName, appSlotHybridConn.ServiceBusNamespace, appSlotHybridConn.SendKeyName)); err != nil && keys.Model != nil {
						appSlotHybridConn.SendKeyValue = pointer.From(keys.Model.PrimaryKey)
						return metadata.Encode(&appSlotHybridConn)
					}

					hybridConnectionsClient := metadata.Client.Relay.HybridConnectionsClient
					ruleID := hybridconnections.NewHybridConnectionAuthorizationRuleID(relayId.SubscriptionId, relayId.ResourceGroupName, appSlotHybridConn.ServiceBusNamespace, pointer.From(model.Name), appSlotHybridConn.SendKeyName)
					keys, err := hybridConnectionsClient.ListKeys(ctx, ruleID)
					if err != nil && keys.Model != nil {
						appSlotHybridConn.SendKeyValue = pointer.From(keys.Model.PrimaryKey)
					}
				}
			}

			return metadata.Encode(&appSlotHybridConn)
		},
	}
}

func (r WebAppSlotHybridConnectionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient

			id, err := webapps.ParseSlotHybridConnectionNamespaceRelayID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.DeleteHybridConnectionSlot(ctx, *id)
			if err != nil {
				if !response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("deleting %s: %+v", id, err)
				}
			}

			return nil
		},
	}
}

func (r WebAppSlotHybridConnectionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient

			id, err := webapps.ParseSlotHybridConnectionNamespaceRelayID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var appHybridConn WebAppSlotHybridConnectionModel
			if err := metadata.Decode(&appHybridConn); err != nil {
				return err
			}

			existing, err := client.GetHybridConnectionSlot(ctx, *id)
			if err != nil || existing.Model == nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("reading %s: %+v", id, err)
			}

			model := *existing.Model

			if metadata.ResourceData.HasChange("hostname") {
				model.Properties.Hostname = pointer.To(appHybridConn.HostName)
			}

			if metadata.ResourceData.HasChange("port") {
				model.Properties.Port = pointer.To(appHybridConn.HostPort)
			}

			if metadata.ResourceData.HasChange("send_key_name") {
				relayId, err := hybridconnections.ParseHybridConnectionID(appHybridConn.RelayId)
				if err != nil {
					return err
				}

				sendKeyValue, err := helpers.GetSendKeyValue(ctx, metadata, *relayId, appHybridConn.SendKeyName)
				if err != nil {
					return err
				}
				model.Properties.SendKeyValue = sendKeyValue
			}

			_, err = client.CreateOrUpdateHybridConnectionSlot(ctx, *id, model)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r WebAppSlotHybridConnectionResource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		id, err := webapps.ParseSlotHybridConnectionNamespaceRelayID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}
		appId := commonids.NewAppServiceID(id.SubscriptionId, id.ResourceGroupName, id.SiteName)

		_, sku, err := helpers.ServicePlanInfoForApp(ctx, metadata, appId)
		if err != nil {
			return err
		}

		if helpers.PlanIsConsumption(sku) || helpers.PlanIsElastic(sku) {
			return fmt.Errorf("unsupported plan type. Hybrid Connections are not supported on Consumption or Elastic service plans")
		}

		return nil
	}
}
