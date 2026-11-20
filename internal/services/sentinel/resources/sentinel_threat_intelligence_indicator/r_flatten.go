// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_threat_intelligence_indicator

import (
	"fmt"
	"maps"
	"regexp"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/threatintelligence"
)

func flattenThreatIntelligenceExternalReferenceModel(input *[]threatintelligence.ThreatIntelligenceExternalReference) []externalReferenceModel {
	output := make([]externalReferenceModel, 0)
	if input == nil {
		return output
	}
	for _, v := range *input {
		o := externalReferenceModel{
			SourceName:  pointer.From(v.SourceName),
			Url:         pointer.From(v.Url),
			Description: pointer.From(v.Description),
			ExternalId:  pointer.From(v.ExternalId),
		}
		if len(pointer.From(v.Hashes)) > 0 {
			o.Hashes = make(map[string]string)
			maps.Copy(o.Hashes, pointer.From(v.Hashes))
		}
		output = append(output, o)
	}
	return output
}

func flattenThreatIntelligenceGranularMarkingModelModel(input *[]threatintelligence.ThreatIntelligenceGranularMarkingModel) []granularMarkingModel {
	output := make([]granularMarkingModel, 0)
	if input == nil {
		return output
	}
	for _, v := range *input {
		output = append(output, granularMarkingModel{
			MarkingRef: pointer.From(v.MarkingRef),
			Selectors:  pointer.From(v.Selectors),
			Language:   pointer.From(v.Language),
		})
	}
	return output
}

func flattenThreatIntelligenceKillChainPhaseModel(input *[]threatintelligence.ThreatIntelligenceKillChainPhase) []killChainPhaseModel {
	output := make([]killChainPhaseModel, 0)
	if input == nil {
		return output
	}
	for _, v := range *input {
		output = append(output, killChainPhaseModel{
			PhaseName: pointer.From(v.PhaseName),
		})
	}
	return output
}

func flattenIndicatorPattern(input string) (string, error) {
	result := regexp.MustCompile(`\s=\s'(.+)'`).FindStringSubmatch(input)
	if len(result) == 2 {
		return result[1], nil
	}
	return "", fmt.Errorf("unable to parse pattern %s", input)
}

func flattenIndicatorParsedPattern(input *[]threatintelligence.ThreatIntelligenceParsedPattern) []parsedPatternModel {
	output := make([]parsedPatternModel, 0)
	if input == nil {
		return output
	}
	for _, v := range *input {
		o := parsedPatternModel{
			PatternTypeKey: pointer.From(v.PatternTypeKey),
		}
		if v.PatternTypeValues != nil {
			for _, p := range *v.PatternTypeValues {
				o.PatternTypeValues = append(o.PatternTypeValues, patternTypeValuesModel{
					Value:     pointer.From(p.Value),
					ValueType: pointer.From(p.ValueType),
				})
			}
		}
		output = append(output, o)
	}
	return output
}
