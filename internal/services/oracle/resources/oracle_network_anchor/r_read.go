// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_network_anchor

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/networkanchors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/resourceanchors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := networkanchors.ParseNetworkAnchorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			client := metadata.Client.Oracle.OracleClient.NetworkAnchors
			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			return r.flatten(metadata, id, resp.Model)
		},
	}
}

func (Resource) flatten(metadata sdk.ResourceMetaData, id *networkanchors.NetworkAnchorId, model *networkanchors.NetworkAnchor) error {
	state := NetworkAnchorResourceModel{
		Name:              id.NetworkAnchorName,
		ResourceGroupName: id.ResourceGroupName,
	}

	if model != nil {
		state.Location = location.Normalize(model.Location)
		state.Tags = pointer.From(model.Tags)
		state.Zones = pointer.From(model.Zones)
		if props := model.Properties; props != nil {
			resourceAnchorID, err := resourceanchors.ParseResourceAnchorIDInsensitively(props.ResourceAnchorId)
			if err != nil {
				return err
			}
			state.ResourceAnchorId = resourceAnchorID.ID()

			subnetID, err := commonids.ParseSubnetIDInsensitively(props.SubnetId)
			if err != nil {
				return err
			}
			state.SubnetId = subnetID.ID()

			state.OciVcnDnsLabel = pointer.From(props.OciVcnDnsLabel)
			state.DnsPrivateZoneReplicationEnabled = pointer.From(props.IsOracleToAzureDnsZoneSyncEnabled)
			state.DnsForwardingEndpointIpAddress = pointer.From(props.DnsForwardingEndpointIPAddress)
			state.DnsForwardingRuleUrl = pointer.From(props.DnsForwardingRulesURL)
			state.DnsForwardingEndpointNsgRuleUrl = pointer.From(props.DnsForwardingEndpointNsgRulesURL)
			state.DnsListeningEndpointIpAddress = pointer.From(props.DnsListeningEndpointIPAddress)
			state.DnsListeningEndpointNsgRuleUrl = pointer.From(props.DnsListeningEndpointNsgRulesURL)

			state.BackupSubnetCidr = pointer.From(props.OciBackupCidrBlock)
			state.DnsForwardingRule = flattenDnsForwardingRules(props.DnsForwardingRules)
			state.DnsListeningEndpointAllowedCidrs = flattenNetworkAnchorCommaSeparatedValues(pointer.From(props.DnsListeningEndpointAllowedCidrs))
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}
