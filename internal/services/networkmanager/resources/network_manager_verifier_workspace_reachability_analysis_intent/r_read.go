// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_verifier_workspace_reachability_analysis_intent

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/reachabilityanalysisintents"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetworkManager.ReachabilityAnalysisIntents

			id, err := reachabilityanalysisintents.ParseReachabilityAnalysisIntentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			workspaceId := reachabilityanalysisintents.NewVerifierWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName, id.VerifierWorkspaceName).ID()
			schema := ManagerVerifierWorkspaceReachabilityAnalysisIntentResourceModel{
				Name:                id.VerifierWorkspaceName,
				VerifierWorkspaceId: workspaceId,
			}

			if model := resp.Model; model != nil {
				props := model.Properties
				schema.Description = pointer.From(props.Description)
				schema.DestinationResourceId = props.DestinationResourceId
				schema.SourceResourceId = props.SourceResourceId
				schema.IpTraffic = flattenReachabilityAnalysisIntentIPTraffic(props.IPTraffic)
			}

			return metadata.Encode(&schema)
		},
	}
}

func flattenReachabilityAnalysisIntentIPTraffic(input reachabilityanalysisintents.IPTraffic) []ManagerVerifierWorkspaceReachabilityAnalysisIntentIpTraffic {
	return []ManagerVerifierWorkspaceReachabilityAnalysisIntentIpTraffic{
		{
			DestinationIps:   input.DestinationIPs,
			DestinationPorts: input.DestinationPorts,
			Protocols:        input.Protocols,
			SourceIps:        input.SourceIPs,
			SourcePorts:      input.SourcePorts,
		},
	}
}

type ManagerVerifierWorkspaceReachabilityAnalysisIntentIpTraffic struct {
	DestinationIps   []string                                      `tfschema:"destination_ips"`
	DestinationPorts []string                                      `tfschema:"destination_ports"`
	Protocols        []reachabilityanalysisintents.NetworkProtocol `tfschema:"protocols"`
	SourceIps        []string                                      `tfschema:"source_ips"`
	SourcePorts      []string                                      `tfschema:"source_ports"`
}
