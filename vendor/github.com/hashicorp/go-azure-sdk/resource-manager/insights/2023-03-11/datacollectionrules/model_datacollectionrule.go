package datacollectionrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataCollectionRule struct {
	AgentSettings            *AgentSettingsSpec                        `json:"agentSettings,omitempty"`
	DataCollectionEndpointId *string                                   `json:"dataCollectionEndpointId,omitempty"`
	DataFlows                *[]DataFlow                               `json:"dataFlows,omitempty"`
	DataSources              *DataSourcesSpec                          `json:"dataSources,omitempty"`
	Description              *string                                   `json:"description,omitempty"`
	Destinations             *DestinationsSpec                         `json:"destinations,omitempty"`
	Endpoints                *EndpointsSpec                            `json:"endpoints,omitempty"`
	ImmutableId              *string                                   `json:"immutableId,omitempty"`
	Metadata                 *Metadata                                 `json:"metadata,omitempty"`
	ProvisioningState        *KnownDataCollectionRuleProvisioningState `json:"provisioningState,omitempty"`
	References               *ReferencesSpec                           `json:"references,omitempty"`
	StreamDeclarations       *map[string]StreamDeclaration             `json:"streamDeclarations,omitempty"`
}
