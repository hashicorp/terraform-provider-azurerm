package scheduledqueryrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzNsActionGroup struct {
	ActionGroup          *[]string `json:"actionGroup,omitempty"`
	CustomWebhookPayload *string   `json:"customWebhookPayload,omitempty"`
	EmailSubject         *string   `json:"emailSubject,omitempty"`
}
