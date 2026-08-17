package jobdefinitions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TargetEndpointProperties struct {
	AzureStorageAccountResourceId *string `json:"azureStorageAccountResourceId,omitempty"`
	AzureStorageBlobContainerName *string `json:"azureStorageBlobContainerName,omitempty"`
	Name                          *string `json:"name,omitempty"`
	TargetEndpointResourceId      *string `json:"targetEndpointResourceId,omitempty"`
}
