package documents

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CaptionResult struct {
	Highlights *string `json:"highlights,omitempty"`
	Text       *string `json:"text,omitempty"`
}
