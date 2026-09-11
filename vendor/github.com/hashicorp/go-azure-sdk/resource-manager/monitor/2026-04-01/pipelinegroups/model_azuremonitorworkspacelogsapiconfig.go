package pipelinegroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureMonitorWorkspaceLogsApiConfig struct {
	DataCollectionEndpointURL string    `json:"dataCollectionEndpointUrl"`
	DataCollectionRule        string    `json:"dataCollectionRule"`
	Schema                    SchemaMap `json:"schema"`
	Stream                    string    `json:"stream"`
}
