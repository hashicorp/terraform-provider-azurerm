package containerapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiagnosticRendering struct {
	Description *string `json:"description,omitempty"`
	IsVisible   *bool   `json:"isVisible,omitempty"`
	Title       *string `json:"title,omitempty"`
	Type        *int64  `json:"type,omitempty"`
}
