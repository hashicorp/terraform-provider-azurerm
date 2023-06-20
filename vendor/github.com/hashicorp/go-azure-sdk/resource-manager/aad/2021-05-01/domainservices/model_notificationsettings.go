package domainservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NotificationSettings struct {
	AdditionalRecipients *[]string           `json:"additionalRecipients,omitempty"`
	NotifyDcAdmins       *NotifyDcAdmins     `json:"notifyDcAdmins,omitempty"`
	NotifyGlobalAdmins   *NotifyGlobalAdmins `json:"notifyGlobalAdmins,omitempty"`
}
