package applicationgatewaywafdynamicmanifests

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayRuleSetStatusOptions string

const (
	ApplicationGatewayRuleSetStatusOptionsDeprecated ApplicationGatewayRuleSetStatusOptions = "Deprecated"
	ApplicationGatewayRuleSetStatusOptionsGA         ApplicationGatewayRuleSetStatusOptions = "GA"
	ApplicationGatewayRuleSetStatusOptionsPreview    ApplicationGatewayRuleSetStatusOptions = "Preview"
	ApplicationGatewayRuleSetStatusOptionsSupported  ApplicationGatewayRuleSetStatusOptions = "Supported"
)

func PossibleValuesForApplicationGatewayRuleSetStatusOptions() []string {
	return []string{
		string(ApplicationGatewayRuleSetStatusOptionsDeprecated),
		string(ApplicationGatewayRuleSetStatusOptionsGA),
		string(ApplicationGatewayRuleSetStatusOptionsPreview),
		string(ApplicationGatewayRuleSetStatusOptionsSupported),
	}
}

func (s *ApplicationGatewayRuleSetStatusOptions) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApplicationGatewayRuleSetStatusOptions(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApplicationGatewayRuleSetStatusOptions(input string) (*ApplicationGatewayRuleSetStatusOptions, error) {
	vals := map[string]ApplicationGatewayRuleSetStatusOptions{
		"deprecated": ApplicationGatewayRuleSetStatusOptionsDeprecated,
		"ga":         ApplicationGatewayRuleSetStatusOptionsGA,
		"preview":    ApplicationGatewayRuleSetStatusOptionsPreview,
		"supported":  ApplicationGatewayRuleSetStatusOptionsSupported,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApplicationGatewayRuleSetStatusOptions(input)
	return &out, nil
}

type ApplicationGatewayTierTypes string

const (
	ApplicationGatewayTierTypesStandard     ApplicationGatewayTierTypes = "Standard"
	ApplicationGatewayTierTypesStandardVTwo ApplicationGatewayTierTypes = "Standard_v2"
	ApplicationGatewayTierTypesWAF          ApplicationGatewayTierTypes = "WAF"
	ApplicationGatewayTierTypesWAFVTwo      ApplicationGatewayTierTypes = "WAF_v2"
)

func PossibleValuesForApplicationGatewayTierTypes() []string {
	return []string{
		string(ApplicationGatewayTierTypesStandard),
		string(ApplicationGatewayTierTypesStandardVTwo),
		string(ApplicationGatewayTierTypesWAF),
		string(ApplicationGatewayTierTypesWAFVTwo),
	}
}

func (s *ApplicationGatewayTierTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApplicationGatewayTierTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApplicationGatewayTierTypes(input string) (*ApplicationGatewayTierTypes, error) {
	vals := map[string]ApplicationGatewayTierTypes{
		"standard":    ApplicationGatewayTierTypesStandard,
		"standard_v2": ApplicationGatewayTierTypesStandardVTwo,
		"waf":         ApplicationGatewayTierTypesWAF,
		"waf_v2":      ApplicationGatewayTierTypesWAFVTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApplicationGatewayTierTypes(input)
	return &out, nil
}

type ApplicationGatewayWafRuleActionTypes string

const (
	ApplicationGatewayWafRuleActionTypesAllow          ApplicationGatewayWafRuleActionTypes = "Allow"
	ApplicationGatewayWafRuleActionTypesAnomalyScoring ApplicationGatewayWafRuleActionTypes = "AnomalyScoring"
	ApplicationGatewayWafRuleActionTypesBlock          ApplicationGatewayWafRuleActionTypes = "Block"
	ApplicationGatewayWafRuleActionTypesLog            ApplicationGatewayWafRuleActionTypes = "Log"
	ApplicationGatewayWafRuleActionTypesNone           ApplicationGatewayWafRuleActionTypes = "None"
)

func PossibleValuesForApplicationGatewayWafRuleActionTypes() []string {
	return []string{
		string(ApplicationGatewayWafRuleActionTypesAllow),
		string(ApplicationGatewayWafRuleActionTypesAnomalyScoring),
		string(ApplicationGatewayWafRuleActionTypesBlock),
		string(ApplicationGatewayWafRuleActionTypesLog),
		string(ApplicationGatewayWafRuleActionTypesNone),
	}
}

func (s *ApplicationGatewayWafRuleActionTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApplicationGatewayWafRuleActionTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApplicationGatewayWafRuleActionTypes(input string) (*ApplicationGatewayWafRuleActionTypes, error) {
	vals := map[string]ApplicationGatewayWafRuleActionTypes{
		"allow":          ApplicationGatewayWafRuleActionTypesAllow,
		"anomalyscoring": ApplicationGatewayWafRuleActionTypesAnomalyScoring,
		"block":          ApplicationGatewayWafRuleActionTypesBlock,
		"log":            ApplicationGatewayWafRuleActionTypesLog,
		"none":           ApplicationGatewayWafRuleActionTypesNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApplicationGatewayWafRuleActionTypes(input)
	return &out, nil
}

type ApplicationGatewayWafRuleStateTypes string

const (
	ApplicationGatewayWafRuleStateTypesDisabled ApplicationGatewayWafRuleStateTypes = "Disabled"
	ApplicationGatewayWafRuleStateTypesEnabled  ApplicationGatewayWafRuleStateTypes = "Enabled"
)

func PossibleValuesForApplicationGatewayWafRuleStateTypes() []string {
	return []string{
		string(ApplicationGatewayWafRuleStateTypesDisabled),
		string(ApplicationGatewayWafRuleStateTypesEnabled),
	}
}

func (s *ApplicationGatewayWafRuleStateTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApplicationGatewayWafRuleStateTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApplicationGatewayWafRuleStateTypes(input string) (*ApplicationGatewayWafRuleStateTypes, error) {
	vals := map[string]ApplicationGatewayWafRuleStateTypes{
		"disabled": ApplicationGatewayWafRuleStateTypesDisabled,
		"enabled":  ApplicationGatewayWafRuleStateTypesEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApplicationGatewayWafRuleStateTypes(input)
	return &out, nil
}
