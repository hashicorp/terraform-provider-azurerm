package webhook

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebhookUpdateProperties struct {
	Description *string            `json:"description,omitempty"`
	IsEnabled   *bool              `json:"isEnabled,omitempty"`
	Parameters  *map[string]string `json:"parameters,omitempty"`
	RunOn       *string            `json:"runOn,omitempty"`
}
