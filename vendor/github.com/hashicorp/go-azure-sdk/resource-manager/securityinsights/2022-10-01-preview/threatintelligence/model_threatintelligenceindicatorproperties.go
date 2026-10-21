package threatintelligence

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ThreatIntelligenceIndicatorProperties struct {
	AdditionalData             *map[string]interface{}                   `json:"additionalData,omitempty"`
	Confidence                 *int64                                    `json:"confidence,omitempty"`
	Created                    *string                                   `json:"created,omitempty"`
	CreatedByRef               *string                                   `json:"createdByRef,omitempty"`
	Defanged                   *bool                                     `json:"defanged,omitempty"`
	Description                *string                                   `json:"description,omitempty"`
	DisplayName                *string                                   `json:"displayName,omitempty"`
	Extensions                 *map[string]interface{}                   `json:"extensions,omitempty"`
	ExternalId                 *string                                   `json:"externalId,omitempty"`
	ExternalLastUpdatedTimeUtc *string                                   `json:"externalLastUpdatedTimeUtc,omitempty"`
	ExternalReferences         *[]ThreatIntelligenceExternalReference    `json:"externalReferences,omitempty"`
	FriendlyName               *string                                   `json:"friendlyName,omitempty"`
	GranularMarkings           *[]ThreatIntelligenceGranularMarkingModel `json:"granularMarkings,omitempty"`
	IndicatorTypes             *[]string                                 `json:"indicatorTypes,omitempty"`
	KillChainPhases            *[]ThreatIntelligenceKillChainPhase       `json:"killChainPhases,omitempty"`
	Labels                     *[]string                                 `json:"labels,omitempty"`
	Language                   *string                                   `json:"language,omitempty"`
	LastUpdatedTimeUtc         *string                                   `json:"lastUpdatedTimeUtc,omitempty"`
	Modified                   *string                                   `json:"modified,omitempty"`
	ObjectMarkingRefs          *[]string                                 `json:"objectMarkingRefs,omitempty"`
	ParsedPattern              *[]ThreatIntelligenceParsedPattern        `json:"parsedPattern,omitempty"`
	Pattern                    *string                                   `json:"pattern,omitempty"`
	PatternType                *string                                   `json:"patternType,omitempty"`
	PatternVersion             *string                                   `json:"patternVersion,omitempty"`
	Revoked                    *bool                                     `json:"revoked,omitempty"`
	Source                     *string                                   `json:"source,omitempty"`
	ThreatIntelligenceTags     *[]string                                 `json:"threatIntelligenceTags,omitempty"`
	ThreatTypes                *[]string                                 `json:"threatTypes,omitempty"`
	ValidFrom                  *string                                   `json:"validFrom,omitempty"`
	ValidUntil                 *string                                   `json:"validUntil,omitempty"`
}
