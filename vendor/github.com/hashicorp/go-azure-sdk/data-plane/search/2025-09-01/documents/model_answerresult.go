package documents

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AnswerResult struct {
	Highlights *string  `json:"highlights,omitempty"`
	Key        *string  `json:"key,omitempty"`
	Score      *float64 `json:"score,omitempty"`
	Text       *string  `json:"text,omitempty"`
}
