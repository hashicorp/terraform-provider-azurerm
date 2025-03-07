package reports

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RequestReportCollection struct {
	Count *int64                         `json:"count,omitempty"`
	Value *[]RequestReportRecordContract `json:"value,omitempty"`
}
