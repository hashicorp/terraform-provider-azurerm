// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventgrid_namespace

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2023-12-15-preview/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/topics"
	"github.com/hashicorp/terraform-provider-azurerm/internal/preflight"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(3, 50),
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"capacity": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 40),
			Default:      1,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"inbound_ip_rule": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 128,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"ip_mask": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsCIDRIPv4,
					},
					"action": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(namespaces.IPActionTypeAllow),
						ValidateFunc: validation.StringInSlice(namespaces.PossibleValuesForIPActionType(), false),
					},
				},
			},
		},

		"public_network_access": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(namespaces.PublicNetworkAccessEnabled),
				string(namespaces.PublicNetworkAccessDisabled),
			}, false),
			Default: string(namespaces.PublicNetworkAccessEnabled),
		},

		"sku": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(namespaces.PossibleValuesForSkuName(), false),
			Default:      namespaces.SkuNameStandard,
		},

		"topic_spaces_configuration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"alternative_authentication_name_source": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(namespaces.PossibleValuesForAlternativeAuthenticationNameSource(), false),
						},
					},

					"maximum_client_sessions_per_authentication_name": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      1,
						ValidateFunc: validation.IntBetween(1, 100),
					},

					"maximum_session_expiry_in_hours": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      1,
						ValidateFunc: validation.IntBetween(1, 8),
					},

					"route_topic_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: topics.ValidateTopicID,
					},

					"dynamic_routing_enrichment": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"key": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringLenBetween(1, 20),
								},

								"value": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringLenBetween(1, 128),
								},
							},
						},
					},

					"static_routing_enrichment": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"key": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringLenBetween(1, 20),
								},

								"value": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringLenBetween(1, 128),
								},
							},
						},
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if metadata.ResourceDiff == nil {
				return nil
			}

			if metadata.Client.Features.EnhancedValidation.PreflightEnabled {
				if len(metadata.ResourceDiff.GetChangedKeysPrefix("")) > 0 || metadata.ResourceDiff.Id() == "" {
					var model EventGridNamespaceResourceModel
					if err := metadata.DecodeDiff(&model); err != nil {
						return err
					}

					req, err := expandCreateForEventGridNamespace(model)
					if err != nil {
						return err
					}

					id := namespaces.NewNamespaceID(metadata.Client.Account.SubscriptionId, model.ResourceGroup, model.Name)
					preflightValidate, err := preflight.NewValidationRequest(pointer.To(model.Location), pointer.To(id), "2023-12-15-preview", req)
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
