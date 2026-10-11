package jobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TaskFailureInformation struct {
	Category ErrorCategory    `json:"category"`
	Code     *string          `json:"code,omitempty"`
	Details  *[]NameValuePair `json:"details,omitempty"`
	Message  *string          `json:"message,omitempty"`
}
