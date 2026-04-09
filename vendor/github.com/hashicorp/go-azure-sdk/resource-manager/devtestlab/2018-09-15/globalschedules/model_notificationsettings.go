package globalschedules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NotificationSettings struct {
	EmailRecipient     *string       `json:"emailRecipient,omitempty"`
	NotificationLocale *string       `json:"notificationLocale,omitempty"`
	Status             *EnableStatus `json:"status,omitempty"`
	TimeInMinutes      *int64        `json:"timeInMinutes,omitempty"`
	WebhookURL         *string       `json:"webhookUrl,omitempty"`
}
