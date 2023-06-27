package autoscalesettings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EmailNotification struct {
	CustomEmails                       *[]string `json:"customEmails,omitempty"`
	SendToSubscriptionAdministrator    *bool     `json:"sendToSubscriptionAdministrator,omitempty"`
	SendToSubscriptionCoAdministrators *bool     `json:"sendToSubscriptionCoAdministrators,omitempty"`
}
