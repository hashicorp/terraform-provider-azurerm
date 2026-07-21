package threatintelligence

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ThreatIntelligenceMetric struct {
	LastUpdatedTimeUtc *string                           `json:"lastUpdatedTimeUtc,omitempty"`
	PatternTypeMetrics *[]ThreatIntelligenceMetricEntity `json:"patternTypeMetrics,omitempty"`
	SourceMetrics      *[]ThreatIntelligenceMetricEntity `json:"sourceMetrics,omitempty"`
	ThreatTypeMetrics  *[]ThreatIntelligenceMetricEntity `json:"threatTypeMetrics,omitempty"`
}
