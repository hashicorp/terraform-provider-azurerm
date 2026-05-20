package staticsites

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StaticSiteTemplateOptions struct {
	Description           *string `json:"description,omitempty"`
	IsPrivate             *bool   `json:"isPrivate,omitempty"`
	Owner                 *string `json:"owner,omitempty"`
	RepositoryName        *string `json:"repositoryName,omitempty"`
	TemplateRepositoryURL *string `json:"templateRepositoryUrl,omitempty"`
}
