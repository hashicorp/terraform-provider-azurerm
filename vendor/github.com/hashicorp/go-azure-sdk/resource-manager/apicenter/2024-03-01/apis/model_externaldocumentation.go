package apis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExternalDocumentation struct {
	Description *string `json:"description,omitempty"`
	Title       *string `json:"title,omitempty"`
	Url         string  `json:"url"`
}
