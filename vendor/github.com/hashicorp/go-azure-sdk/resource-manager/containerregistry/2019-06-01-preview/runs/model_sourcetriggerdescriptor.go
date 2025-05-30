package runs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SourceTriggerDescriptor struct {
	BranchName    *string `json:"branchName,omitempty"`
	CommitId      *string `json:"commitId,omitempty"`
	EventType     *string `json:"eventType,omitempty"`
	Id            *string `json:"id,omitempty"`
	ProviderType  *string `json:"providerType,omitempty"`
	PullRequestId *string `json:"pullRequestId,omitempty"`
	RepositoryURL *string `json:"repositoryUrl,omitempty"`
}
