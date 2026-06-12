package deployments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TemplateLink struct {
	ContentVersion *string `json:"contentVersion,omitempty"`
	Id             *string `json:"id,omitempty"`
	QueryString    *string `json:"queryString,omitempty"`
	RelativePath   *string `json:"relativePath,omitempty"`
	Uri            *string `json:"uri,omitempty"`
}
