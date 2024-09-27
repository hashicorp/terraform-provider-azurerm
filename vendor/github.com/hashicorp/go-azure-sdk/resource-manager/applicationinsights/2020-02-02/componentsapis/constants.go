package componentsapis

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationType string

const (
	ApplicationTypeOther ApplicationType = "other"
	ApplicationTypeWeb   ApplicationType = "web"
)

func PossibleValuesForApplicationType() []string {
	return []string{
		string(ApplicationTypeOther),
		string(ApplicationTypeWeb),
	}
}

func (s *ApplicationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApplicationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApplicationType(input string) (*ApplicationType, error) {
	vals := map[string]ApplicationType{
		"other": ApplicationTypeOther,
		"web":   ApplicationTypeWeb,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApplicationType(input)
	return &out, nil
}

type FlowType string

const (
	FlowTypeBluefield FlowType = "Bluefield"
)

func PossibleValuesForFlowType() []string {
	return []string{
		string(FlowTypeBluefield),
	}
}

func (s *FlowType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFlowType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFlowType(input string) (*FlowType, error) {
	vals := map[string]FlowType{
		"bluefield": FlowTypeBluefield,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FlowType(input)
	return &out, nil
}

type IngestionMode string

const (
	IngestionModeApplicationInsights                       IngestionMode = "ApplicationInsights"
	IngestionModeApplicationInsightsWithDiagnosticSettings IngestionMode = "ApplicationInsightsWithDiagnosticSettings"
	IngestionModeLogAnalytics                              IngestionMode = "LogAnalytics"
)

func PossibleValuesForIngestionMode() []string {
	return []string{
		string(IngestionModeApplicationInsights),
		string(IngestionModeApplicationInsightsWithDiagnosticSettings),
		string(IngestionModeLogAnalytics),
	}
}

func (s *IngestionMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIngestionMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIngestionMode(input string) (*IngestionMode, error) {
	vals := map[string]IngestionMode{
		"applicationinsights":                       IngestionModeApplicationInsights,
		"applicationinsightswithdiagnosticsettings": IngestionModeApplicationInsightsWithDiagnosticSettings,
		"loganalytics":                              IngestionModeLogAnalytics,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IngestionMode(input)
	return &out, nil
}

type PublicNetworkAccessType string

const (
	PublicNetworkAccessTypeDisabled PublicNetworkAccessType = "Disabled"
	PublicNetworkAccessTypeEnabled  PublicNetworkAccessType = "Enabled"
)

func PossibleValuesForPublicNetworkAccessType() []string {
	return []string{
		string(PublicNetworkAccessTypeDisabled),
		string(PublicNetworkAccessTypeEnabled),
	}
}

func (s *PublicNetworkAccessType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePublicNetworkAccessType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePublicNetworkAccessType(input string) (*PublicNetworkAccessType, error) {
	vals := map[string]PublicNetworkAccessType{
		"disabled": PublicNetworkAccessTypeDisabled,
		"enabled":  PublicNetworkAccessTypeEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PublicNetworkAccessType(input)
	return &out, nil
}

type PurgeState string

const (
	PurgeStateCompleted PurgeState = "completed"
	PurgeStatePending   PurgeState = "pending"
)

func PossibleValuesForPurgeState() []string {
	return []string{
		string(PurgeStateCompleted),
		string(PurgeStatePending),
	}
}

func (s *PurgeState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePurgeState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePurgeState(input string) (*PurgeState, error) {
	vals := map[string]PurgeState{
		"completed": PurgeStateCompleted,
		"pending":   PurgeStatePending,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PurgeState(input)
	return &out, nil
}

type RequestSource string

const (
	RequestSourceRest RequestSource = "rest"
)

func PossibleValuesForRequestSource() []string {
	return []string{
		string(RequestSourceRest),
	}
}

func (s *RequestSource) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRequestSource(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRequestSource(input string) (*RequestSource, error) {
	vals := map[string]RequestSource{
		"rest": RequestSourceRest,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RequestSource(input)
	return &out, nil
}
