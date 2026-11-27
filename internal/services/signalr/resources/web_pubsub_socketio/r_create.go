// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_pubsub_socketio

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2024-03-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	signalrClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (w Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := signalrClient.From(metadata.Client).WebPubSubClient.WebPubSub
			subscriptionId := account.From(metadata.Client).SubscriptionId

			var config WebPubSubSocketIOResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := webpubsub.NewWebPubSubID(subscriptionId, config.ResourceGroupName, config.Name)

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			if !features.From(metadata.Client).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return tf.ImportAsExistsError(w.ResourceType(), id.ID())
				}
			}

			expandedIdentity, err := identity.ExpandSystemOrUserAssignedMapFromModel(config.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			parameters := webpubsub.WebPubSubResource{
				Kind:     pointer.To(webpubsub.ServiceKindSocketIO),
				Location: location.Normalize(config.Location),
				Identity: expandedIdentity,
				Properties: &webpubsub.WebPubSubProperties{
					DisableAadAuth:         pointer.To(!config.AADAuthEnabled),
					DisableLocalAuth:       pointer.To(!config.LocalAuthEnabled),
					LiveTraceConfiguration: expandLiveTraceConfigFromModel(config),
					PublicNetworkAccess:    pointer.To(config.PublicNetworkAccess),
					SocketIO: &webpubsub.WebPubSubSocketIOSettings{
						ServiceMode: pointer.To(config.ServiceMode),
					},
					Tls: &webpubsub.WebPubSubTlsSettings{
						ClientCertEnabled: pointer.To(config.TlsClientCertEnabled),
					},
				},
				Sku:  expandWebPubSubSocketIOSkuFromModel(config.Sku),
				Tags: pointer.To(config.Tags),
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}

			return nil
		},
	}
}

func expandLiveTraceConfigFromModel(input WebPubSubSocketIOResourceModel) *webpubsub.LiveTraceConfiguration {
	resourceCategories := make([]webpubsub.LiveTraceCategory, 0, 3)

	enabled := pointer.To("false")
	if input.LiveTraceEnabled {
		enabled = pointer.To("true")
	}

	messageLogEnabled := "false"
	if input.LiveTraceMessagingLogsEnabled {
		messageLogEnabled = "true"
	}
	resourceCategories = append(resourceCategories, webpubsub.LiveTraceCategory{
		Name:    pointer.To("MessagingLogs"),
		Enabled: pointer.To(messageLogEnabled),
	})

	connectivityLogEnabled := "false"
	if input.LiveTraceConnectivityLogsEnabled {
		connectivityLogEnabled = "true"
	}
	resourceCategories = append(resourceCategories, webpubsub.LiveTraceCategory{
		Name:    pointer.To("ConnectivityLogs"),
		Enabled: pointer.To(connectivityLogEnabled),
	})

	httpLogEnabled := "false"
	if input.LiveTraceHttpRequestLogsEnabled {
		httpLogEnabled = "true"
	}
	resourceCategories = append(resourceCategories, webpubsub.LiveTraceCategory{
		Name:    pointer.To("HttpRequestLogs"),
		Enabled: pointer.To(httpLogEnabled),
	})

	return &webpubsub.LiveTraceConfiguration{
		Enabled:    enabled,
		Categories: &resourceCategories,
	}
}

func expandWebPubSubSocketIOSkuFromModel(input []WebPubSubSocketIOSkuModel) *webpubsub.ResourceSku {
	result := webpubsub.ResourceSku{}
	if len(input) == 0 {
		return &result
	}

	v := input[0]
	result.Name = v.Name
	result.Capacity = pointer.To(v.Capacity)

	return &result
}
