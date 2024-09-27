package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomHostnameAnalysisResultProperties struct {
	ARecords                            *[]string                  `json:"aRecords,omitempty"`
	AlternateCNameRecords               *[]string                  `json:"alternateCNameRecords,omitempty"`
	AlternateTxtRecords                 *[]string                  `json:"alternateTxtRecords,omitempty"`
	CNameRecords                        *[]string                  `json:"cNameRecords,omitempty"`
	ConflictingAppResourceId            *string                    `json:"conflictingAppResourceId,omitempty"`
	CustomDomainVerificationFailureInfo *ErrorEntity               `json:"customDomainVerificationFailureInfo,omitempty"`
	CustomDomainVerificationTest        *DnsVerificationTestResult `json:"customDomainVerificationTest,omitempty"`
	HasConflictAcrossSubscription       *bool                      `json:"hasConflictAcrossSubscription,omitempty"`
	HasConflictOnScaleUnit              *bool                      `json:"hasConflictOnScaleUnit,omitempty"`
	IsHostnameAlreadyVerified           *bool                      `json:"isHostnameAlreadyVerified,omitempty"`
	TxtRecords                          *[]string                  `json:"txtRecords,omitempty"`
}
