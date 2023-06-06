package azuremonitorworkspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IngestionSettings struct {
	DataCollectionEndpointResourceId *string `json:"dataCollectionEndpointResourceId,omitempty"`
	DataCollectionRuleResourceId     *string `json:"dataCollectionRuleResourceId,omitempty"`
}
