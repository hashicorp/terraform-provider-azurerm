package hubs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NotificationHubPatchParameters struct {
	Properties *NotificationHubProperties `json:"properties,omitempty"`
	Sku        *Sku                       `json:"sku,omitempty"`
	Tags       *map[string]string         `json:"tags,omitempty"`
}
