package chatmodeldeployments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ChatModelDeploymentProperties struct {
	Capacity          *int64             `json:"capacity,omitempty"`
	ModelFormat       string             `json:"modelFormat"`
	ModelName         string             `json:"modelName"`
	ModelVersion      *string            `json:"modelVersion,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	SkuName           *string            `json:"skuName,omitempty"`
}
