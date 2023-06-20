package scheduledactions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NotificationProperties struct {
	Message *string  `json:"message,omitempty"`
	Subject string   `json:"subject"`
	To      []string `json:"to"`
}
