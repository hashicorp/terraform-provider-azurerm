package webhooks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventResponseMessage struct {
	Content      *string            `json:"content,omitempty"`
	Headers      *map[string]string `json:"headers,omitempty"`
	ReasonPhrase *string            `json:"reasonPhrase,omitempty"`
	StatusCode   *string            `json:"statusCode,omitempty"`
	Version      *string            `json:"version,omitempty"`
}
