// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package signalr

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2024-03-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type WebPubSubSocketIOResource struct{}

type WebPubSubSocketIOResourceModel struct {
	Name                             string                                     `tfschema:"name"`
	ResourceGroupName                string                                     `tfschema:"resource_group_name"`
	Location                         string                                     `tfschema:"location"`
	Sku                              []WebPubSubSocketIOSkuModel                `tfschema:"sku"`
	AADAuthEnabled                   bool                                       `tfschema:"aad_auth_enabled"`
	Identity                         []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	LiveTraceEnabled                 bool                                       `tfschema:"live_trace_enabled"`
	LiveTraceConnectivityLogsEnabled bool                                       `tfschema:"live_trace_connectivity_logs_enabled"`
	LiveTraceHttpRequestLogsEnabled  bool                                       `tfschema:"live_trace_http_request_logs_enabled"`
	LiveTraceMessagingLogsEnabled    bool                                       `tfschema:"live_trace_messaging_logs_enabled"`
	LocalAuthEnabled                 bool                                       `tfschema:"local_auth_enabled"`
	PublicNetworkAccess              string                                     `tfschema:"public_network_access"`
	ServiceMode                      string                                     `tfschema:"service_mode"`
	Tags                             map[string]string                          `tfschema:"tags"`
	TlsClientCertEnabled             bool                                       `tfschema:"tls_client_cert_enabled"`
	ExternalIP                       string                                     `tfschema:"external_ip"`
	HostName                         string                                     `tfschema:"hostname"`
	PrimaryAccessKey                 string                                     `tfschema:"primary_access_key"`
	PrimaryConnectionString          string                                     `tfschema:"primary_connection_string"`
	PublicPort                       int64                                      `tfschema:"public_port"`
	ServerPort                       int64                                      `tfschema:"server_port"`
	SecondaryAccessKey               string                                     `tfschema:"secondary_access_key"`
	SecondaryConnectionString        string                                     `tfschema:"secondary_connection_string"`
}

type WebPubSubSocketIOSkuModel struct {
	Name     string `tfschema:"name"`
	Capacity int64  `tfschema:"capacity"`
}

var (
	_ sdk.ResourceWithUpdate        = WebPubSubSocketIOResource{}
	_ sdk.ResourceWithCustomizeDiff = WebPubSubSocketIOResource{}
)

func (w WebPubSubSocketIOResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.WebPubSubName(),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"sku": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(helpers.PossibleValuesForSkuName(), false),
					},
					"capacity": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  1,
						ValidateFunc: validation.IntInSlice([]int{
							1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 200,
							300, 400, 500, 600, 700, 800, 900, 1000,
						}),
					},
				},
			},
		},

		"aad_auth_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"identity": commonschema.SystemOrUserAssignedIdentityOptional(),

		"live_trace_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"live_trace_connectivity_logs_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"live_trace_http_request_logs_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"live_trace_messaging_logs_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"local_auth_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"public_network_access": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      helpers.PublicNetworkAccessEnabled,
			ValidateFunc: validation.StringInSlice(helpers.PossibleValuesForPublicNetworkAccess(), false),
		},

		"service_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      helpers.SocketIOServiceModeDefault,
			ValidateFunc: validation.StringInSlice(helpers.PossibleValuesForSocketIOServiceMode(), false),
		},

		"tags": commonschema.Tags(),

		"tls_client_cert_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},
	}
}

func (w WebPubSubSocketIOResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"external_ip": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"hostname": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"primary_access_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"primary_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"public_port": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"server_port": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"secondary_access_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"secondary_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},
	}
}

func (w WebPubSubSocketIOResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var config WebPubSubSocketIOResourceModel
			if err := metadata.DecodeDiff(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if len(config.Sku) > 0 {
				capacity := config.Sku[0].Capacity
				sku := config.Sku[0].Name

				switch sku {
				case helpers.SkuNameFreeF1:
					if capacity > 1 {
						return fmt.Errorf("`capacity` can only be `1` when `sku` is set to `%s`", sku)
					}

					if config.PublicNetworkAccess == helpers.PublicNetworkAccessDisabled {
						return fmt.Errorf("`public_network_access` cannot be set to `Disabled` when `sku` is set to `%s`", sku)
					}

					if config.TlsClientCertEnabled {
						return fmt.Errorf("tls_client_cert_enabled` cannot be set to `true` when `sku` is set to `%s`", sku)
					}
				case helpers.SkuNameStandardS1:
					if capacity > 100 {
						return fmt.Errorf("`capacity` cannot be greater than `100` when `sku` is set to `%s`", sku)
					}
				case helpers.SkuNamePremiumP1:
					if capacity > 100 {
						return fmt.Errorf("`capacity` cannot be greater than `100` when `sku` is set to `%s`", sku)
					}
				case helpers.SkuNamePremiumP2:
					if capacity < 100 {
						return fmt.Errorf("`capacity` must be greater than or equal to `100` when `sku` is set to `%s`", sku)
					}
				}
			}

			return nil
		},
	}
}

func (w WebPubSubSocketIOResource) ModelObject() interface{} {
	return &WebPubSubSocketIOResourceModel{}
}

func (w WebPubSubSocketIOResource) ResourceType() string {
	return "azurerm_web_pubsub_socketio"
}

func (w WebPubSubSocketIOResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SignalR.WebPubSubClient.WebPubSub
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config WebPubSubSocketIOResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := webpubsub.NewWebPubSubID(subscriptionId, config.ResourceGroupName, config.Name)

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_web_pubsub_socketio", id.ID())
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

			if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (w WebPubSubSocketIOResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SignalR.WebPubSubClient.WebPubSub

			id, err := webpubsub.ParseWebPubSubID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			keys, err := client.ListKeys(ctx, *id)
			if err != nil {
				return fmt.Errorf("listing keys for %s: %+v", *id, err)
			}

			state := WebPubSubSocketIOResourceModel{
				Name:              id.WebPubSubName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)

				flattenedIdentity, err := identity.FlattenSystemOrUserAssignedMapToModel(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}
				state.Identity = pointer.From(flattenedIdentity)

				if sku := model.Sku; sku != nil {
					state.Sku = flattenWebPubSubSocketIOSkuToModel(sku)
				}

				if props := model.Properties; props != nil {
					liveTrace := flattenLiveTraceConfigToMap(props.LiveTraceConfiguration)

					state.AADAuthEnabled = !pointer.From(props.DisableAadAuth)
					state.ExternalIP = pointer.From(props.ExternalIP)
					state.HostName = pointer.From(props.HostName)
					state.PublicPort = pointer.From(props.PublicPort)
					state.ServerPort = pointer.From(props.ServerPort)

					state.LiveTraceEnabled = liveTrace["enabled"]
					state.LiveTraceConnectivityLogsEnabled = liveTrace["connectivityLogsEnabled"]
					state.LiveTraceHttpRequestLogsEnabled = liveTrace["httpLogsEnabled"]
					state.LiveTraceMessagingLogsEnabled = liveTrace["messagingLogsEnabled"]

					state.LocalAuthEnabled = !pointer.From(props.DisableLocalAuth)
					state.PublicNetworkAccess = pointer.From(props.PublicNetworkAccess)

					if socketio := props.SocketIO; socketio != nil {
						state.ServiceMode = pointer.From(socketio.ServiceMode)
					}

					if tls := props.Tls; tls != nil {
						state.TlsClientCertEnabled = pointer.From(tls.ClientCertEnabled)
					}
				}
			}

			if keyModel := keys.Model; keyModel != nil {
				state.PrimaryAccessKey = pointer.From(keyModel.PrimaryKey)
				state.PrimaryConnectionString = pointer.From(keyModel.PrimaryConnectionString)
				state.SecondaryAccessKey = pointer.From(keyModel.SecondaryKey)
				state.SecondaryConnectionString = pointer.From(keyModel.SecondaryConnectionString)
			}

			return metadata.Encode(&state)
		},
	}
}

func (w WebPubSubSocketIOResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SignalR.WebPubSubClient.WebPubSub

			id, err := webpubsub.ParseWebPubSubID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config WebPubSubSocketIOResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			if existing.Model.Sku == nil {
				return fmt.Errorf("retrieving %s: `model.Sku` was nil", *id)
			}

			payload := webpubsub.WebPubSubResource{
				Location:   location.Normalize(config.Location),
				Properties: &webpubsub.WebPubSubProperties{},
				Sku:        existing.Model.Sku,
			}
			props := payload.Properties

			rd := metadata.ResourceData
			if rd.HasChange("sku") {
				payload.Sku = expandWebPubSubSocketIOSkuFromModel(config.Sku)
			}

			if rd.HasChange("aad_auth_enabled") {
				props.DisableAadAuth = pointer.To(!config.AADAuthEnabled)
			}

			if rd.HasChange("identity") {
				expandedIdentity, err := identity.ExpandSystemOrUserAssignedMapFromModel(config.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				payload.Identity = expandedIdentity
			}

			if rd.HasChanges("live_trace_enabled", "live_trace_connectivity_logs_enabled", "live_trace_http_request_logs_enabled", "live_trace_messaging_logs_enabled") {
				props.LiveTraceConfiguration = expandLiveTraceConfigFromModel(config)
			}

			if rd.HasChange("local_auth_enabled") {
				props.DisableLocalAuth = pointer.To(!config.LocalAuthEnabled)
			}

			if rd.HasChange("public_network_access") {
				props.PublicNetworkAccess = pointer.To(config.PublicNetworkAccess)
			}

			if rd.HasChange("service_mode") {
				props.SocketIO = &webpubsub.WebPubSubSocketIOSettings{
					ServiceMode: pointer.To(config.ServiceMode),
				}
			}

			if rd.HasChange("tags") {
				payload.Tags = pointer.To(config.Tags)
			}

			if rd.HasChange("tls_client_cert_enabled") {
				props.Tls = &webpubsub.WebPubSubTlsSettings{
					ClientCertEnabled: pointer.To(config.TlsClientCertEnabled),
				}
			}

			if err := client.UpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (w WebPubSubSocketIOResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SignalR.WebPubSubClient.WebPubSub

			id, err := webpubsub.ParseWebPubSubID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (w WebPubSubSocketIOResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return webpubsub.ValidateWebPubSubID
}

func expandLiveTraceConfigFromModel(input WebPubSubSocketIOResourceModel) *webpubsub.LiveTraceConfiguration {
	resourceCategories := make([]webpubsub.LiveTraceCategory, 0)

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

func flattenLiveTraceConfigToMap(input *webpubsub.LiveTraceConfiguration) map[string]bool {
	result := make(map[string]bool)
	if input == nil {
		return result
	}

	result["enabled"] = strings.EqualFold(pointer.From(input.Enabled), "true")

	if input.Categories != nil {
		for _, item := range *input.Categories {
			name := pointer.From(item.Name)
			categoryEnabled := pointer.From(item.Enabled)

			switch name {
			case "ConnectivityLogs":
				result["connectivityLogsEnabled"] = strings.EqualFold(categoryEnabled, "true")
			case "HttpRequestLogs":
				result["httpLogsEnabled"] = strings.EqualFold(categoryEnabled, "true")
			case "MessagingLogs":
				result["messagingLogsEnabled"] = strings.EqualFold(categoryEnabled, "true")
			default:
				continue
			}
		}
	}

	return result
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

func flattenWebPubSubSocketIOSkuToModel(input *webpubsub.ResourceSku) []WebPubSubSocketIOSkuModel {
	result := make([]WebPubSubSocketIOSkuModel, 0)
	if input == nil {
		return result
	}

	result = append(result, WebPubSubSocketIOSkuModel{
		Name:     input.Name,
		Capacity: pointer.From(input.Capacity),
	})

	return result
}
