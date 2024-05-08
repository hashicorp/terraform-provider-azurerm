package staticsites

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StaticSiteBuildProperties struct {
	ApiBuildCommand                    *string `json:"apiBuildCommand,omitempty"`
	ApiLocation                        *string `json:"apiLocation,omitempty"`
	AppArtifactLocation                *string `json:"appArtifactLocation,omitempty"`
	AppBuildCommand                    *string `json:"appBuildCommand,omitempty"`
	AppLocation                        *string `json:"appLocation,omitempty"`
	GitHubActionSecretNameOverride     *string `json:"githubActionSecretNameOverride,omitempty"`
	OutputLocation                     *string `json:"outputLocation,omitempty"`
	SkipGithubActionWorkflowGeneration *bool   `json:"skipGithubActionWorkflowGeneration,omitempty"`
}
