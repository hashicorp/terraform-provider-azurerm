package documents

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutocompleteMode string

const (
	AutocompleteModeOneTerm            AutocompleteMode = "oneTerm"
	AutocompleteModeOneTermWithContext AutocompleteMode = "oneTermWithContext"
	AutocompleteModeTwoTerms           AutocompleteMode = "twoTerms"
)

func PossibleValuesForAutocompleteMode() []string {
	return []string{
		string(AutocompleteModeOneTerm),
		string(AutocompleteModeOneTermWithContext),
		string(AutocompleteModeTwoTerms),
	}
}

func (s *AutocompleteMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAutocompleteMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAutocompleteMode(input string) (*AutocompleteMode, error) {
	vals := map[string]AutocompleteMode{
		"oneterm":            AutocompleteModeOneTerm,
		"onetermwithcontext": AutocompleteModeOneTermWithContext,
		"twoterms":           AutocompleteModeTwoTerms,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AutocompleteMode(input)
	return &out, nil
}

type IndexActionType string

const (
	IndexActionTypeDelete        IndexActionType = "delete"
	IndexActionTypeMerge         IndexActionType = "merge"
	IndexActionTypeMergeOrUpload IndexActionType = "mergeOrUpload"
	IndexActionTypeUpload        IndexActionType = "upload"
)

func PossibleValuesForIndexActionType() []string {
	return []string{
		string(IndexActionTypeDelete),
		string(IndexActionTypeMerge),
		string(IndexActionTypeMergeOrUpload),
		string(IndexActionTypeUpload),
	}
}

func (s *IndexActionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIndexActionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIndexActionType(input string) (*IndexActionType, error) {
	vals := map[string]IndexActionType{
		"delete":        IndexActionTypeDelete,
		"merge":         IndexActionTypeMerge,
		"mergeorupload": IndexActionTypeMergeOrUpload,
		"upload":        IndexActionTypeUpload,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IndexActionType(input)
	return &out, nil
}

type QueryAnswerType string

const (
	QueryAnswerTypeExtractive QueryAnswerType = "extractive"
	QueryAnswerTypeNone       QueryAnswerType = "none"
)

func PossibleValuesForQueryAnswerType() []string {
	return []string{
		string(QueryAnswerTypeExtractive),
		string(QueryAnswerTypeNone),
	}
}

func (s *QueryAnswerType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseQueryAnswerType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseQueryAnswerType(input string) (*QueryAnswerType, error) {
	vals := map[string]QueryAnswerType{
		"extractive": QueryAnswerTypeExtractive,
		"none":       QueryAnswerTypeNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := QueryAnswerType(input)
	return &out, nil
}

type QueryCaptionType string

const (
	QueryCaptionTypeExtractive QueryCaptionType = "extractive"
	QueryCaptionTypeNone       QueryCaptionType = "none"
)

func PossibleValuesForQueryCaptionType() []string {
	return []string{
		string(QueryCaptionTypeExtractive),
		string(QueryCaptionTypeNone),
	}
}

func (s *QueryCaptionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseQueryCaptionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseQueryCaptionType(input string) (*QueryCaptionType, error) {
	vals := map[string]QueryCaptionType{
		"extractive": QueryCaptionTypeExtractive,
		"none":       QueryCaptionTypeNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := QueryCaptionType(input)
	return &out, nil
}

type QueryDebugMode string

const (
	QueryDebugModeDisabled QueryDebugMode = "disabled"
	QueryDebugModeVector   QueryDebugMode = "vector"
)

func PossibleValuesForQueryDebugMode() []string {
	return []string{
		string(QueryDebugModeDisabled),
		string(QueryDebugModeVector),
	}
}

func (s *QueryDebugMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseQueryDebugMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseQueryDebugMode(input string) (*QueryDebugMode, error) {
	vals := map[string]QueryDebugMode{
		"disabled": QueryDebugModeDisabled,
		"vector":   QueryDebugModeVector,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := QueryDebugMode(input)
	return &out, nil
}

type QueryType string

const (
	QueryTypeFull     QueryType = "full"
	QueryTypeSemantic QueryType = "semantic"
	QueryTypeSimple   QueryType = "simple"
)

func PossibleValuesForQueryType() []string {
	return []string{
		string(QueryTypeFull),
		string(QueryTypeSemantic),
		string(QueryTypeSimple),
	}
}

func (s *QueryType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseQueryType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseQueryType(input string) (*QueryType, error) {
	vals := map[string]QueryType{
		"full":     QueryTypeFull,
		"semantic": QueryTypeSemantic,
		"simple":   QueryTypeSimple,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := QueryType(input)
	return &out, nil
}

type ScoringStatistics string

const (
	ScoringStatisticsGlobal ScoringStatistics = "global"
	ScoringStatisticsLocal  ScoringStatistics = "local"
)

func PossibleValuesForScoringStatistics() []string {
	return []string{
		string(ScoringStatisticsGlobal),
		string(ScoringStatisticsLocal),
	}
}

func (s *ScoringStatistics) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseScoringStatistics(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseScoringStatistics(input string) (*ScoringStatistics, error) {
	vals := map[string]ScoringStatistics{
		"global": ScoringStatisticsGlobal,
		"local":  ScoringStatisticsLocal,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScoringStatistics(input)
	return &out, nil
}

type SearchMode string

const (
	SearchModeAll SearchMode = "all"
	SearchModeAny SearchMode = "any"
)

func PossibleValuesForSearchMode() []string {
	return []string{
		string(SearchModeAll),
		string(SearchModeAny),
	}
}

func (s *SearchMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSearchMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSearchMode(input string) (*SearchMode, error) {
	vals := map[string]SearchMode{
		"all": SearchModeAll,
		"any": SearchModeAny,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SearchMode(input)
	return &out, nil
}

type SemanticErrorMode string

const (
	SemanticErrorModeFail    SemanticErrorMode = "fail"
	SemanticErrorModePartial SemanticErrorMode = "partial"
)

func PossibleValuesForSemanticErrorMode() []string {
	return []string{
		string(SemanticErrorModeFail),
		string(SemanticErrorModePartial),
	}
}

func (s *SemanticErrorMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSemanticErrorMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSemanticErrorMode(input string) (*SemanticErrorMode, error) {
	vals := map[string]SemanticErrorMode{
		"fail":    SemanticErrorModeFail,
		"partial": SemanticErrorModePartial,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SemanticErrorMode(input)
	return &out, nil
}

type SemanticErrorReason string

const (
	SemanticErrorReasonCapacityOverloaded SemanticErrorReason = "capacityOverloaded"
	SemanticErrorReasonMaxWaitExceeded    SemanticErrorReason = "maxWaitExceeded"
	SemanticErrorReasonTransient          SemanticErrorReason = "transient"
)

func PossibleValuesForSemanticErrorReason() []string {
	return []string{
		string(SemanticErrorReasonCapacityOverloaded),
		string(SemanticErrorReasonMaxWaitExceeded),
		string(SemanticErrorReasonTransient),
	}
}

func (s *SemanticErrorReason) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSemanticErrorReason(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSemanticErrorReason(input string) (*SemanticErrorReason, error) {
	vals := map[string]SemanticErrorReason{
		"capacityoverloaded": SemanticErrorReasonCapacityOverloaded,
		"maxwaitexceeded":    SemanticErrorReasonMaxWaitExceeded,
		"transient":          SemanticErrorReasonTransient,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SemanticErrorReason(input)
	return &out, nil
}

type SemanticSearchResultsType string

const (
	SemanticSearchResultsTypeBaseResults     SemanticSearchResultsType = "baseResults"
	SemanticSearchResultsTypeRerankedResults SemanticSearchResultsType = "rerankedResults"
)

func PossibleValuesForSemanticSearchResultsType() []string {
	return []string{
		string(SemanticSearchResultsTypeBaseResults),
		string(SemanticSearchResultsTypeRerankedResults),
	}
}

func (s *SemanticSearchResultsType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSemanticSearchResultsType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSemanticSearchResultsType(input string) (*SemanticSearchResultsType, error) {
	vals := map[string]SemanticSearchResultsType{
		"baseresults":     SemanticSearchResultsTypeBaseResults,
		"rerankedresults": SemanticSearchResultsTypeRerankedResults,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SemanticSearchResultsType(input)
	return &out, nil
}

type VectorFilterMode string

const (
	VectorFilterModePostFilter VectorFilterMode = "postFilter"
	VectorFilterModePreFilter  VectorFilterMode = "preFilter"
)

func PossibleValuesForVectorFilterMode() []string {
	return []string{
		string(VectorFilterModePostFilter),
		string(VectorFilterModePreFilter),
	}
}

func (s *VectorFilterMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVectorFilterMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVectorFilterMode(input string) (*VectorFilterMode, error) {
	vals := map[string]VectorFilterMode{
		"postfilter": VectorFilterModePostFilter,
		"prefilter":  VectorFilterModePreFilter,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VectorFilterMode(input)
	return &out, nil
}

type VectorQueryKind string

const (
	VectorQueryKindText   VectorQueryKind = "text"
	VectorQueryKindVector VectorQueryKind = "vector"
)

func PossibleValuesForVectorQueryKind() []string {
	return []string{
		string(VectorQueryKindText),
		string(VectorQueryKindVector),
	}
}

func (s *VectorQueryKind) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVectorQueryKind(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVectorQueryKind(input string) (*VectorQueryKind, error) {
	vals := map[string]VectorQueryKind{
		"text":   VectorQueryKindText,
		"vector": VectorQueryKindVector,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VectorQueryKind(input)
	return &out, nil
}
