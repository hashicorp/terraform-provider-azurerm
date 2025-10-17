package securitypartnerproviders

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProvisioningState string

const (
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUpdating  ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUpdating),
	}
}

func (s *ProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"deleting":  ProvisioningStateDeleting,
		"failed":    ProvisioningStateFailed,
		"succeeded": ProvisioningStateSucceeded,
		"updating":  ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type SecurityPartnerProviderConnectionStatus string

const (
	SecurityPartnerProviderConnectionStatusConnected          SecurityPartnerProviderConnectionStatus = "Connected"
	SecurityPartnerProviderConnectionStatusNotConnected       SecurityPartnerProviderConnectionStatus = "NotConnected"
	SecurityPartnerProviderConnectionStatusPartiallyConnected SecurityPartnerProviderConnectionStatus = "PartiallyConnected"
	SecurityPartnerProviderConnectionStatusUnknown            SecurityPartnerProviderConnectionStatus = "Unknown"
)

func PossibleValuesForSecurityPartnerProviderConnectionStatus() []string {
	return []string{
		string(SecurityPartnerProviderConnectionStatusConnected),
		string(SecurityPartnerProviderConnectionStatusNotConnected),
		string(SecurityPartnerProviderConnectionStatusPartiallyConnected),
		string(SecurityPartnerProviderConnectionStatusUnknown),
	}
}

func (s *SecurityPartnerProviderConnectionStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSecurityPartnerProviderConnectionStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSecurityPartnerProviderConnectionStatus(input string) (*SecurityPartnerProviderConnectionStatus, error) {
	vals := map[string]SecurityPartnerProviderConnectionStatus{
		"connected":          SecurityPartnerProviderConnectionStatusConnected,
		"notconnected":       SecurityPartnerProviderConnectionStatusNotConnected,
		"partiallyconnected": SecurityPartnerProviderConnectionStatusPartiallyConnected,
		"unknown":            SecurityPartnerProviderConnectionStatusUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SecurityPartnerProviderConnectionStatus(input)
	return &out, nil
}

type SecurityProviderName string

const (
	SecurityProviderNameCheckpoint SecurityProviderName = "Checkpoint"
	SecurityProviderNameIBoss      SecurityProviderName = "IBoss"
	SecurityProviderNameZScaler    SecurityProviderName = "ZScaler"
)

func PossibleValuesForSecurityProviderName() []string {
	return []string{
		string(SecurityProviderNameCheckpoint),
		string(SecurityProviderNameIBoss),
		string(SecurityProviderNameZScaler),
	}
}

func (s *SecurityProviderName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSecurityProviderName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSecurityProviderName(input string) (*SecurityProviderName, error) {
	vals := map[string]SecurityProviderName{
		"checkpoint": SecurityProviderNameCheckpoint,
		"iboss":      SecurityProviderNameIBoss,
		"zscaler":    SecurityProviderNameZScaler,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SecurityProviderName(input)
	return &out, nil
}
