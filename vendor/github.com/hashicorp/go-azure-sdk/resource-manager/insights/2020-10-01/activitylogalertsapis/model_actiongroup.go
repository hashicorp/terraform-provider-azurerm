package activitylogalertsapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActionGroup struct {
	ActionGroupId     string             `json:"actionGroupId"`
	WebhookProperties *map[string]string `json:"webhookProperties,omitempty"`
}
