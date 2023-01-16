package webtestsapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebTestPropertiesValidationRules struct {
	ContentValidation             *WebTestPropertiesValidationRulesContentValidation `json:"ContentValidation,omitempty"`
	ExpectedHTTPStatusCode        *int64                                             `json:"ExpectedHttpStatusCode,omitempty"`
	IgnoreHTTPSStatusCode         *bool                                              `json:"IgnoreHttpsStatusCode,omitempty"`
	SSLCertRemainingLifetimeCheck *int64                                             `json:"SSLCertRemainingLifetimeCheck,omitempty"`
	SSLCheck                      *bool                                              `json:"SSLCheck,omitempty"`
}
