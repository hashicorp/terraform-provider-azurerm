package namespaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NamespaceUpdateParameterProperties struct {
	InboundIPRules           *[]InboundIPRule                    `json:"inboundIpRules,omitempty"`
	PublicNetworkAccess      *PublicNetworkAccess                `json:"publicNetworkAccess,omitempty"`
	TopicSpacesConfiguration *UpdateTopicSpacesConfigurationInfo `json:"topicSpacesConfiguration,omitempty"`
	TopicsConfiguration      *UpdateTopicsConfigurationInfo      `json:"topicsConfiguration,omitempty"`
}
