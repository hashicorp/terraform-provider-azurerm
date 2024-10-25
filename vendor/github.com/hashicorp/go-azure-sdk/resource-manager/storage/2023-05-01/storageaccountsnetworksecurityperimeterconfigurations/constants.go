package storageaccountsnetworksecurityperimeterconfigurations

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IssueType string

const (
	IssueTypeConfigurationPropagationFailure IssueType = "ConfigurationPropagationFailure"
	IssueTypeUnknown                         IssueType = "Unknown"
)

func PossibleValuesForIssueType() []string {
	return []string{
		string(IssueTypeConfigurationPropagationFailure),
		string(IssueTypeUnknown),
	}
}

func (s *IssueType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIssueType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIssueType(input string) (*IssueType, error) {
	vals := map[string]IssueType{
		"configurationpropagationfailure": IssueTypeConfigurationPropagationFailure,
		"unknown":                         IssueTypeUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IssueType(input)
	return &out, nil
}

type NetworkSecurityPerimeterConfigurationProvisioningState string

const (
	NetworkSecurityPerimeterConfigurationProvisioningStateAccepted  NetworkSecurityPerimeterConfigurationProvisioningState = "Accepted"
	NetworkSecurityPerimeterConfigurationProvisioningStateCanceled  NetworkSecurityPerimeterConfigurationProvisioningState = "Canceled"
	NetworkSecurityPerimeterConfigurationProvisioningStateDeleting  NetworkSecurityPerimeterConfigurationProvisioningState = "Deleting"
	NetworkSecurityPerimeterConfigurationProvisioningStateFailed    NetworkSecurityPerimeterConfigurationProvisioningState = "Failed"
	NetworkSecurityPerimeterConfigurationProvisioningStateSucceeded NetworkSecurityPerimeterConfigurationProvisioningState = "Succeeded"
)

func PossibleValuesForNetworkSecurityPerimeterConfigurationProvisioningState() []string {
	return []string{
		string(NetworkSecurityPerimeterConfigurationProvisioningStateAccepted),
		string(NetworkSecurityPerimeterConfigurationProvisioningStateCanceled),
		string(NetworkSecurityPerimeterConfigurationProvisioningStateDeleting),
		string(NetworkSecurityPerimeterConfigurationProvisioningStateFailed),
		string(NetworkSecurityPerimeterConfigurationProvisioningStateSucceeded),
	}
}

func (s *NetworkSecurityPerimeterConfigurationProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNetworkSecurityPerimeterConfigurationProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNetworkSecurityPerimeterConfigurationProvisioningState(input string) (*NetworkSecurityPerimeterConfigurationProvisioningState, error) {
	vals := map[string]NetworkSecurityPerimeterConfigurationProvisioningState{
		"accepted":  NetworkSecurityPerimeterConfigurationProvisioningStateAccepted,
		"canceled":  NetworkSecurityPerimeterConfigurationProvisioningStateCanceled,
		"deleting":  NetworkSecurityPerimeterConfigurationProvisioningStateDeleting,
		"failed":    NetworkSecurityPerimeterConfigurationProvisioningStateFailed,
		"succeeded": NetworkSecurityPerimeterConfigurationProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NetworkSecurityPerimeterConfigurationProvisioningState(input)
	return &out, nil
}

type NspAccessRuleDirection string

const (
	NspAccessRuleDirectionInbound  NspAccessRuleDirection = "Inbound"
	NspAccessRuleDirectionOutbound NspAccessRuleDirection = "Outbound"
)

func PossibleValuesForNspAccessRuleDirection() []string {
	return []string{
		string(NspAccessRuleDirectionInbound),
		string(NspAccessRuleDirectionOutbound),
	}
}

func (s *NspAccessRuleDirection) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNspAccessRuleDirection(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNspAccessRuleDirection(input string) (*NspAccessRuleDirection, error) {
	vals := map[string]NspAccessRuleDirection{
		"inbound":  NspAccessRuleDirectionInbound,
		"outbound": NspAccessRuleDirectionOutbound,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NspAccessRuleDirection(input)
	return &out, nil
}

type ResourceAssociationAccessMode string

const (
	ResourceAssociationAccessModeAudit    ResourceAssociationAccessMode = "Audit"
	ResourceAssociationAccessModeEnforced ResourceAssociationAccessMode = "Enforced"
	ResourceAssociationAccessModeLearning ResourceAssociationAccessMode = "Learning"
)

func PossibleValuesForResourceAssociationAccessMode() []string {
	return []string{
		string(ResourceAssociationAccessModeAudit),
		string(ResourceAssociationAccessModeEnforced),
		string(ResourceAssociationAccessModeLearning),
	}
}

func (s *ResourceAssociationAccessMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseResourceAssociationAccessMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseResourceAssociationAccessMode(input string) (*ResourceAssociationAccessMode, error) {
	vals := map[string]ResourceAssociationAccessMode{
		"audit":    ResourceAssociationAccessModeAudit,
		"enforced": ResourceAssociationAccessModeEnforced,
		"learning": ResourceAssociationAccessModeLearning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceAssociationAccessMode(input)
	return &out, nil
}

type Severity string

const (
	SeverityError   Severity = "Error"
	SeverityWarning Severity = "Warning"
)

func PossibleValuesForSeverity() []string {
	return []string{
		string(SeverityError),
		string(SeverityWarning),
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
		"error":   SeverityError,
		"warning": SeverityWarning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Severity(input)
	return &out, nil
}
