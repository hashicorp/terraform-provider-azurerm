package appserviceenvironments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StatusCodesBasedTrigger struct {
	Count        *int64  `json:"count,omitempty"`
	Path         *string `json:"path,omitempty"`
	Status       *int64  `json:"status,omitempty"`
	SubStatus    *int64  `json:"subStatus,omitempty"`
	TimeInterval *string `json:"timeInterval,omitempty"`
	Win32Status  *int64  `json:"win32Status,omitempty"`
}
