package tasks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SourceUpdateParameters struct {
	Branch                      *string                   `json:"branch,omitempty"`
	RepositoryUrl               *string                   `json:"repositoryUrl,omitempty"`
	SourceControlAuthProperties *AuthInfoUpdateParameters `json:"sourceControlAuthProperties,omitempty"`
	SourceControlType           *SourceControlType        `json:"sourceControlType,omitempty"`
}
