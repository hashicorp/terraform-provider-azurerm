package datacollectionendpoints

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataCollectionEndpoint struct {
	ConfigurationAccess *ConfigurationAccessEndpointSpec              `json:"configurationAccess"`
	Description         *string                                       `json:"description,omitempty"`
	ImmutableId         *string                                       `json:"immutableId,omitempty"`
	LogsIngestion       *LogsIngestionEndpointSpec                    `json:"logsIngestion"`
	NetworkAcls         *NetworkRuleSet                               `json:"networkAcls"`
	ProvisioningState   *KnownDataCollectionEndpointProvisioningState `json:"provisioningState,omitempty"`
}
