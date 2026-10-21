// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventgrid_namespace

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2023-12-15-preview/namespaces"
)

func expandCreateForEventGridNamespace(model EventGridNamespaceResourceModel) (namespaces.Namespace, error) {
	identity, err := identity.ExpandSystemAndUserAssignedMapFromModel(model.Identity)
	if err != nil {
		return namespaces.Namespace{}, fmt.Errorf("expanding `identity`: %+v", err)
	}

	namespace := namespaces.Namespace{
		Identity: identity,
		Location: location.Normalize(model.Location),
		Name:     pointer.To(model.Name),
		Properties: &namespaces.NamespaceProperties{
			InboundIPRules:      expandInboundIPRules(model.InboundIpRules),
			PublicNetworkAccess: pointer.ToEnum[namespaces.PublicNetworkAccess](model.PublicNetworkAccess),
		},
		Sku: &namespaces.NamespaceSku{
			Capacity: pointer.To(model.Capacity),
			Name:     pointer.ToEnum[namespaces.SkuName](model.Sku),
		},
		Tags: pointer.To(model.Tags),
	}

	if len(model.TopicSpacesConfiguration) > 0 {
		namespace.Properties.TopicSpacesConfiguration = expandTopicSpacesConfiguration(model.TopicSpacesConfiguration)
	}

	return namespace, nil
}

func expandInboundIPRules(input []InboundIpRuleModel) *[]namespaces.InboundIPRule {
	if len(input) == 0 {
		return nil
	}

	ipRules := make([]namespaces.InboundIPRule, 0)
	for _, v := range input {
		ipRules = append(ipRules, namespaces.InboundIPRule{
			Action: pointer.ToEnum[namespaces.IPActionType](v.Action),
			IPMask: pointer.To(v.IpMask),
		})
	}
	return &ipRules
}

func expandTopicSpacesConfiguration(input []TopicSpacesConfigurationModel) *namespaces.TopicSpacesConfiguration {
	if input == nil {
		return nil
	}
	topicSpacesConfig := namespaces.TopicSpacesConfiguration{
		State: pointer.To(namespaces.TopicSpacesConfigurationStateEnabled),
		ClientAuthentication: &namespaces.ClientAuthenticationSettings{
			AlternativeAuthenticationNameSources: expandAlternativeAuthenticationNameSources(input[0].AlternativeAuthenticationNameSources),
		},
		MaximumClientSessionsPerAuthenticationName: pointer.To(input[0].MaximumClientSessionsPerAuthenticationName),
		MaximumSessionExpiryInHours:                pointer.To(input[0].MaximumSessionExpiryInHours),
		RouteTopicResourceId:                       pointer.To(input[0].RouteTopicResourceId),
		RoutingEnrichments: &namespaces.RoutingEnrichments{
			Dynamic: expandDynamicRoutingEnrichments(input[0].DynamicRoutingEnrichment),
			Static:  expandStaticRoutingEnrichments(input[0].StaticRoutingEnrichment),
		},
	}
	return &topicSpacesConfig
}

func expandTopicSpacesConfigurationUpdate(input []TopicSpacesConfigurationModel) *namespaces.UpdateTopicSpacesConfigurationInfo {
	if input == nil {
		return nil
	}

	topicSpacesConfig := namespaces.UpdateTopicSpacesConfigurationInfo{
		State: pointer.To(namespaces.TopicSpacesConfigurationStateEnabled),
		ClientAuthentication: &namespaces.ClientAuthenticationSettings{
			AlternativeAuthenticationNameSources: expandAlternativeAuthenticationNameSources(input[0].AlternativeAuthenticationNameSources),
		},
		MaximumClientSessionsPerAuthenticationName: pointer.To(input[0].MaximumClientSessionsPerAuthenticationName),
		MaximumSessionExpiryInHours:                pointer.To(input[0].MaximumSessionExpiryInHours),
		RouteTopicResourceId:                       pointer.To(input[0].RouteTopicResourceId),
		RoutingEnrichments: &namespaces.RoutingEnrichments{
			Dynamic: expandDynamicRoutingEnrichments(input[0].DynamicRoutingEnrichment),
			Static:  expandStaticRoutingEnrichments(input[0].StaticRoutingEnrichment),
		},
	}
	return &topicSpacesConfig
}

func expandAlternativeAuthenticationNameSources(input []string) *[]namespaces.AlternativeAuthenticationNameSource {
	if len(input) == 0 {
		return nil
	}

	nameSources := make([]namespaces.AlternativeAuthenticationNameSource, 0)
	for _, v := range input {
		nameSources = append(nameSources, namespaces.AlternativeAuthenticationNameSource(v))
	}
	return &nameSources
}

func expandDynamicRoutingEnrichments(input []RoutingEnrichmentModel) *[]namespaces.DynamicRoutingEnrichment {
	if len(input) == 0 {
		return nil
	}
	dynamicRoutingEnrichments := make([]namespaces.DynamicRoutingEnrichment, 0)
	for _, v := range input {
		dynamicRoutingEnrichments = append(dynamicRoutingEnrichments, namespaces.DynamicRoutingEnrichment{
			Value: pointer.To(v.Value),
			Key:   pointer.To(v.Key),
		})
	}
	return &dynamicRoutingEnrichments
}

func expandStaticRoutingEnrichments(input []RoutingEnrichmentModel) *[]namespaces.StaticRoutingEnrichment {
	if len(input) == 0 {
		return nil
	}

	staticRoutingEnrichments := make([]namespaces.StaticRoutingEnrichment, 0)
	for _, v := range input {
		staticRoutingEnrichments = append(staticRoutingEnrichments, namespaces.StaticStringRoutingEnrichment{
			Value:     pointer.To(v.Value),
			Key:       pointer.To(v.Key),
			ValueType: namespaces.StaticRoutingEnrichmentTypeString,
		})
	}
	return &staticRoutingEnrichments
}
