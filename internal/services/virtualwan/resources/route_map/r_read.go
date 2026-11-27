// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package route_map

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.VirtualWan.VirtualWANs

			id, err := virtualwans.ParseRouteMapID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.RouteMapsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := RouteMapModel{
				Name:         id.RouteMapName,
				VirtualHubId: virtualwans.NewVirtualHubID(id.SubscriptionId, id.ResourceGroupName, id.VirtualHubName).ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.Rules = flattenRules(props.Rules)
				}
			}

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
			}

			return metadata.Encode(&state)
		},
	}
}

type Rule struct {
	Actions           []Action             `tfschema:"action"`
	MatchCriteria     []Criterion          `tfschema:"match_criterion"`
	Name              string               `tfschema:"name"`
	NextStepIfMatched virtualwans.NextStep `tfschema:"next_step_if_matched"`
}

type Action struct {
	Parameters []Parameter                    `tfschema:"parameter"`
	Type       virtualwans.RouteMapActionType `tfschema:"type"`
}

type Parameter struct {
	AsPath      []string `tfschema:"as_path"`
	Community   []string `tfschema:"community"`
	RoutePrefix []string `tfschema:"route_prefix"`
}

type Criterion struct {
	AsPath         []string                           `tfschema:"as_path"`
	Community      []string                           `tfschema:"community"`
	MatchCondition virtualwans.RouteMapMatchCondition `tfschema:"match_condition"`
	RoutePrefix    []string                           `tfschema:"route_prefix"`
}
