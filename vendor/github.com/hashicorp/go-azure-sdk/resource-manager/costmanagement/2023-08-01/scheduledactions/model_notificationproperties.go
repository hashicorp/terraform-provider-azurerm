package scheduledactions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NotificationProperties struct {
	Language       *string  `json:"language,omitempty"`
	Message        *string  `json:"message,omitempty"`
	RegionalFormat *string  `json:"regionalFormat,omitempty"`
	Subject        string   `json:"subject"`
	To             []string `json:"to"`
}
