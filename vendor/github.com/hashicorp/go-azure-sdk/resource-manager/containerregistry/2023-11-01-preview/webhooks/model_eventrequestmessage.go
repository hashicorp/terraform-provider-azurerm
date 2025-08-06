package webhooks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventRequestMessage struct {
	Content    *EventContent      `json:"content,omitempty"`
	Headers    *map[string]string `json:"headers,omitempty"`
	Method     *string            `json:"method,omitempty"`
	RequestUri *string            `json:"requestUri,omitempty"`
	Version    *string            `json:"version,omitempty"`
}
