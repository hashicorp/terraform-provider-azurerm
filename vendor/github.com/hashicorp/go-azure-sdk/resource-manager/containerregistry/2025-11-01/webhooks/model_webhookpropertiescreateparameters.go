package webhooks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebhookPropertiesCreateParameters struct {
	Actions       []WebhookAction    `json:"actions"`
	CustomHeaders *map[string]string `json:"customHeaders,omitempty"`
	Scope         *string            `json:"scope,omitempty"`
	ServiceUri    string             `json:"serviceUri"`
	Status        *WebhookStatus     `json:"status,omitempty"`
}
