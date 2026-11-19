// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_pubsub_socketio

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (w Resource) Arguments() map[string]*pluginsdk.Schema {
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
						ValidateFunc: validation.StringInSlice(PossibleValuesForSkuName(), false),
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
			Default:      PublicNetworkAccessEnabled,
			ValidateFunc: validation.StringInSlice(PossibleValuesForPublicNetworkAccess(), false),
		},

		"service_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      SocketIOServiceModeDefault,
			ValidateFunc: validation.StringInSlice(PossibleValuesForSocketIOServiceMode(), false),
		},

		"tags": commonschema.Tags(),

		"tls_client_cert_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},
	}
}

func (w Resource) Attributes() map[string]*pluginsdk.Schema {
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

func (w Resource) CustomizeDiff() sdk.ResourceFunc {
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
				case SkuNameFreeF1:
					if capacity > 1 {
						return fmt.Errorf("`capacity` can only be `1` when `sku` is set to `%s`", sku)
					}

					if config.PublicNetworkAccess == PublicNetworkAccessDisabled {
						return fmt.Errorf("`public_network_access` cannot be set to `Disabled` when `sku` is set to `%s`", sku)
					}

					if config.TlsClientCertEnabled {
						return fmt.Errorf("tls_client_cert_enabled` cannot be set to `true` when `sku` is set to `%s`", sku)
					}
				case SkuNameStandardS1:
					if capacity > 100 {
						return fmt.Errorf("`capacity` cannot be greater than `100` when `sku` is set to `%s`", sku)
					}
				case SkuNamePremiumP1:
					if capacity > 100 {
						return fmt.Errorf("`capacity` cannot be greater than `100` when `sku` is set to `%s`", sku)
					}
				case SkuNamePremiumP2:
					if capacity < 100 {
						return fmt.Errorf("`capacity` must be greater than or equal to `100` when `sku` is set to `%s`", sku)
					}
				}
			}

			return nil
		},
	}
}
