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
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type FunctionAppHybridConnectionResource struct{}

type FunctionAppHybridConnectionModel struct {
	FunctionAppId       string `tfschema:"function_app_id"`
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

var _ sdk.ResourceWithUpdate = FunctionAppHybridConnectionResource{}

var _ sdk.ResourceWithCustomImporter = FunctionAppHybridConnectionResource{}

func (r FunctionAppHybridConnectionResource) ModelObject() interface{} {
	return &FunctionAppHybridConnectionModel{}
}

func (r FunctionAppHybridConnectionResource) ResourceType() string {
	return "azurerm_function_app_hybrid_connection"
}

func (r FunctionAppHybridConnectionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return webapps.ValidateRelayID
}

func (r FunctionAppHybridConnectionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"function_app_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateFunctionAppID,
			Description:  "The ID of the Function App for this Hybrid Connection.",
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

func (r FunctionAppHybridConnectionResource) Attributes() map[string]*pluginsdk.Schema {
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

func (r FunctionAppHybridConnectionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var appHybridConn FunctionAppHybridConnectionModel

			client := metadata.Client.AppService.WebAppsClient

			if err := metadata.Decode(&appHybridConn); err != nil {
				return err
			}
			appId, err := commonids.ParseFunctionAppID(appHybridConn.FunctionAppId)
			if err != nil {
				return err
			}
			relayId, err := hybridconnections.ParseHybridConnectionID(appHybridConn.RelayId)
			if err != nil {
				return err
			}

			id := webapps.NewRelayID(appId.SubscriptionId, appId.ResourceGroupName, appId.SiteName, relayId.NamespaceName, relayId.HybridConnectionName)

			existing, err := client.GetHybridConnection(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %s", id, err)
				}
			}
			if existing.Model != nil && existing.Model.Id != nil && *existing.Model.Id != "" {
				return tf.ImportAsExistsError(r.ResourceType(), id.ID())
			}

			sendKeyValue, err := helpers.GetSendKeyValue(ctx, metadata, *relayId, appHybridConn.SendKeyName)
			if err != nil {
				return err
			}

			envelope := webapps.HybridConnection{
				Properties: &webapps.HybridConnectionProperties{
					RelayArmUri:  pointer.To(relayId.ID()),
					Hostname:     pointer.To(appHybridConn.HostName),
					Port:         pointer.To(appHybridConn.HostPort),
					SendKeyName:  pointer.To(appHybridConn.SendKeyName),
					SendKeyValue: sendKeyValue,
				},
			}

			_, err = client.CreateOrUpdateHybridConnection(ctx, id, envelope)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r FunctionAppHybridConnectionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient

			id, err := webapps.ParseRelayID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.GetHybridConnection(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("reading %s: %+v", id, err)
			}

			appHybridConn := FunctionAppHybridConnectionModel{
				FunctionAppId: commonids.NewAppServiceID(id.SubscriptionId, id.ResourceGroupName, id.SiteName).ID(),
				RelayName:     id.RelayName,
				NamespaceName: id.HybridConnectionNamespaceName,
			}

			if model := existing.Model; model != nil {

				if props := model.Properties; props != nil {
					appHybridConn.RelayId = pointer.From(props.RelayArmUri)
					appHybridConn.HostName = pointer.From(props.Hostname)
					appHybridConn.HostPort = pointer.From(props.Port)
					appHybridConn.SendKeyName = pointer.From(props.SendKeyName)
					appHybridConn.ServiceBusNamespace = pointer.From(props.ServiceBusNamespace)
					appHybridConn.ServiceBusSuffix = pointer.From(props.ServiceBusSuffix)
					appHybridConn.SendKeyValue = pointer.From(props.SendKeyValue)
				}

				if appHybridConn.ServiceBusNamespace != "" && appHybridConn.SendKeyName != "" {
					relayNamespaceClient := metadata.Client.Relay.NamespacesClient
					relayId, err := hybridconnections.ParseHybridConnectionIDInsensitively(appHybridConn.RelayId)
					if err != nil {
						return err
					}

					if keys, err := relayNamespaceClient.ListKeys(ctx, namespaces.NewAuthorizationRuleID(relayId.SubscriptionId, relayId.ResourceGroupName, appHybridConn.ServiceBusNamespace, appHybridConn.SendKeyName)); err != nil && keys.Model != nil {
						appHybridConn.SendKeyValue = pointer.From(keys.Model.PrimaryKey)
						return metadata.Encode(&appHybridConn)
					}

					hybridConnectionsClient := metadata.Client.Relay.HybridConnectionsClient
					ruleID := hybridconnections.NewHybridConnectionAuthorizationRuleID(relayId.SubscriptionId, relayId.ResourceGroupName, appHybridConn.ServiceBusNamespace, *model.Name, appHybridConn.SendKeyName)
					keys, err := hybridConnectionsClient.ListKeys(ctx, ruleID)
					if err != nil && keys.Model != nil {
						appHybridConn.SendKeyValue = pointer.From(keys.Model.PrimaryKey)
					}
				}
			}

			return metadata.Encode(&appHybridConn)
		},
	}
}

func (r FunctionAppHybridConnectionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient

			id, err := webapps.ParseRelayID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.DeleteHybridConnection(ctx, *id)
			if err != nil {
				if !response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("deleting %s: %+v", *id, err)
				}
			}

			return nil
		},
	}
}

func (r FunctionAppHybridConnectionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient

			id, err := webapps.ParseRelayID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var appHybridConn FunctionAppHybridConnectionModel
			if err := metadata.Decode(&appHybridConn); err != nil {
				return err
			}

			existing, err := client.GetHybridConnection(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("reading %s: %+v", *id, err)
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

				key, err := helpers.GetSendKeyValue(ctx, metadata, *relayId, appHybridConn.SendKeyName)
				if err != nil {
					return err
				}
				model.Properties.SendKeyValue = key
			}

			_, err = client.CreateOrUpdateHybridConnection(ctx, *id, model)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r FunctionAppHybridConnectionResource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		id, err := webapps.ParseRelayID(metadata.ResourceData.Id())
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
