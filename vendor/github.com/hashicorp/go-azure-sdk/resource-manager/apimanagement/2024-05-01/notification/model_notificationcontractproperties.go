package notification

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NotificationContractProperties struct {
	Description *string                       `json:"description,omitempty"`
	Recipients  *RecipientsContractProperties `json:"recipients,omitempty"`
	Title       string                        `json:"title"`
}
