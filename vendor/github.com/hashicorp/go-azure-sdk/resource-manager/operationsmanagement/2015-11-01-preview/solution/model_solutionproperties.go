package solution

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SolutionProperties struct {
	ContainedResources  *[]string `json:"containedResources,omitempty"`
	ProvisioningState   *string   `json:"provisioningState,omitempty"`
	ReferencedResources *[]string `json:"referencedResources,omitempty"`
	WorkspaceResourceId string    `json:"workspaceResourceId"`
}
