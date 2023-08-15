// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/hybridconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/web/2022-09-01/web"
)

type FunctionAppHybridConnectionResource struct{}

type FunctionAppHybridConnectionModel struct {
	FunctionAppId       string `tfschema:"function_app_id"`
	RelayId             string `tfschema:"relay_id"`
	HostName            string `tfschema:"hostname"`
	HostPort            int    `tfschema:"port"`
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
	return validate.AppHybridConnectionID
}

func (r FunctionAppHybridConnectionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"function_app_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FunctionAppID,
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
			appId, err := parse.FunctionAppID(appHybridConn.FunctionAppId)
			if err != nil {
				return err
			}
			relayId, err := hybridconnections.ParseHybridConnectionID(appHybridConn.RelayId)
			if err != nil {
				return err
			}

			id := parse.NewAppHybridConnectionID(appId.SubscriptionId, appId.ResourceGroup, appId.SiteName, relayId.NamespaceName, relayId.HybridConnectionName)

			existing, err := client.GetHybridConnection(ctx, id.ResourceGroup, id.SiteName, id.HybridConnectionNamespaceName, id.RelayName)
			if err != nil {
				if !utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("checking for presence of existing %s: %s", id, err)
				}
			}
			if existing.ID != nil && *existing.ID != "" {
				return tf.ImportAsExistsError(r.ResourceType(), *existing.ID)
			}

			envelope := web.HybridConnection{
				HybridConnectionProperties: &web.HybridConnectionProperties{
					RelayArmURI:  utils.String(relayId.ID()),
					Hostname:     utils.String(appHybridConn.HostName),
					Port:         utils.Int32(int32(appHybridConn.HostPort)),
					SendKeyName:  utils.String(appHybridConn.SendKeyName),
					SendKeyValue: utils.String(""),
				},
			}

			_, err = client.CreateOrUpdateHybridConnection(ctx, id.ResourceGroup, id.SiteName, id.HybridConnectionNamespaceName, id.RelayName, envelope)
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

			id, err := parse.AppHybridConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.GetHybridConnection(ctx, id.ResourceGroup, id.SiteName, id.HybridConnectionNamespaceName, id.RelayName)
			if err != nil {
				if utils.ResponseWasNotFound(existing.Response) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("reading %s: %+v", id, err)
			}

			appHybridConn := FunctionAppHybridConnectionModel{
				FunctionAppId: parse.NewFunctionAppID(id.SubscriptionId, id.ResourceGroup, id.SiteName).ID(),
				RelayName:     id.RelayName,
				NamespaceName: id.HybridConnectionNamespaceName,
			}

			if props := existing.HybridConnectionProperties; props != nil {
				appHybridConn.RelayId = utils.NormalizeNilableString(props.RelayArmURI)
				appHybridConn.HostName = utils.NormalizeNilableString(props.Hostname)
				appHybridConn.HostPort = int(utils.NormaliseNilableInt32(props.Port))
				appHybridConn.SendKeyName = utils.NormalizeNilableString(existing.SendKeyName)
				appHybridConn.ServiceBusNamespace = utils.NormalizeNilableString(props.ServiceBusNamespace)
				appHybridConn.ServiceBusSuffix = utils.NormalizeNilableString(props.ServiceBusSuffix)
				appHybridConn.SendKeyValue = utils.NormalizeNilableString(props.SendKeyValue)
			}

			if appHybridConn.ServiceBusNamespace != "" && appHybridConn.SendKeyName != "" {
				relayNamespaceClient := metadata.Client.Relay.NamespacesClient
				relayId, err := hybridconnections.ParseHybridConnectionIDInsensitively(appHybridConn.RelayId)
				if err != nil {
					return err
				}

				if keys, err := relayNamespaceClient.ListKeys(ctx, namespaces.NewAuthorizationRuleID(relayId.SubscriptionId, relayId.ResourceGroupName, appHybridConn.ServiceBusNamespace, appHybridConn.SendKeyName)); err != nil && keys.Model != nil {
					appHybridConn.SendKeyValue = utils.NormalizeNilableString(keys.Model.PrimaryKey)
					return metadata.Encode(&appHybridConn)
				}

				hybridConnectionsClient := metadata.Client.Relay.HybridConnectionsClient
				ruleID := hybridconnections.NewHybridConnectionAuthorizationRuleID(relayId.SubscriptionId, relayId.ResourceGroupName, appHybridConn.ServiceBusNamespace, *existing.Name, appHybridConn.SendKeyName)
				keys, err := hybridConnectionsClient.ListKeys(ctx, ruleID)
				if err != nil && keys.Model != nil {
					appHybridConn.SendKeyValue = utils.NormalizeNilableString(keys.Model.PrimaryKey)
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

			id, err := parse.AppHybridConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.DeleteHybridConnection(ctx, id.ResourceGroup, id.SiteName, id.HybridConnectionNamespaceName, id.RelayName)
			if err != nil {
				if !response.WasNotFound(resp.Response) {
					return fmt.Errorf("deleting %s: %+v", id, err)
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

			id, err := parse.AppHybridConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var appHybridConn FunctionAppHybridConnectionModel
			if err := metadata.Decode(&appHybridConn); err != nil {
				return err
			}

			existing, err := client.GetHybridConnection(ctx, id.ResourceGroup, id.SiteName, id.HybridConnectionNamespaceName, id.RelayName)
			if err != nil {
				if utils.ResponseWasNotFound(existing.Response) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("reading %s: %+v", id, err)
			}

			if metadata.ResourceData.HasChange("hostname") {
				existing.HybridConnectionProperties.Hostname = utils.String(appHybridConn.HostName)
			}

			if metadata.ResourceData.HasChange("port") {
				existing.HybridConnectionProperties.Port = utils.Int32(int32(appHybridConn.HostPort))
			}

			if metadata.ResourceData.HasChange("send_key_name") {
				key, err := helpers.GetSendKeyValue(ctx, metadata, *id, appHybridConn.SendKeyName)
				if err != nil {
					return err
				}
				existing.HybridConnectionProperties.SendKeyValue = key
			}

			_, err = client.CreateOrUpdateHybridConnection(ctx, id.ResourceGroup, id.SiteName, id.HybridConnectionNamespaceName, id.RelayName, existing)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r FunctionAppHybridConnectionResource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		id, err := parse.AppHybridConnectionID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}

		appId := parse.NewFunctionAppID(id.SubscriptionId, id.ResourceGroup, id.SiteName)

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
