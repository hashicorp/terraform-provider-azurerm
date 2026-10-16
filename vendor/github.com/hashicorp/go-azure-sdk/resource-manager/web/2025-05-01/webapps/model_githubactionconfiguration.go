package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GitHubActionConfiguration struct {
	CodeConfiguration      *GitHubActionCodeConfiguration      `json:"codeConfiguration,omitempty"`
	ContainerConfiguration *GitHubActionContainerConfiguration `json:"containerConfiguration,omitempty"`
	GenerateWorkflowFile   *bool                               `json:"generateWorkflowFile,omitempty"`
	IsLinux                *bool                               `json:"isLinux,omitempty"`
}
