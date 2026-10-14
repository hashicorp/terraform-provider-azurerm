// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package app_service_environment_v3

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/appserviceenvironments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/preflight"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AppServiceEnvironmentName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"subnet_id": { // (@jackofallops) - This _should_ be updatable via `ChangeVnet`, but the service returns Code="NotImplemented" Message="The requested method is not implemented."
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"allow_new_private_endpoint_connections": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"cluster_setting": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"value": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"dedicated_host_count": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(2, 2), // Docs suggest is limited to 2 physical hosts at this time
			ConflictsWith: []string{
				"zone_redundant",
			},
		},

		"internal_load_balancing_mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  string(appserviceenvironments.LoadBalancingModeNone),
			ValidateFunc: validation.StringInSlice([]string{
				string(appserviceenvironments.LoadBalancingModeNone),
				string(appserviceenvironments.LoadBalancingModeWebPublishing),
			}, false),
		},

		"remote_debugging_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"zone_redundant": {
			Type:     pluginsdk.TypeBool,
			ForceNew: true,
			Optional: true,
			Default:  false,
			ConflictsWith: []string{
				"dedicated_host_count",
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"dns_suffix": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"external_inbound_ip_addresses": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"inbound_network_dependencies": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"description": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"ip_addresses": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"ports": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"internal_inbound_ip_addresses": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"ip_ssl_address_count": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"linux_outbound_ip_addresses": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"location": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"pricing_tier": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"windows_outbound_ip_addresses": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff
			if rd == nil {
				return nil
			}

			// Only run preflight validation when it's enabled and the diff contains changes (or this is a new resource).
			if metadata.Client.Features.EnhancedValidation.PreflightEnabled {
				if len(rd.GetChangedKeysPrefix("")) > 0 || rd.Id() == "" {
					var model AppServiceEnvironmentV3Model
					if err := metadata.DecodeDiff(&model); err != nil {
						return err
					}

					vnetLoc, skip := resolvePreflightVnetLocation(ctx, metadata, model)
					if skip {
						return nil
					}

					req := expandCreateForAppServiceEnvironmentV3(model, vnetLoc)
					id := commonids.NewAppServiceEnvironmentID(metadata.Client.Account.SubscriptionId, model.ResourceGroup, model.Name)

					preflightValidate, err := preflight.NewValidationRequest(pointer.To(vnetLoc), pointer.To(id), "2023-01-01", req)
					if err != nil {
						return fmt.Errorf("constructing preflight validation request: %w", err)
					}

					if err = preflightValidate.ValidateResource(ctx, metadata); err != nil {
						return err
					}
				}
			}

			return nil
		},
	}
}
