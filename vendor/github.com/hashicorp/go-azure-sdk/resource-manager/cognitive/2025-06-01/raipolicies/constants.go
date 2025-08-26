package raipolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContentLevel string

const (
	ContentLevelHigh   ContentLevel = "High"
	ContentLevelLow    ContentLevel = "Low"
	ContentLevelMedium ContentLevel = "Medium"
)

func PossibleValuesForContentLevel() []string {
	return []string{
		string(ContentLevelHigh),
		string(ContentLevelLow),
		string(ContentLevelMedium),
	}
}

func (s *ContentLevel) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseContentLevel(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseContentLevel(input string) (*ContentLevel, error) {
	vals := map[string]ContentLevel{
		"high":   ContentLevelHigh,
		"low":    ContentLevelLow,
		"medium": ContentLevelMedium,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ContentLevel(input)
	return &out, nil
}

type RaiPolicyContentSource string

const (
	RaiPolicyContentSourceCompletion RaiPolicyContentSource = "Completion"
	RaiPolicyContentSourcePrompt     RaiPolicyContentSource = "Prompt"
)

func PossibleValuesForRaiPolicyContentSource() []string {
	return []string{
		string(RaiPolicyContentSourceCompletion),
		string(RaiPolicyContentSourcePrompt),
	}
}

func (s *RaiPolicyContentSource) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRaiPolicyContentSource(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRaiPolicyContentSource(input string) (*RaiPolicyContentSource, error) {
	vals := map[string]RaiPolicyContentSource{
		"completion": RaiPolicyContentSourceCompletion,
		"prompt":     RaiPolicyContentSourcePrompt,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RaiPolicyContentSource(input)
	return &out, nil
}

type RaiPolicyMode string

const (
	RaiPolicyModeAsynchronousFilter RaiPolicyMode = "Asynchronous_filter"
	RaiPolicyModeBlocking           RaiPolicyMode = "Blocking"
	RaiPolicyModeDefault            RaiPolicyMode = "Default"
	RaiPolicyModeDeferred           RaiPolicyMode = "Deferred"
)

func PossibleValuesForRaiPolicyMode() []string {
	return []string{
		string(RaiPolicyModeAsynchronousFilter),
		string(RaiPolicyModeBlocking),
		string(RaiPolicyModeDefault),
		string(RaiPolicyModeDeferred),
	}
}

func (s *RaiPolicyMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRaiPolicyMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRaiPolicyMode(input string) (*RaiPolicyMode, error) {
	vals := map[string]RaiPolicyMode{
		"asynchronous_filter": RaiPolicyModeAsynchronousFilter,
		"blocking":            RaiPolicyModeBlocking,
		"default":             RaiPolicyModeDefault,
		"deferred":            RaiPolicyModeDeferred,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RaiPolicyMode(input)
	return &out, nil
}

type RaiPolicyType string

const (
	RaiPolicyTypeSystemManaged RaiPolicyType = "SystemManaged"
	RaiPolicyTypeUserManaged   RaiPolicyType = "UserManaged"
)

func PossibleValuesForRaiPolicyType() []string {
	return []string{
		string(RaiPolicyTypeSystemManaged),
		string(RaiPolicyTypeUserManaged),
	}
}

func (s *RaiPolicyType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRaiPolicyType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRaiPolicyType(input string) (*RaiPolicyType, error) {
	vals := map[string]RaiPolicyType{
		"systemmanaged": RaiPolicyTypeSystemManaged,
		"usermanaged":   RaiPolicyTypeUserManaged,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RaiPolicyType(input)
	return &out, nil
}
