package resources

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthorizationScopeFilter string

const (
	AuthorizationScopeFilterAtScopeAboveAndBelow AuthorizationScopeFilter = "AtScopeAboveAndBelow"
	AuthorizationScopeFilterAtScopeAndAbove      AuthorizationScopeFilter = "AtScopeAndAbove"
	AuthorizationScopeFilterAtScopeAndBelow      AuthorizationScopeFilter = "AtScopeAndBelow"
	AuthorizationScopeFilterAtScopeExact         AuthorizationScopeFilter = "AtScopeExact"
)

func PossibleValuesForAuthorizationScopeFilter() []string {
	return []string{
		string(AuthorizationScopeFilterAtScopeAboveAndBelow),
		string(AuthorizationScopeFilterAtScopeAndAbove),
		string(AuthorizationScopeFilterAtScopeAndBelow),
		string(AuthorizationScopeFilterAtScopeExact),
	}
}

func (s *AuthorizationScopeFilter) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAuthorizationScopeFilter(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAuthorizationScopeFilter(input string) (*AuthorizationScopeFilter, error) {
	vals := map[string]AuthorizationScopeFilter{
		"atscopeaboveandbelow": AuthorizationScopeFilterAtScopeAboveAndBelow,
		"atscopeandabove":      AuthorizationScopeFilterAtScopeAndAbove,
		"atscopeandbelow":      AuthorizationScopeFilterAtScopeAndBelow,
		"atscopeexact":         AuthorizationScopeFilterAtScopeExact,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AuthorizationScopeFilter(input)
	return &out, nil
}

type FacetSortOrder string

const (
	FacetSortOrderAsc  FacetSortOrder = "asc"
	FacetSortOrderDesc FacetSortOrder = "desc"
)

func PossibleValuesForFacetSortOrder() []string {
	return []string{
		string(FacetSortOrderAsc),
		string(FacetSortOrderDesc),
	}
}

func (s *FacetSortOrder) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFacetSortOrder(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFacetSortOrder(input string) (*FacetSortOrder, error) {
	vals := map[string]FacetSortOrder{
		"asc":  FacetSortOrderAsc,
		"desc": FacetSortOrderDesc,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FacetSortOrder(input)
	return &out, nil
}

type ResultFormat string

const (
	ResultFormatObjectArray ResultFormat = "objectArray"
	ResultFormatTable       ResultFormat = "table"
)

func PossibleValuesForResultFormat() []string {
	return []string{
		string(ResultFormatObjectArray),
		string(ResultFormatTable),
	}
}

func (s *ResultFormat) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseResultFormat(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseResultFormat(input string) (*ResultFormat, error) {
	vals := map[string]ResultFormat{
		"objectarray": ResultFormatObjectArray,
		"table":       ResultFormatTable,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResultFormat(input)
	return &out, nil
}

type ResultTruncated string

const (
	ResultTruncatedFalse ResultTruncated = "false"
	ResultTruncatedTrue  ResultTruncated = "true"
)

func PossibleValuesForResultTruncated() []string {
	return []string{
		string(ResultTruncatedFalse),
		string(ResultTruncatedTrue),
	}
}

func (s *ResultTruncated) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseResultTruncated(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseResultTruncated(input string) (*ResultTruncated, error) {
	vals := map[string]ResultTruncated{
		"false": ResultTruncatedFalse,
		"true":  ResultTruncatedTrue,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResultTruncated(input)
	return &out, nil
}
