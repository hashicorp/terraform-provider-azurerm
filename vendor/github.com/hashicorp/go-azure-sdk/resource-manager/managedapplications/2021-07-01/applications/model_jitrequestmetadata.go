package applications

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JitRequestMetadata struct {
	OriginRequestId    *string `json:"originRequestId,omitempty"`
	RequestorId        *string `json:"requestorId,omitempty"`
	SubjectDisplayName *string `json:"subjectDisplayName,omitempty"`
	TenantDisplayName  *string `json:"tenantDisplayName,omitempty"`
}
