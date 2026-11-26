package securitymlanalyticssettings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AnomalySecurityMLAnalyticsCustomizableObservations struct {
	MultiSelectObservations       *[]AnomalySecurityMLAnalyticsMultiSelectObservations       `json:"multiSelectObservations,omitempty"`
	PrioritizeExcludeObservations *[]AnomalySecurityMLAnalyticsPrioritizeExcludeObservations `json:"prioritizeExcludeObservations,omitempty"`
	SingleSelectObservations      *[]AnomalySecurityMLAnalyticsSingleSelectObservations      `json:"singleSelectObservations,omitempty"`
	ThresholdObservations         *[]AnomalySecurityMLAnalyticsThresholdObservations         `json:"thresholdObservations,omitempty"`
}
