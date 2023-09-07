package containerapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomHostnameAnalysisResultCustomDomainVerificationFailureInfo struct {
	Code    *string                                                                          `json:"code,omitempty"`
	Details *[]CustomHostnameAnalysisResultCustomDomainVerificationFailureInfoDetailsInlined `json:"details,omitempty"`
	Message *string                                                                          `json:"message,omitempty"`
	Target  *string                                                                          `json:"target,omitempty"`
}
