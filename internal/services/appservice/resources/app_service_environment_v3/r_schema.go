// Copyright IBM Corp. 2014, 2026
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
			Description:  "The name of the App Service Environment.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AppServiceEnvironmentName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"subnet_id": { // (@jackofallops) - This _should_ be updatable via `ChangeVnet`, but the service returns Code="NotImplemented" Message="The requested method is not implemented."
			Description:  "The ID of the Subnet which the App Service Environment should be connected to.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"allow_new_private_endpoint_connections": {
			Description: "Whether new Private Endpoint Connections should be allowed.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"cluster_setting": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name of the Cluster Setting.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"value": {
						Description: "The value for the Cluster Setting.",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},
				},
			},
		},

		"dedicated_host_count": {
			Description:  "This ASEv3 should use Dedicated Hosts.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(2, 2), // Docs suggest is limited to 2 physical hosts at this time
			ConflictsWith: []string{
				"zone_redundant",
			},
		},

		"internal_load_balancing_mode": {
			Description: "Which endpoints to serve internally in the Virtual Network for the App Service Environment.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true,
			Default:     string(appserviceenvironments.LoadBalancingModeNone),
			ValidateFunc: validation.StringInSlice([]string{
				string(appserviceenvironments.LoadBalancingModeNone),
				string(appserviceenvironments.LoadBalancingModeWebPublishing),
			}, false),
		},

		"remote_debugging_enabled": {
			Description: "Whether to enable remote debug.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"zone_redundant": {
			Description: "Set to `true` to deploy the ASEv3 with Availability Zones supported. Zonal ASEs can be deployed in some regions, you can refer to [Availability Zone support for App Service Environments](https://learn.microsoft.com/azure/app-service/environment/zone-redundancy). You can only set either `dedicated_host_count` or `zone_redundant` but not both.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
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
			Description: "The DNS suffix for this App Service Environment V3.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"external_inbound_ip_addresses": {
			Description: "The external inbound IP addresses of the App Service Environment V3.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
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
						Description: "A short description of the purpose of the network traffic.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"ip_addresses": {
						Description: "A list of IP addresses that network traffic will originate from in CIDR notation.",
						Type:        pluginsdk.TypeList,
						Computed:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"ports": {
						Description: "The ports that network traffic will arrive to the App Service Environment V3 on.",
						Type:        pluginsdk.TypeList,
						Computed:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"internal_inbound_ip_addresses": {
			Description: "The internal inbound IP addresses of the App Service Environment V3.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"ip_ssl_address_count": {
			Description: "The number of IP SSL addresses reserved for the App Service Environment V3.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"linux_outbound_ip_addresses": {
			Description: "Outbound addresses of Linux based Apps in this App Service Environment V3.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"location": {
			Description: "The location where the App Service Environment exists.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"pricing_tier": {
			Description: "Pricing tier for the front end instances.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"windows_outbound_ip_addresses": {
			Description: "Outbound addresses of Windows based Apps in this App Service Environment V3.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
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
