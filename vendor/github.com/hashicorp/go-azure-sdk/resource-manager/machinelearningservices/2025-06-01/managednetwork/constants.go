package managednetwork

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedNetworkStatus string

const (
	ManagedNetworkStatusActive   ManagedNetworkStatus = "Active"
	ManagedNetworkStatusInactive ManagedNetworkStatus = "Inactive"
)

func PossibleValuesForManagedNetworkStatus() []string {
	return []string{
		string(ManagedNetworkStatusActive),
		string(ManagedNetworkStatusInactive),
	}
}

func (s *ManagedNetworkStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseManagedNetworkStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseManagedNetworkStatus(input string) (*ManagedNetworkStatus, error) {
	vals := map[string]ManagedNetworkStatus{
		"active":   ManagedNetworkStatusActive,
		"inactive": ManagedNetworkStatusInactive,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ManagedNetworkStatus(input)
	return &out, nil
}

type RuleAction string

const (
	RuleActionAllow RuleAction = "Allow"
	RuleActionDeny  RuleAction = "Deny"
)

func PossibleValuesForRuleAction() []string {
	return []string{
		string(RuleActionAllow),
		string(RuleActionDeny),
	}
}

func (s *RuleAction) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRuleAction(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRuleAction(input string) (*RuleAction, error) {
	vals := map[string]RuleAction{
		"allow": RuleActionAllow,
		"deny":  RuleActionDeny,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RuleAction(input)
	return &out, nil
}

type RuleCategory string

const (
	RuleCategoryDependency  RuleCategory = "Dependency"
	RuleCategoryRecommended RuleCategory = "Recommended"
	RuleCategoryRequired    RuleCategory = "Required"
	RuleCategoryUserDefined RuleCategory = "UserDefined"
)

func PossibleValuesForRuleCategory() []string {
	return []string{
		string(RuleCategoryDependency),
		string(RuleCategoryRecommended),
		string(RuleCategoryRequired),
		string(RuleCategoryUserDefined),
	}
}

func (s *RuleCategory) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRuleCategory(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRuleCategory(input string) (*RuleCategory, error) {
	vals := map[string]RuleCategory{
		"dependency":  RuleCategoryDependency,
		"recommended": RuleCategoryRecommended,
		"required":    RuleCategoryRequired,
		"userdefined": RuleCategoryUserDefined,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RuleCategory(input)
	return &out, nil
}

type RuleStatus string

const (
	RuleStatusActive       RuleStatus = "Active"
	RuleStatusDeleting     RuleStatus = "Deleting"
	RuleStatusFailed       RuleStatus = "Failed"
	RuleStatusInactive     RuleStatus = "Inactive"
	RuleStatusProvisioning RuleStatus = "Provisioning"
)

func PossibleValuesForRuleStatus() []string {
	return []string{
		string(RuleStatusActive),
		string(RuleStatusDeleting),
		string(RuleStatusFailed),
		string(RuleStatusInactive),
		string(RuleStatusProvisioning),
	}
}

func (s *RuleStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRuleStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRuleStatus(input string) (*RuleStatus, error) {
	vals := map[string]RuleStatus{
		"active":       RuleStatusActive,
		"deleting":     RuleStatusDeleting,
		"failed":       RuleStatusFailed,
		"inactive":     RuleStatusInactive,
		"provisioning": RuleStatusProvisioning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RuleStatus(input)
	return &out, nil
}

type RuleType string

const (
	RuleTypeFQDN            RuleType = "FQDN"
	RuleTypePrivateEndpoint RuleType = "PrivateEndpoint"
	RuleTypeServiceTag      RuleType = "ServiceTag"
)

func PossibleValuesForRuleType() []string {
	return []string{
		string(RuleTypeFQDN),
		string(RuleTypePrivateEndpoint),
		string(RuleTypeServiceTag),
	}
}

func (s *RuleType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRuleType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRuleType(input string) (*RuleType, error) {
	vals := map[string]RuleType{
		"fqdn":            RuleTypeFQDN,
		"privateendpoint": RuleTypePrivateEndpoint,
		"servicetag":      RuleTypeServiceTag,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RuleType(input)
	return &out, nil
}
