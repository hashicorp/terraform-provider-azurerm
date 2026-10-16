package dashboards

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MarkdownPartMetadataSettingsContent struct {
	Content        *string `json:"content,omitempty"`
	MarkdownSource *int64  `json:"markdownSource,omitempty"`
	MarkdownUri    *string `json:"markdownUri,omitempty"`
	Subtitle       *string `json:"subtitle,omitempty"`
	Title          *string `json:"title,omitempty"`
}
