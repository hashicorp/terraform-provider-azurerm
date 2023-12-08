package assessmentsmetadata

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecurityAssessmentMetadataPropertiesResponse struct {
	AssessmentType         AssessmentType                                            `json:"assessmentType"`
	Categories             *[]Categories                                             `json:"categories,omitempty"`
	Description            *string                                                   `json:"description,omitempty"`
	DisplayName            string                                                    `json:"displayName"`
	ImplementationEffort   *ImplementationEffort                                     `json:"implementationEffort,omitempty"`
	PartnerData            *SecurityAssessmentMetadataPartnerData                    `json:"partnerData,omitempty"`
	PlannedDeprecationDate *string                                                   `json:"plannedDeprecationDate,omitempty"`
	PolicyDefinitionId     *string                                                   `json:"policyDefinitionId,omitempty"`
	Preview                *bool                                                     `json:"preview,omitempty"`
	PublishDates           *SecurityAssessmentMetadataPropertiesResponsePublishDates `json:"publishDates,omitempty"`
	RemediationDescription *string                                                   `json:"remediationDescription,omitempty"`
	Severity               Severity                                                  `json:"severity"`
	Tactics                *[]Tactics                                                `json:"tactics,omitempty"`
	Techniques             *[]Techniques                                             `json:"techniques,omitempty"`
	Threats                *[]Threats                                                `json:"threats,omitempty"`
	UserImpact             *UserImpact                                               `json:"userImpact,omitempty"`
}
