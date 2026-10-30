package protectioncontainers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkloadInquiryDetails struct {
	InquiryValidation *InquiryValidation `json:"inquiryValidation,omitempty"`
	ItemCount         *int64             `json:"itemCount,omitempty"`
	Type              *string            `json:"type,omitempty"`
}
