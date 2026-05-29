package smartdetectoralertrules

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AlertRuleState string

const (
	AlertRuleStateDisabled AlertRuleState = "Disabled"
	AlertRuleStateEnabled  AlertRuleState = "Enabled"
)

func PossibleValuesForAlertRuleState() []string {
	return []string{
		string(AlertRuleStateDisabled),
		string(AlertRuleStateEnabled),
	}
}

func (s *AlertRuleState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAlertRuleState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAlertRuleState(input string) (*AlertRuleState, error) {
	vals := map[string]AlertRuleState{
		"disabled": AlertRuleStateDisabled,
		"enabled":  AlertRuleStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AlertRuleState(input)
	return &out, nil
}

type Severity string

const (
	SeveritySevFour  Severity = "Sev4"
	SeveritySevOne   Severity = "Sev1"
	SeveritySevThree Severity = "Sev3"
	SeveritySevTwo   Severity = "Sev2"
	SeveritySevZero  Severity = "Sev0"
)

func PossibleValuesForSeverity() []string {
	return []string{
		string(SeveritySevFour),
		string(SeveritySevOne),
		string(SeveritySevThree),
		string(SeveritySevTwo),
		string(SeveritySevZero),
	}
}

func (s *Severity) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSeverity(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSeverity(input string) (*Severity, error) {
	vals := map[string]Severity{
		"sev4": SeveritySevFour,
		"sev1": SeveritySevOne,
		"sev3": SeveritySevThree,
		"sev2": SeveritySevTwo,
		"sev0": SeveritySevZero,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Severity(input)
	return &out, nil
}
