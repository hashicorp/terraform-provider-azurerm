// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package network_manager_routing_rule

import (
	"context"
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/routingrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9\_\.\-]{1,64}$`),
				"`name` must be between 1 and 64 characters long and can only contain letters, numbers, underscores(_), periods(.), and hyphens(-).",
			),
		},

		"rule_collection_id": commonschema.ResourceIDReferenceRequiredForceNew(&routingrules.RuleCollectionId{}),

		"destination": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"address": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"type": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(routingrules.PossibleValuesForRoutingRuleDestinationType(), false),
					},
				},
			},
		},

		"next_hop": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(routingrules.PossibleValuesForRoutingRuleNextHopType(), false),
					},
					"address": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsIPAddress,
					},
				},
			},
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ManagerRoutingRuleResourceModel
			if err := metadata.DecodeDiff(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if len(model.Destination) > 0 {
				v := model.Destination[0]
				if strings.EqualFold(v.Type, string(routingrules.RoutingRuleDestinationTypeAddressPrefix)) {
					if _, _, err := net.ParseCIDR(v.Address); err != nil {
						return fmt.Errorf("expanding `destination`: `address` must be a valid CIDR when `type` is `AddressPrefix`: %+v", err)
					}
				}
			}

			if len(model.NextHop) > 0 {
				v := model.NextHop[0]
				if strings.EqualFold(v.Type, string(routingrules.RoutingRuleNextHopTypeVirtualAppliance)) && v.Address == "" {
					return fmt.Errorf("expanding `next_hop`: `address` is required when `type` is `VirtualAppliance`")
				}
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
