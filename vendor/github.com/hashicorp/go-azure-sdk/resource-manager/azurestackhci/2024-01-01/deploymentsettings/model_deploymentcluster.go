package deploymentsettings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentCluster struct {
	AzureServiceEndpoint *string `json:"azureServiceEndpoint,omitempty"`
	CloudAccountName     *string `json:"cloudAccountName,omitempty"`
	Name                 *string `json:"name,omitempty"`
	WitnessPath          *string `json:"witnessPath,omitempty"`
	WitnessType          *string `json:"witnessType,omitempty"`
}
