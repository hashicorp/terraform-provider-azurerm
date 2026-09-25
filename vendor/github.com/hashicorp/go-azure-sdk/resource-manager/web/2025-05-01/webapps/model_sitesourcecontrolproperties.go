package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SiteSourceControlProperties struct {
	Branch                    *string                    `json:"branch,omitempty"`
	DeploymentRollbackEnabled *bool                      `json:"deploymentRollbackEnabled,omitempty"`
	GitHubActionConfiguration *GitHubActionConfiguration `json:"gitHubActionConfiguration,omitempty"`
	IsGitHubAction            *bool                      `json:"isGitHubAction,omitempty"`
	IsManualIntegration       *bool                      `json:"isManualIntegration,omitempty"`
	IsMercurial               *bool                      `json:"isMercurial,omitempty"`
	RepoURL                   *string                    `json:"repoUrl,omitempty"`
}
