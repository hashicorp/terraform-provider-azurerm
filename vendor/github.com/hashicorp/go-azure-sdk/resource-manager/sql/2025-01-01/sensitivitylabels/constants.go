package sensitivitylabels

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClientClassificationSource string

const (
	ClientClassificationSourceMIP         ClientClassificationSource = "MIP"
	ClientClassificationSourceNative      ClientClassificationSource = "Native"
	ClientClassificationSourceNone        ClientClassificationSource = "None"
	ClientClassificationSourceRecommended ClientClassificationSource = "Recommended"
)

func PossibleValuesForClientClassificationSource() []string {
	return []string{
		string(ClientClassificationSourceMIP),
		string(ClientClassificationSourceNative),
		string(ClientClassificationSourceNone),
		string(ClientClassificationSourceRecommended),
	}
}

func (s *ClientClassificationSource) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseClientClassificationSource(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseClientClassificationSource(input string) (*ClientClassificationSource, error) {
	vals := map[string]ClientClassificationSource{
		"mip":         ClientClassificationSourceMIP,
		"native":      ClientClassificationSourceNative,
		"none":        ClientClassificationSourceNone,
		"recommended": ClientClassificationSourceRecommended,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClientClassificationSource(input)
	return &out, nil
}

type RecommendedSensitivityLabelUpdateKind string

const (
	RecommendedSensitivityLabelUpdateKindDisable RecommendedSensitivityLabelUpdateKind = "disable"
	RecommendedSensitivityLabelUpdateKindEnable  RecommendedSensitivityLabelUpdateKind = "enable"
)

func PossibleValuesForRecommendedSensitivityLabelUpdateKind() []string {
	return []string{
		string(RecommendedSensitivityLabelUpdateKindDisable),
		string(RecommendedSensitivityLabelUpdateKindEnable),
	}
}

func (s *RecommendedSensitivityLabelUpdateKind) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRecommendedSensitivityLabelUpdateKind(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRecommendedSensitivityLabelUpdateKind(input string) (*RecommendedSensitivityLabelUpdateKind, error) {
	vals := map[string]RecommendedSensitivityLabelUpdateKind{
		"disable": RecommendedSensitivityLabelUpdateKindDisable,
		"enable":  RecommendedSensitivityLabelUpdateKindEnable,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RecommendedSensitivityLabelUpdateKind(input)
	return &out, nil
}

type SensitivityLabelRank string

const (
	SensitivityLabelRankCritical SensitivityLabelRank = "Critical"
	SensitivityLabelRankHigh     SensitivityLabelRank = "High"
	SensitivityLabelRankLow      SensitivityLabelRank = "Low"
	SensitivityLabelRankMedium   SensitivityLabelRank = "Medium"
	SensitivityLabelRankNone     SensitivityLabelRank = "None"
)

func PossibleValuesForSensitivityLabelRank() []string {
	return []string{
		string(SensitivityLabelRankCritical),
		string(SensitivityLabelRankHigh),
		string(SensitivityLabelRankLow),
		string(SensitivityLabelRankMedium),
		string(SensitivityLabelRankNone),
	}
}

func (s *SensitivityLabelRank) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSensitivityLabelRank(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSensitivityLabelRank(input string) (*SensitivityLabelRank, error) {
	vals := map[string]SensitivityLabelRank{
		"critical": SensitivityLabelRankCritical,
		"high":     SensitivityLabelRankHigh,
		"low":      SensitivityLabelRankLow,
		"medium":   SensitivityLabelRankMedium,
		"none":     SensitivityLabelRankNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SensitivityLabelRank(input)
	return &out, nil
}

type SensitivityLabelSource string

const (
	SensitivityLabelSourceCurrent     SensitivityLabelSource = "current"
	SensitivityLabelSourceRecommended SensitivityLabelSource = "recommended"
)

func PossibleValuesForSensitivityLabelSource() []string {
	return []string{
		string(SensitivityLabelSourceCurrent),
		string(SensitivityLabelSourceRecommended),
	}
}

func (s *SensitivityLabelSource) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSensitivityLabelSource(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSensitivityLabelSource(input string) (*SensitivityLabelSource, error) {
	vals := map[string]SensitivityLabelSource{
		"current":     SensitivityLabelSourceCurrent,
		"recommended": SensitivityLabelSourceRecommended,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SensitivityLabelSource(input)
	return &out, nil
}

type SensitivityLabelUpdateKind string

const (
	SensitivityLabelUpdateKindRemove SensitivityLabelUpdateKind = "remove"
	SensitivityLabelUpdateKindSet    SensitivityLabelUpdateKind = "set"
)

func PossibleValuesForSensitivityLabelUpdateKind() []string {
	return []string{
		string(SensitivityLabelUpdateKindRemove),
		string(SensitivityLabelUpdateKindSet),
	}
}

func (s *SensitivityLabelUpdateKind) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSensitivityLabelUpdateKind(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSensitivityLabelUpdateKind(input string) (*SensitivityLabelUpdateKind, error) {
	vals := map[string]SensitivityLabelUpdateKind{
		"remove": SensitivityLabelUpdateKindRemove,
		"set":    SensitivityLabelUpdateKindSet,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SensitivityLabelUpdateKind(input)
	return &out, nil
}
