package networksecurityperimeterconfiguration

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccessMode string

const (
	AccessModeAudit    AccessMode = "audit"
	AccessModeEnforced AccessMode = "enforced"
	AccessModeLearning AccessMode = "learning"
)

func PossibleValuesForAccessMode() []string {
	return []string{
		string(AccessModeAudit),
		string(AccessModeEnforced),
		string(AccessModeLearning),
	}
}

func (s *AccessMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAccessMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAccessMode(input string) (*AccessMode, error) {
	vals := map[string]AccessMode{
		"audit":    AccessModeAudit,
		"enforced": AccessModeEnforced,
		"learning": AccessModeLearning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AccessMode(input)
	return &out, nil
}

type AccessRuleDirection string

const (
	AccessRuleDirectionInbound  AccessRuleDirection = "Inbound"
	AccessRuleDirectionOutbound AccessRuleDirection = "Outbound"
)

func PossibleValuesForAccessRuleDirection() []string {
	return []string{
		string(AccessRuleDirectionInbound),
		string(AccessRuleDirectionOutbound),
	}
}

func (s *AccessRuleDirection) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAccessRuleDirection(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAccessRuleDirection(input string) (*AccessRuleDirection, error) {
	vals := map[string]AccessRuleDirection{
		"inbound":  AccessRuleDirectionInbound,
		"outbound": AccessRuleDirectionOutbound,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AccessRuleDirection(input)
	return &out, nil
}

type ProvisioningIssueSeverity string

const (
	ProvisioningIssueSeverityError   ProvisioningIssueSeverity = "Error"
	ProvisioningIssueSeverityWarning ProvisioningIssueSeverity = "Warning"
)

func PossibleValuesForProvisioningIssueSeverity() []string {
	return []string{
		string(ProvisioningIssueSeverityError),
		string(ProvisioningIssueSeverityWarning),
	}
}

func (s *ProvisioningIssueSeverity) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProvisioningIssueSeverity(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProvisioningIssueSeverity(input string) (*ProvisioningIssueSeverity, error) {
	vals := map[string]ProvisioningIssueSeverity{
		"error":   ProvisioningIssueSeverityError,
		"warning": ProvisioningIssueSeverityWarning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningIssueSeverity(input)
	return &out, nil
}

type ProvisioningIssueType string

const (
	ProvisioningIssueTypeConfigurationPropagationFailure ProvisioningIssueType = "ConfigurationPropagationFailure"
	ProvisioningIssueTypeMissingIdentityConfiguration    ProvisioningIssueType = "MissingIdentityConfiguration"
	ProvisioningIssueTypeMissingPerimeterConfiguration   ProvisioningIssueType = "MissingPerimeterConfiguration"
	ProvisioningIssueTypeOther                           ProvisioningIssueType = "Other"
)

func PossibleValuesForProvisioningIssueType() []string {
	return []string{
		string(ProvisioningIssueTypeConfigurationPropagationFailure),
		string(ProvisioningIssueTypeMissingIdentityConfiguration),
		string(ProvisioningIssueTypeMissingPerimeterConfiguration),
		string(ProvisioningIssueTypeOther),
	}
}

func (s *ProvisioningIssueType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProvisioningIssueType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProvisioningIssueType(input string) (*ProvisioningIssueType, error) {
	vals := map[string]ProvisioningIssueType{
		"configurationpropagationfailure": ProvisioningIssueTypeConfigurationPropagationFailure,
		"missingidentityconfiguration":    ProvisioningIssueTypeMissingIdentityConfiguration,
		"missingperimeterconfiguration":   ProvisioningIssueTypeMissingPerimeterConfiguration,
		"other":                           ProvisioningIssueTypeOther,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningIssueType(input)
	return &out, nil
}
