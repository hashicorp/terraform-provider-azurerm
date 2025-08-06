package replicationprotectioncontainermappings

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AgentAutoUpdateStatus string

const (
	AgentAutoUpdateStatusDisabled AgentAutoUpdateStatus = "Disabled"
	AgentAutoUpdateStatusEnabled  AgentAutoUpdateStatus = "Enabled"
)

func PossibleValuesForAgentAutoUpdateStatus() []string {
	return []string{
		string(AgentAutoUpdateStatusDisabled),
		string(AgentAutoUpdateStatusEnabled),
	}
}

func (s *AgentAutoUpdateStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAgentAutoUpdateStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAgentAutoUpdateStatus(input string) (*AgentAutoUpdateStatus, error) {
	vals := map[string]AgentAutoUpdateStatus{
		"disabled": AgentAutoUpdateStatusDisabled,
		"enabled":  AgentAutoUpdateStatusEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AgentAutoUpdateStatus(input)
	return &out, nil
}

type AutomationAccountAuthenticationType string

const (
	AutomationAccountAuthenticationTypeRunAsAccount           AutomationAccountAuthenticationType = "RunAsAccount"
	AutomationAccountAuthenticationTypeSystemAssignedIdentity AutomationAccountAuthenticationType = "SystemAssignedIdentity"
)

func PossibleValuesForAutomationAccountAuthenticationType() []string {
	return []string{
		string(AutomationAccountAuthenticationTypeRunAsAccount),
		string(AutomationAccountAuthenticationTypeSystemAssignedIdentity),
	}
}

func (s *AutomationAccountAuthenticationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAutomationAccountAuthenticationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAutomationAccountAuthenticationType(input string) (*AutomationAccountAuthenticationType, error) {
	vals := map[string]AutomationAccountAuthenticationType{
		"runasaccount":           AutomationAccountAuthenticationTypeRunAsAccount,
		"systemassignedidentity": AutomationAccountAuthenticationTypeSystemAssignedIdentity,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AutomationAccountAuthenticationType(input)
	return &out, nil
}

type HealthErrorCustomerResolvability string

const (
	HealthErrorCustomerResolvabilityAllowed    HealthErrorCustomerResolvability = "Allowed"
	HealthErrorCustomerResolvabilityNotAllowed HealthErrorCustomerResolvability = "NotAllowed"
)

func PossibleValuesForHealthErrorCustomerResolvability() []string {
	return []string{
		string(HealthErrorCustomerResolvabilityAllowed),
		string(HealthErrorCustomerResolvabilityNotAllowed),
	}
}

func (s *HealthErrorCustomerResolvability) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHealthErrorCustomerResolvability(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHealthErrorCustomerResolvability(input string) (*HealthErrorCustomerResolvability, error) {
	vals := map[string]HealthErrorCustomerResolvability{
		"allowed":    HealthErrorCustomerResolvabilityAllowed,
		"notallowed": HealthErrorCustomerResolvabilityNotAllowed,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HealthErrorCustomerResolvability(input)
	return &out, nil
}
