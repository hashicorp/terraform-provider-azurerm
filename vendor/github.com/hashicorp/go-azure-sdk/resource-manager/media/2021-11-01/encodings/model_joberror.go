package encodings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobError struct {
	Category *JobErrorCategory `json:"category,omitempty"`
	Code     *JobErrorCode     `json:"code,omitempty"`
	Details  *[]JobErrorDetail `json:"details,omitempty"`
	Message  *string           `json:"message,omitempty"`
	Retry    *JobRetry         `json:"retry,omitempty"`
}
