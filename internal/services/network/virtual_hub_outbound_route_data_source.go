package network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/virtualwans"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type VirtualHubOutboundRouteDataSource struct{}

type VirtualHubOutboundRouteModel struct {
	VirtualHubId   string                                 `tfschema:"virtual_hub_id"`
	ResourceUri    string                                 `tfschema:"target_resource_id"`
	ConnectionType string                                 `tfschema:"connection_type"`
	RouteMap       []VirtualHubOutboundRouteRouteMapModel `tfschema:"route_map"`
}

type VirtualHubOutboundRouteRouteMapModel struct {
	Prefix         string `tfschema:"prefix"`
	BgpCommunities string `tfschema:"bgp_communities"`
	AsPath         string `tfschema:"as_path"`
}

var _ sdk.DataSource = VirtualHubOutboundRouteDataSource{}

func (v VirtualHubOutboundRouteDataSource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"virtual_hub_id": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validate.VirtualHubID,
		},

		"target_resource_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: azure.ValidateResourceID,
		},

		"connection_type": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (v VirtualHubOutboundRouteDataSource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"route_map": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"prefix": {
						Type:     schema.TypeString,
						Computed: true,
					},

					"bgp_communities": {
						Type:     schema.TypeString,
						Computed: true,
					},

					"as_path": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}

func (v VirtualHubOutboundRouteDataSource) ModelObject() interface{} {
	return &VirtualHubOutboundRouteModel{}
}

func (v VirtualHubOutboundRouteDataSource) ResourceType() string {
	return "azurerm_virtual_hub_outbound_route"
}

func (v VirtualHubOutboundRouteDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model VirtualHubOutboundRouteModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			client := metadata.Client.Network.VirtualWANs
			virtualHubId, err := virtualwans.ParseVirtualHubID(model.VirtualHubId)
			if err != nil {
				return fmt.Errorf("parsing %q: %+v", model.VirtualHubId, err)
			}

			param := virtualwans.GetOutboundRoutesParameters{}

			if model.ResourceUri != "" {
				param.ResourceUri = &model.ResourceUri
			}

			if model.ConnectionType != "" {
				param.ConnectionType = &model.ConnectionType
			}

			resp, err := client.VirtualHubsGetOutboundRoutes(ctx, *virtualHubId, param)
			if err != nil {
				return fmt.Errorf("retrieving Virtual Hub Inbound Route: %+v", err)
			}

			err = resp.Poller.PollUntilDone(ctx)
			if err != nil {
				return fmt.Errorf("polling Virtual Hub Inbound Route: %+v", err)
			}

			response := resp.Poller.LatestResponse()
			if response == nil {
				return fmt.Errorf("no response from Virtual Hub Inbound Route")
			}

			var respModel VirtualHubsOutboundRouteResponse
			if err := response.Unmarshal(&respModel); err != nil {
				return fmt.Errorf("unmarshalling LRO response for %s: %+v", virtualHubId, err)
			}

			routeMap := make([]VirtualHubOutboundRouteRouteMapModel, 0)
			for _, r := range respModel.Properties.Output.Value {
				routeMap = append(routeMap, VirtualHubOutboundRouteRouteMapModel{
					Prefix:         r.Prefix,
					BgpCommunities: r.BgpCommunities,
					AsPath:         r.AsPath,
				})
			}

			state := VirtualHubOutboundRouteModel{
				VirtualHubId:   virtualHubId.ID(),
				ResourceUri:    model.ResourceUri,
				ConnectionType: model.ConnectionType,
				RouteMap:       routeMap,
			}

			metadata.SetID(virtualHubId)

			return metadata.Encode(&state)
		},
	}
}

// a workaround for https://github.com/hashicorp/pandora/issues/2828
// the GET operation is defined as a LRO, while it should be a simple GET.
type VirtualHubsOutboundRouteResponse struct {
	Status     string `json:"status"`
	Properties struct {
		Output struct {
			Value []struct {
				Prefix         string `json:"prefix"`
				BgpCommunities string `json:"bgpCommunities"`
				AsPath         string `json:"asPath"`
			} `json:"value"`
		} `json:"output"`
	} `json:"properties"`
}
