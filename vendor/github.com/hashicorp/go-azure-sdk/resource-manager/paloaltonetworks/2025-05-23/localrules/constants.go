package localrules

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActionEnum string

const (
	ActionEnumAllow           ActionEnum = "Allow"
	ActionEnumDenyResetBoth   ActionEnum = "DenyResetBoth"
	ActionEnumDenyResetServer ActionEnum = "DenyResetServer"
	ActionEnumDenySilent      ActionEnum = "DenySilent"
)

func PossibleValuesForActionEnum() []string {
	return []string{
		string(ActionEnumAllow),
		string(ActionEnumDenyResetBoth),
		string(ActionEnumDenyResetServer),
		string(ActionEnumDenySilent),
	}
}

func (s *ActionEnum) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseActionEnum(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseActionEnum(input string) (*ActionEnum, error) {
	vals := map[string]ActionEnum{
		"allow":           ActionEnumAllow,
		"denyresetboth":   ActionEnumDenyResetBoth,
		"denyresetserver": ActionEnumDenyResetServer,
		"denysilent":      ActionEnumDenySilent,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ActionEnum(input)
	return &out, nil
}

type BooleanEnum string

const (
	BooleanEnumFALSE BooleanEnum = "FALSE"
	BooleanEnumTRUE  BooleanEnum = "TRUE"
)

func PossibleValuesForBooleanEnum() []string {
	return []string{
		string(BooleanEnumFALSE),
		string(BooleanEnumTRUE),
	}
}

func (s *BooleanEnum) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBooleanEnum(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBooleanEnum(input string) (*BooleanEnum, error) {
	vals := map[string]BooleanEnum{
		"false": BooleanEnumFALSE,
		"true":  BooleanEnumTRUE,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BooleanEnum(input)
	return &out, nil
}

type DecryptionRuleTypeEnum string

const (
	DecryptionRuleTypeEnumNone                  DecryptionRuleTypeEnum = "None"
	DecryptionRuleTypeEnumSSLInboundInspection  DecryptionRuleTypeEnum = "SSLInboundInspection"
	DecryptionRuleTypeEnumSSLOutboundInspection DecryptionRuleTypeEnum = "SSLOutboundInspection"
)

func PossibleValuesForDecryptionRuleTypeEnum() []string {
	return []string{
		string(DecryptionRuleTypeEnumNone),
		string(DecryptionRuleTypeEnumSSLInboundInspection),
		string(DecryptionRuleTypeEnumSSLOutboundInspection),
	}
}

func (s *DecryptionRuleTypeEnum) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDecryptionRuleTypeEnum(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDecryptionRuleTypeEnum(input string) (*DecryptionRuleTypeEnum, error) {
	vals := map[string]DecryptionRuleTypeEnum{
		"none":                  DecryptionRuleTypeEnumNone,
		"sslinboundinspection":  DecryptionRuleTypeEnumSSLInboundInspection,
		"ssloutboundinspection": DecryptionRuleTypeEnumSSLOutboundInspection,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DecryptionRuleTypeEnum(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateAccepted     ProvisioningState = "Accepted"
	ProvisioningStateCanceled     ProvisioningState = "Canceled"
	ProvisioningStateCreating     ProvisioningState = "Creating"
	ProvisioningStateDeleted      ProvisioningState = "Deleted"
	ProvisioningStateDeleting     ProvisioningState = "Deleting"
	ProvisioningStateFailed       ProvisioningState = "Failed"
	ProvisioningStateNotSpecified ProvisioningState = "NotSpecified"
	ProvisioningStateSucceeded    ProvisioningState = "Succeeded"
	ProvisioningStateUpdating     ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleted),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateNotSpecified),
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
		"accepted":     ProvisioningStateAccepted,
		"canceled":     ProvisioningStateCanceled,
		"creating":     ProvisioningStateCreating,
		"deleted":      ProvisioningStateDeleted,
		"deleting":     ProvisioningStateDeleting,
		"failed":       ProvisioningStateFailed,
		"notspecified": ProvisioningStateNotSpecified,
		"succeeded":    ProvisioningStateSucceeded,
		"updating":     ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type StateEnum string

const (
	StateEnumDISABLED StateEnum = "DISABLED"
	StateEnumENABLED  StateEnum = "ENABLED"
)

func PossibleValuesForStateEnum() []string {
	return []string{
		string(StateEnumDISABLED),
		string(StateEnumENABLED),
	}
}

func (s *StateEnum) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStateEnum(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStateEnum(input string) (*StateEnum, error) {
	vals := map[string]StateEnum{
		"disabled": StateEnumDISABLED,
		"enabled":  StateEnumENABLED,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StateEnum(input)
	return &out, nil
}
