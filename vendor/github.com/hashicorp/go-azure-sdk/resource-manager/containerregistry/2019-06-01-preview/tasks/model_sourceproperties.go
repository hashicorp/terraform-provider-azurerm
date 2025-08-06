package tasks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SourceProperties struct {
	Branch                      *string           `json:"branch,omitempty"`
	RepositoryURL               string            `json:"repositoryUrl"`
	SourceControlAuthProperties *AuthInfo         `json:"sourceControlAuthProperties,omitempty"`
	SourceControlType           SourceControlType `json:"sourceControlType"`
}
