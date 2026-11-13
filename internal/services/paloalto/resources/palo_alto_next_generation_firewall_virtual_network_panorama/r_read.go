// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_next_generation_firewall_virtual_network_panorama

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/firewallresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/schema"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.FirewallResources

			id, err := firewallresources.ParseFirewallID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state NextGenerationFirewallVnetPanoramaModel

			existing, err := client.FirewallsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			state.Name = id.FirewallName
			state.ResourceGroupName = id.ResourceGroupName
			if model := existing.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				props := model.Properties

				state.DNSSettings = schema.FlattenDNSSettings(props.DnsSettings)

				state.NetworkProfile = schema.FlattenNetworkProfileVnet(props.NetworkProfile)

				state.FrontEnd = schema.FlattenDestinationNAT(props.FrontEndSettings)

				if panoramaConfig := props.PanoramaConfig; panoramaConfig != nil {
					state.PanoramaBase64Config = panoramaConfig.ConfigString
					state.PanoramaConfig = []schema.Panorama{{
						Name:            pointer.From(panoramaConfig.CgName),
						DeviceGroupName: pointer.From(panoramaConfig.DgName),
						HostName:        pointer.From(panoramaConfig.HostName),
						PanoramaServer:  pointer.From(panoramaConfig.PanoramaServer),
						PanoramaServer2: pointer.From(panoramaConfig.PanoramaServer2),
						TplName:         pointer.From(panoramaConfig.TplName),
						VMAuthKey:       pointer.From(panoramaConfig.VMAuthKey),
					}}
				}

				state.MarketplaceOfferId = props.MarketplaceDetails.OfferId

				state.PlanId = props.PlanData.PlanId

				state.Tags = tags.Flatten(model.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}
