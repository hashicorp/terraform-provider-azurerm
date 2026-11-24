// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sre_agent

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2026-07-01/agents"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SREAgents.AgentsClient

			id, err := agents.ParseAgentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config SREAgentModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			// Azure omits the connection string from GET responses, but PUT requires it alongside AppId.
			existing.Model.Properties.LogConfiguration = expandSREAgentApplicationInsights(config.ApplicationInsights)

			// Azure omits the connection key from GET responses, so resend the configured secret.
			existing.Model.Properties.IncidentManagementConfiguration = expandSREAgentIncidentManagement(config.IncidentManagement)

			if metadata.ResourceData.HasChange("identity") {
				expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(config.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				existing.Model.Identity = expandedIdentity
			}

			if metadata.ResourceData.HasChange("tags") {
				existing.Model.Tags = pointer.To(config.Tags)
			}

			if metadata.ResourceData.HasChange("action") {
				existing.Model.Properties.ActionConfiguration = expandSREAgentAction(config.Action)
			}

			if metadata.ResourceData.HasChange("default_model") {
				existing.Model.Properties.DefaultModel = expandSREAgentDefaultModel(config.DefaultModel)
			}

			if metadata.ResourceData.HasChange("knowledge_graph") {
				existing.Model.Properties.KnowledgeGraphConfiguration = expandSREAgentKnowledgeGraph(config.KnowledgeGraph)
			}

			if metadata.ResourceData.HasChange("sandbox") {
				existing.Model.Properties.SandboxConfiguration = expandSREAgentSandbox(config.Sandbox)
			}

			if metadata.ResourceData.HasChange("subnet_id") {
				existing.Model.Properties.VnetConfiguration = nil
				if config.SubnetId != "" {
					existing.Model.Properties.VnetConfiguration = &agents.VnetConfiguration{
						SubnetResourceId: pointer.To(config.SubnetId),
					}
				}
			}

			if metadata.ResourceData.HasChange("upgrade_channel") {
				existing.Model.Properties.UpgradeChannel = pointer.ToEnum[agents.UpgradeChannel](config.UpgradeChannel)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
