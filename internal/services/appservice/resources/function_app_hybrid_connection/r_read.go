// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package function_app_hybrid_connection

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
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
)

func (r Resource) Read() sdk.ResourceFunc {
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

func (r Resource) CustomImporter() sdk.ResourceRunFunc {
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

		if helpers.PlanIsConsumption(sku) {
			return fmt.Errorf("unsupported plan type. Hybrid Connections are not supported on Consumption service plans")
		}

		return nil
	}
}
