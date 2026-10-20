// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventgrid_namespace

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2023-12-15-preview/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/topics"
)

func flattenInboundIPRules(ipRules *[]namespaces.InboundIPRule) []InboundIpRuleModel {
	output := make([]InboundIpRuleModel, 0)

	if ipRules == nil || len(*ipRules) == 0 {
		return output
	}

	for _, v := range *ipRules {
		output = append(output, InboundIpRuleModel{
			IpMask: pointer.From(v.IPMask),
			Action: pointer.FromEnum(v.Action),
		})
	}
	return output
}

func flattenTopicSpacesConfiguration(topicSpacesConfig *namespaces.TopicSpacesConfiguration) ([]TopicSpacesConfigurationModel, error) {
	var output TopicSpacesConfigurationModel
	if topicSpacesConfig == nil {
		return []TopicSpacesConfigurationModel{}, nil
	}

	output.MaximumSessionExpiryInHours = pointer.From(topicSpacesConfig.MaximumSessionExpiryInHours)
	output.MaximumClientSessionsPerAuthenticationName = pointer.From(topicSpacesConfig.MaximumClientSessionsPerAuthenticationName)
	var routeId string
	if topicSpacesConfig.RouteTopicResourceId != nil && *topicSpacesConfig.RouteTopicResourceId != "" {
		id, err := topics.ParseTopicID(*topicSpacesConfig.RouteTopicResourceId)
		if err != nil {
			return nil, err
		}
		routeId = id.ID()
	}
	output.RouteTopicResourceId = routeId

	if topicSpacesConfig.ClientAuthentication != nil {
		output.AlternativeAuthenticationNameSources = flattenAlternativeAuthenticationNameSources(topicSpacesConfig.ClientAuthentication.AlternativeAuthenticationNameSources)
	}
	if topicSpacesConfig.RoutingEnrichments != nil {
		output.DynamicRoutingEnrichment = flattenDynamicRoutingEnrichments(topicSpacesConfig.RoutingEnrichments.Dynamic)
		output.StaticRoutingEnrichment = flattenStaticRoutingEnrichments(topicSpacesConfig.RoutingEnrichments.Static)
	}

	return []TopicSpacesConfigurationModel{output}, nil
}

func flattenAlternativeAuthenticationNameSources(nameSources *[]namespaces.AlternativeAuthenticationNameSource) []string {
	output := make([]string, 0)

	if nameSources == nil || len(*nameSources) == 0 {
		return output
	}

	for _, v := range *nameSources {
		output = append(output, string(v))
	}
	return output
}

func flattenDynamicRoutingEnrichments(dynamicRoutingEnrichments *[]namespaces.DynamicRoutingEnrichment) []RoutingEnrichmentModel {
	output := make([]RoutingEnrichmentModel, 0)
	if dynamicRoutingEnrichments == nil || len(*dynamicRoutingEnrichments) == 0 {
		return output
	}

	for _, v := range *dynamicRoutingEnrichments {
		output = append(output, RoutingEnrichmentModel{
			Value: pointer.From(v.Value),
			Key:   pointer.From(v.Key),
		})
	}
	return output
}

func flattenStaticRoutingEnrichments(staticRoutingEnrichments *[]namespaces.StaticRoutingEnrichment) []RoutingEnrichmentModel {
	output := make([]RoutingEnrichmentModel, 0)
	if staticRoutingEnrichments == nil || len(*staticRoutingEnrichments) == 0 {
		return output
	}

	for _, v := range *staticRoutingEnrichments {
		output = append(output, RoutingEnrichmentModel{
			Value: pointer.From(v.(namespaces.StaticStringRoutingEnrichment).Value),
			Key:   pointer.From(v.(namespaces.StaticStringRoutingEnrichment).Key),
		})
	}
	return output
}
