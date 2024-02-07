package dscpconfiguration

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DscpConfigurationPropertiesFormat struct {
	AssociatedNetworkInterfaces *[]NetworkInterface `json:"associatedNetworkInterfaces,omitempty"`
	DestinationIPRanges         *[]QosIPRange       `json:"destinationIpRanges,omitempty"`
	DestinationPortRanges       *[]QosPortRange     `json:"destinationPortRanges,omitempty"`
	Markings                    *[]int64            `json:"markings,omitempty"`
	Protocol                    *ProtocolType       `json:"protocol,omitempty"`
	ProvisioningState           *ProvisioningState  `json:"provisioningState,omitempty"`
	QosCollectionId             *string             `json:"qosCollectionId,omitempty"`
	QosDefinitionCollection     *[]QosDefinition    `json:"qosDefinitionCollection,omitempty"`
	ResourceGuid                *string             `json:"resourceGuid,omitempty"`
	SourceIPRanges              *[]QosIPRange       `json:"sourceIpRanges,omitempty"`
	SourcePortRanges            *[]QosPortRange     `json:"sourcePortRanges,omitempty"`
}
