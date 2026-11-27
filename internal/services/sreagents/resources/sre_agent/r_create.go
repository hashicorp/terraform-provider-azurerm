// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sre_agent

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2026-07-01/agents"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SREAgents.AgentsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config SREAgentModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := agents.NewAgentID(subscriptionId, config.ResourceGroupName, config.Name)
			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(config.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			param := agents.Agent{
				Location: location.Normalize(config.Location),
				Identity: expandedIdentity,
				Properties: &agents.AgentProperties{
					ActionConfiguration:             expandSREAgentAction(config.Action),
					LogConfiguration:                expandSREAgentApplicationInsights(config.ApplicationInsights),
					DefaultModel:                    expandSREAgentDefaultModel(config.DefaultModel),
					IncidentManagementConfiguration: expandSREAgentIncidentManagement(config.IncidentManagement),
					KnowledgeGraphConfiguration:     expandSREAgentKnowledgeGraph(config.KnowledgeGraph),
					SandboxConfiguration:            expandSREAgentSandbox(config.Sandbox),
					UpgradeChannel:                  pointer.ToEnum[agents.UpgradeChannel](config.UpgradeChannel),
				},
				Tags: pointer.To(config.Tags),
			}

			if config.SubnetId != "" {
				param.Properties.VnetConfiguration = &agents.VnetConfiguration{
					SubnetResourceId: pointer.To(config.SubnetId),
				}
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, param, metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}

func expandSREAgentAction(input []ActionModel) *agents.ActionConfiguration {
	if len(input) == 0 {
		return nil
	}

	action := input[0]
	result := &agents.ActionConfiguration{}
	if action.AccessLevel != "" {
		result.AccessLevel = pointer.ToEnum[agents.AgentAccessLevel](action.AccessLevel)
	}

	if action.Mode != "" {
		result.Mode = pointer.ToEnum[agents.AgentMode](action.Mode)
	}

	if action.UserAssignedIdentityId != "" {
		result.Identity = pointer.To(action.UserAssignedIdentityId)
	}

	return result
}

func expandSREAgentApplicationInsights(input []ApplicationInsightsModel) *agents.LogConfiguration {
	if len(input) == 0 {
		return nil
	}

	applicationInsights := input[0]
	return &agents.LogConfiguration{
		ApplicationInsightsConfiguration: &agents.ApplicationInsightsConfiguration{
			AppId:            pointer.To(applicationInsights.AppId),
			ConnectionString: pointer.To(applicationInsights.ConnectionString),
		},
	}
}

func expandSREAgentDefaultModel(input []DefaultModel) *agents.DefaultModel {
	if len(input) == 0 {
		return nil
	}

	defaultModel := input[0]
	return &agents.DefaultModel{
		Name:     pointer.To(defaultModel.Name),
		Provider: pointer.To(defaultModel.Provider),
	}
}

func expandSREAgentIncidentManagement(input []IncidentManagementModel) *agents.IncidentManagementConfiguration {
	if len(input) == 0 {
		return nil
	}

	incidentManagement := input[0]
	result := &agents.IncidentManagementConfiguration{
		Type:          pointer.To(incidentManagement.Type),
		ConnectionKey: pointer.To(incidentManagement.ConnectionKey),
		ConnectionURL: pointer.To(incidentManagement.ConnectionURL),
		OboUser:       pointer.To(incidentManagement.OnBehalfOfUser),
	}
	if incidentManagement.ConnectionName != "" {
		result.ConnectionName = pointer.To(incidentManagement.ConnectionName)
	}

	return result
}

func expandSREAgentKnowledgeGraph(input []KnowledgeGraphModel) *agents.KnowledgeGraphConfiguration {
	if len(input) == 0 {
		return nil
	}

	knowledgeGraph := input[0]
	result := &agents.KnowledgeGraphConfiguration{
		ManagedResources: pointer.To(knowledgeGraph.ManagedResourceIds),
	}
	if knowledgeGraph.UserAssignedIdentityId != "" {
		result.Identity = pointer.To(knowledgeGraph.UserAssignedIdentityId)
	}

	return result
}

func expandSREAgentSandbox(input []SandboxModel) *agents.SandboxConfiguration {
	if len(input) == 0 {
		return nil
	}

	sandbox := input[0]
	packages := make([]agents.SandboxPackageSpec, 0, len(sandbox.Package))
	for _, item := range sandbox.Package {
		spec := agents.SandboxPackageSpec{
			Name:           item.Name,
			PackageManager: agents.SandboxPackageManager(item.PackageManager),
		}
		if item.Version != "" {
			spec.Version = pointer.To(item.Version)
		}

		packages = append(packages, spec)
	}

	result := &agents.SandboxConfiguration{
		Packages: pointer.To(packages),
	}
	if len(sandbox.Egress) > 0 {
		egress := sandbox.Egress[0]
		result.Egress = &agents.SandboxEgressConfiguration{
			AllowHTTPMcpServerNetworkAccess: pointer.To(egress.HTTPMCPServerNetworkAccessEnabled),
			AllowedCodeRepositories:         pointer.To(egress.AllowedCodeRepositories),
			AllowedHosts:                    pointer.To(egress.AllowedHosts),
			AllowedRegistries:               pointer.To(egress.AllowedRegistries),
		}
		if egress.Mode != "" {
			result.Egress.Mode = pointer.ToEnum[agents.SandboxEgressMode](egress.Mode)
		}

		if len(egress.HTTPProxy) > 0 {
			proxy := egress.HTTPProxy[0]
			result.Egress.HTTPProxy = &agents.SandboxHTTPProxyConfiguration{
				HTTPSURL: pointer.To(proxy.HTTPSURL),
				NoProxy:  pointer.To(proxy.NoProxy),
			}
			if proxy.HTTPURL != "" {
				result.Egress.HTTPProxy.HTTPURL = pointer.To(proxy.HTTPURL)
			}
		}
	}

	return result
}
