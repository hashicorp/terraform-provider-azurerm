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

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetworkManager.ReachabilityAnalysisIntents
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config ManagerVerifierWorkspaceReachabilityAnalysisIntentResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			workspaceId, err := reachabilityanalysisintents.ParseVerifierWorkspaceID(config.VerifierWorkspaceId)
			if err != nil {
				return err
			}

			id := reachabilityanalysisintents.NewReachabilityAnalysisIntentID(subscriptionId, workspaceId.ResourceGroupName, workspaceId.NetworkManagerName, workspaceId.VerifierWorkspaceName, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			payload := reachabilityanalysisintents.ReachabilityAnalysisIntent{
				Name: pointer.To(config.Name),
				Properties: reachabilityanalysisintents.ReachabilityAnalysisIntentProperties{
					Description:           pointer.To(config.Description),
					SourceResourceId:      config.SourceResourceId,
					DestinationResourceId: config.DestinationResourceId,
					IPTraffic:             expandReachabilityAnalysisIntentIPTraffic(config.IpTraffic),
				},
			}

			if _, err := client.Create(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func expandReachabilityAnalysisIntentIPTraffic(input []ManagerVerifierWorkspaceReachabilityAnalysisIntentIpTraffic) reachabilityanalysisintents.IPTraffic {
	if len(input) == 0 {
		return reachabilityanalysisintents.IPTraffic{}
	}

	return reachabilityanalysisintents.IPTraffic{
		DestinationIPs:   input[0].DestinationIps,
		DestinationPorts: input[0].DestinationPorts,
		Protocols:        input[0].Protocols,
		SourceIPs:        input[0].SourceIps,
		SourcePorts:      input[0].SourcePorts,
	}
}
