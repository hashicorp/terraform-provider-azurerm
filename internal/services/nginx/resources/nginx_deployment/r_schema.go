// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package nginx_deployment

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-11-01-preview/nginxdeployment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/preflight"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (m Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupName(),

		"name": {
			Description:  "The name which should be used for this NGINX Deployment.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"sku": {
			Description: "The NGINX Deployment SKU.",
			// docs: <https://docs.nginx.com/nginxaas/azure/billing/overview/>
			// we will not be forcing validation of SKU as there are internal SKUs
			// used for testing and for F5 NGINX private offers
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"location": commonschema.Location(),

		"capacity": {
			Description:   "The number of NGINX capacity units for this NGINX deployment.",
			Type:          pluginsdk.TypeInt,
			Optional:      true,
			ValidateFunc:  validation.IntPositive,
			ConflictsWith: []string{"auto_scale_profile"},
		},

		"auto_scale_profile": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			ConflictsWith: []string{"capacity"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name of the autoscaling profile.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"min_capacity": {
						Description:  "The minimum number of NGINX capacity units for this NGINX Deployment.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntPositive,
					},

					"max_capacity": {
						Description:  "The maximum number of NGINX capacity units for this NGINX Deployment.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntPositive,
					},
				},
			},
		},

		"email": {
			Description:  "The preferred support contact email address for receiving alerts and notifications.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"frontend_public": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"frontend_private"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"ip_address": {
						Description: "A list of Public IP Resource IDs for this NGINX Deployment.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
		},

		"frontend_private": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			ConflictsWith: []string{"frontend_public"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"ip_address": {
						Description: "The private IP Address.",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},

					"allocation_method": {
						Description:  "The method for allocating the private IP.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(nginxdeployment.PossibleValuesForNginxPrivateIPAllocationMethod(), false),
					},

					"subnet_id": {
						Description: "The Subnet Resource ID for this NGINX Deployment.",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},
				},
			},
		},

		"network_interface": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"subnet_id": {
						Description: "The Subnet Resource ID for this NGINX Deployment.",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},
				},
			},
		},

		"automatic_upgrade_channel": {
			Description: "The automatic upgrade channel for the NGINX deployment.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     "stable",
			ValidateFunc: validation.StringInSlice(
				[]string{
					"stable",
					"preview",
				}, false,
			),
		},

		"web_application_firewall": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"activation_state_enabled": {
						Description: "Whether WAF is enabled/disabled for this NGINX Deployment.",
						Type:        pluginsdk.TypeBool,
						Required:    true,
					},
					"status": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"attack_signatures_package": webApplicationFirewallPackageComputed(),
								"bot_signatures_package":    webApplicationFirewallPackageComputed(),
								"threat_campaigns_package":  webApplicationFirewallPackageComputed(),
								"component_versions":        webApplicationFirewallComponentVersionsComputed(),
							},
						},
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (m Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"nginx_version": {
			Description: "The version of the NGINX Deployment.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"ip_address": {
			Description: "The IP address of the NGINX Deployment.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
		"dataplane_api_endpoint": {
			Description: "The dataplane API endpoint of the NGINX Deployment.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}

func (m Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff

			if metadata.Client.Features.EnhancedValidation.PreflightEnabled {
				// Only perform preflight validation if there are changes. This avoids validation failures and
				// additional API calls for resources that are unchanged between plan invocations
				if len(rd.GetChangedKeysPrefix("")) > 0 || rd.Id() == "" {
					var model DeploymentModel
					if err := metadata.DecodeDiff(&model); err != nil {
						return err
					}

					req, err := expandCreateForNginxDeployment(model)
					if err != nil {
						return err
					}

					id := nginxdeployment.NewNginxDeploymentID(metadata.Client.Account.SubscriptionId, model.ResourceGroupName, model.Name)

					preflightValidate, err := preflight.NewValidationRequest(pointer.To(model.Location), pointer.To(id), "2024-11-01-preview", req)
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
