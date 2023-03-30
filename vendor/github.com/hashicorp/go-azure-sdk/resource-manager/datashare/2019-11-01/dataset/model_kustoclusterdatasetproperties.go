package dataset

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KustoClusterDataSetProperties struct {
	DataSetId              *string            `json:"dataSetId,omitempty"`
	KustoClusterResourceId string             `json:"kustoClusterResourceId"`
	Location               *string            `json:"location,omitempty"`
	ProvisioningState      *ProvisioningState `json:"provisioningState,omitempty"`
}
