// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sre_agent

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2026-07-01/agents"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SREAgents.AgentsClient

			id, err := agents.ParseAgentID(metadata.ResourceData.Id())
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

			return r.flatten(metadata, id, resp.Model)
		},
	}
}

func (Resource) flatten(metadata sdk.ResourceMetaData, id *agents.AgentId, input *agents.Agent) error {
	var config SREAgentModel
	if err := metadata.Decode(&config); err != nil {
		return fmt.Errorf("decoding: %+v", err)
	}

	state := SREAgentModel{
		Name:              id.AgentName,
		ResourceGroupName: id.ResourceGroupName,
	}

	if model := input; model != nil {
		state.Location = location.Normalize(model.Location)
		state.Tags = pointer.From(model.Tags)

		flattenedIdentity, err := identity.FlattenLegacySystemAndUserAssignedMapToModel(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		state.Identity = flattenedIdentity

		if props := model.Properties; props != nil {
			state.Endpoint = pointer.From(props.AgentEndpoint)
			state.OutboundIPAddresses = pointer.From(props.OutboundIPAddresses)
			state.UpgradeChannel = pointer.FromEnum(props.UpgradeChannel)
			state.ApplicationInsights = flattenSREAgentApplicationInsights(props.LogConfiguration, config.ApplicationInsights)
			state.DefaultModel = flattenSREAgentDefaultModel(props.DefaultModel)
			state.IncidentManagement = flattenSREAgentIncidentManagement(props.IncidentManagementConfiguration, config.IncidentManagement)
			state.Sandbox = flattenSREAgentSandbox(props.SandboxConfiguration)

			state.Action, err = flattenSREAgentAction(props.ActionConfiguration)
			if err != nil {
				return err
			}

			state.KnowledgeGraph, err = flattenSREAgentKnowledgeGraph(props.KnowledgeGraphConfiguration)
			if err != nil {
				return err
			}

			if props.VnetConfiguration != nil && pointer.From(props.VnetConfiguration.SubnetResourceId) != "" {
				parsed, err := commonids.ParseSubnetIDInsensitively(*props.VnetConfiguration.SubnetResourceId)
				if err != nil {
					return err
				}

				state.SubnetId = parsed.ID()
			}
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func flattenSREAgentAction(input *agents.ActionConfiguration) ([]ActionModel, error) {
	if input == nil {
		return make([]ActionModel, 0), nil
	}

	action := ActionModel{
		AccessLevel: pointer.FromEnum(input.AccessLevel),
		Mode:        pointer.FromEnum(input.Mode),
	}
	if pointer.From(input.Identity) != "" {
		parsed, err := commonids.ParseUserAssignedIdentityIDInsensitively(*input.Identity)
		if err != nil {
			return make([]ActionModel, 0), err
		}

		action.UserAssignedIdentityId = parsed.ID()
	}

	return []ActionModel{
		action,
	}, nil
}

func flattenSREAgentApplicationInsights(input *agents.LogConfiguration, config []ApplicationInsightsModel) []ApplicationInsightsModel {
	if input == nil || input.ApplicationInsightsConfiguration == nil {
		return make([]ApplicationInsightsModel, 0)
	}

	applicationInsights := ApplicationInsightsModel{
		AppId: pointer.From(input.ApplicationInsightsConfiguration.AppId),
	}

	// Azure does not return the create/update-only connection string.
	if len(config) > 0 {
		applicationInsights.ConnectionString = config[0].ConnectionString
	}

	return []ApplicationInsightsModel{
		applicationInsights,
	}
}

func flattenSREAgentDefaultModel(input *agents.DefaultModel) []DefaultModel {
	if input == nil {
		return make([]DefaultModel, 0)
	}

	return []DefaultModel{
		{
			Name:     pointer.From(input.Name),
			Provider: pointer.From(input.Provider),
		},
	}
}

func flattenSREAgentIncidentManagement(input *agents.IncidentManagementConfiguration, config []IncidentManagementModel) []IncidentManagementModel {
	if input == nil {
		return make([]IncidentManagementModel, 0)
	}

	incidentManagement := IncidentManagementModel{
		ConnectionName: pointer.From(input.ConnectionName),
		ConnectionURL:  pointer.From(input.ConnectionURL),
		OnBehalfOfUser: pointer.From(input.OboUser),
		Type:           pointer.From(input.Type),
	}

	// Azure masks the connection key on GET, so retain the configured secret.
	if len(config) > 0 {
		incidentManagement.ConnectionKey = config[0].ConnectionKey
	}

	return []IncidentManagementModel{
		incidentManagement,
	}
}

func flattenSREAgentKnowledgeGraph(input *agents.KnowledgeGraphConfiguration) ([]KnowledgeGraphModel, error) {
	if input == nil {
		return make([]KnowledgeGraphModel, 0), nil
	}

	knowledgeGraph := KnowledgeGraphModel{
		ManagedResourceIds: pointer.From(input.ManagedResources),
	}
	if pointer.From(input.Identity) != "" {
		parsed, err := commonids.ParseUserAssignedIdentityIDInsensitively(*input.Identity)
		if err != nil {
			return make([]KnowledgeGraphModel, 0), err
		}

		knowledgeGraph.UserAssignedIdentityId = parsed.ID()
	}

	return []KnowledgeGraphModel{
		knowledgeGraph,
	}, nil
}

func flattenSREAgentSandbox(input *agents.SandboxConfiguration) []SandboxModel {
	if input == nil {
		return make([]SandboxModel, 0)
	}

	sandbox := SandboxModel{
		Package: make([]SandboxPackageModel, 0),
	}
	for _, item := range pointer.From(input.Packages) {
		sandbox.Package = append(sandbox.Package, SandboxPackageModel{
			Name:           item.Name,
			PackageManager: string(item.PackageManager),
			Version:        pointer.From(item.Version),
		})
	}

	if egress := input.Egress; egress != nil {
		sandbox.Egress = []SandboxEgressModel{
			{
				Mode:                              pointer.FromEnum(egress.Mode),
				HTTPMCPServerNetworkAccessEnabled: pointer.From(egress.AllowHTTPMcpServerNetworkAccess),
				AllowedCodeRepositories:           pointer.From(egress.AllowedCodeRepositories),
				AllowedHosts:                      pointer.From(egress.AllowedHosts),
				AllowedRegistries:                 pointer.From(egress.AllowedRegistries),
			},
		}
		if proxy := egress.HTTPProxy; proxy != nil {
			sandbox.Egress[0].HTTPProxy = []SandboxHTTPProxyModel{
				{
					HTTPSURL: pointer.From(proxy.HTTPSURL),
					HTTPURL:  pointer.From(proxy.HTTPURL),
					NoProxy:  pointer.From(proxy.NoProxy),
				},
			}
		}
	}

	return []SandboxModel{
		sandbox,
	}
}
