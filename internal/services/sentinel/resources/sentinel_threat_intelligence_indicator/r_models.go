// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_threat_intelligence_indicator

type IndicatorModel struct {
	Name                       string                   `tfschema:"guid"`
	WorkspaceId                string                   `tfschema:"workspace_id"`
	Confidence                 int64                    `tfschema:"confidence"`
	CreatedByRef               string                   `tfschema:"created_by"`
	Description                string                   `tfschema:"description"`
	DisplayName                string                   `tfschema:"display_name"`
	Extensions                 string                   `tfschema:"extension"`
	ExternalReference          []externalReferenceModel `tfschema:"external_reference"`
	GranularMarkings           []granularMarkingModel   `tfschema:"granular_marking"`
	IndicatorTypes             []string                 `tfschema:"indicator_type"`
	KillChainPhases            []killChainPhaseModel    `tfschema:"kill_chain_phase"`
	Labels                     []string                 `tfschema:"tags"`
	Language                   string                   `tfschema:"language"`
	ObjectMarking              []string                 `tfschema:"object_marking_refs"`
	ParsedPattern              []parsedPatternModel     `tfschema:"parsed_pattern"`
	Pattern                    string                   `tfschema:"pattern"`
	PatternType                string                   `tfschema:"pattern_type"`
	PatternVersion             string                   `tfschema:"pattern_version"`
	Revoked                    bool                     `tfschema:"revoked"`
	Source                     string                   `tfschema:"source"`
	ThreatTypes                []string                 `tfschema:"threat_types"`
	ValidFrom                  string                   `tfschema:"validate_from_utc"`
	ValidUntil                 string                   `tfschema:"validate_until_utc"`
	CreatedOn                  string                   `tfschema:"created_on"`
	ExternalId                 string                   `tfschema:"external_id"`
	ExternalLastUpdatedTimeUtc string                   `tfschema:"external_last_updated_time_utc"`
	LastUpdatedTimeUtc         string                   `tfschema:"last_updated_time_utc"`
	Defanged                   bool                     `tfschema:"defanged"`
}

type externalReferenceModel struct {
	SourceName  string            `tfschema:"source_name"`
	Url         string            `tfschema:"url"`
	Hashes      map[string]string `tfschema:"hashes"`
	Description string            `tfschema:"description"`
	ExternalId  string            `tfschema:"id"`
}

type granularMarkingModel struct {
	MarkingRef string   `tfschema:"marking_ref"`
	Selectors  []string `tfschema:"selectors"`
	Language   string   `tfschema:"language"`
}

type killChainPhaseModel struct {
	PhaseName string `tfschema:"name"`
}

type parsedPatternModel struct {
	PatternTypeValues []patternTypeValuesModel `tfschema:"pattern_type_values"`
	PatternTypeKey    string                   `tfschema:"pattern_type_key"`
}

type patternTypeValuesModel struct {
	Value     string `tfschema:"value"`
	ValueType string `tfschema:"value_type"`
}

func (r Resource) ModelObject() any {
	return &IndicatorModel{}
}
