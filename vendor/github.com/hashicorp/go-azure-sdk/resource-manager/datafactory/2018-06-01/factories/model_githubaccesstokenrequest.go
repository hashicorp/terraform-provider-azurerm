package factories

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GitHubAccessTokenRequest struct {
	GitHubAccessCode         string              `json:"gitHubAccessCode"`
	GitHubAccessTokenBaseUrl string              `json:"gitHubAccessTokenBaseUrl"`
	GitHubClientId           *string             `json:"gitHubClientId,omitempty"`
	GitHubClientSecret       *GitHubClientSecret `json:"gitHubClientSecret,omitempty"`
}
