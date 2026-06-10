package protectioncontainers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InquiryInfo struct {
	ErrorDetail    *ErrorDetail              `json:"errorDetail,omitempty"`
	InquiryDetails *[]WorkloadInquiryDetails `json:"inquiryDetails,omitempty"`
	Status         *string                   `json:"status,omitempty"`
}
