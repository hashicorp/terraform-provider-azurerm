package externalauths

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConditionType string

const (
	ConditionTypeAvailable   ConditionType = "Available"
	ConditionTypeDegraded    ConditionType = "Degraded"
	ConditionTypeProgressing ConditionType = "Progressing"
)

func PossibleValuesForConditionType() []string {
	return []string{
		string(ConditionTypeAvailable),
		string(ConditionTypeDegraded),
		string(ConditionTypeProgressing),
	}
}

func (s *ConditionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConditionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConditionType(input string) (*ConditionType, error) {
	vals := map[string]ConditionType{
		"available":   ConditionTypeAvailable,
		"degraded":    ConditionTypeDegraded,
		"progressing": ConditionTypeProgressing,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConditionType(input)
	return &out, nil
}

type ExternalAuthClientType string

const (
	ExternalAuthClientTypeConfidential ExternalAuthClientType = "Confidential"
	ExternalAuthClientTypePublic       ExternalAuthClientType = "Public"
)

func PossibleValuesForExternalAuthClientType() []string {
	return []string{
		string(ExternalAuthClientTypeConfidential),
		string(ExternalAuthClientTypePublic),
	}
}

func (s *ExternalAuthClientType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseExternalAuthClientType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseExternalAuthClientType(input string) (*ExternalAuthClientType, error) {
	vals := map[string]ExternalAuthClientType{
		"confidential": ExternalAuthClientTypeConfidential,
		"public":       ExternalAuthClientTypePublic,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExternalAuthClientType(input)
	return &out, nil
}

type ExternalAuthProvisioningState string

const (
	ExternalAuthProvisioningStateAccepted       ExternalAuthProvisioningState = "Accepted"
	ExternalAuthProvisioningStateAwaitingSecret ExternalAuthProvisioningState = "AwaitingSecret"
	ExternalAuthProvisioningStateCanceled       ExternalAuthProvisioningState = "Canceled"
	ExternalAuthProvisioningStateDeleting       ExternalAuthProvisioningState = "Deleting"
	ExternalAuthProvisioningStateFailed         ExternalAuthProvisioningState = "Failed"
	ExternalAuthProvisioningStateProvisioning   ExternalAuthProvisioningState = "Provisioning"
	ExternalAuthProvisioningStateSucceeded      ExternalAuthProvisioningState = "Succeeded"
	ExternalAuthProvisioningStateUpdating       ExternalAuthProvisioningState = "Updating"
)

func PossibleValuesForExternalAuthProvisioningState() []string {
	return []string{
		string(ExternalAuthProvisioningStateAccepted),
		string(ExternalAuthProvisioningStateAwaitingSecret),
		string(ExternalAuthProvisioningStateCanceled),
		string(ExternalAuthProvisioningStateDeleting),
		string(ExternalAuthProvisioningStateFailed),
		string(ExternalAuthProvisioningStateProvisioning),
		string(ExternalAuthProvisioningStateSucceeded),
		string(ExternalAuthProvisioningStateUpdating),
	}
}

func (s *ExternalAuthProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseExternalAuthProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseExternalAuthProvisioningState(input string) (*ExternalAuthProvisioningState, error) {
	vals := map[string]ExternalAuthProvisioningState{
		"accepted":       ExternalAuthProvisioningStateAccepted,
		"awaitingsecret": ExternalAuthProvisioningStateAwaitingSecret,
		"canceled":       ExternalAuthProvisioningStateCanceled,
		"deleting":       ExternalAuthProvisioningStateDeleting,
		"failed":         ExternalAuthProvisioningStateFailed,
		"provisioning":   ExternalAuthProvisioningStateProvisioning,
		"succeeded":      ExternalAuthProvisioningStateSucceeded,
		"updating":       ExternalAuthProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExternalAuthProvisioningState(input)
	return &out, nil
}

type StatusType string

const (
	StatusTypeFalse   StatusType = "False"
	StatusTypeTrue    StatusType = "True"
	StatusTypeUnknown StatusType = "Unknown"
)

func PossibleValuesForStatusType() []string {
	return []string{
		string(StatusTypeFalse),
		string(StatusTypeTrue),
		string(StatusTypeUnknown),
	}
}

func (s *StatusType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStatusType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStatusType(input string) (*StatusType, error) {
	vals := map[string]StatusType{
		"false":   StatusTypeFalse,
		"true":    StatusTypeTrue,
		"unknown": StatusTypeUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StatusType(input)
	return &out, nil
}

type TokenValidationRuleType string

const (
	TokenValidationRuleTypeRequiredClaim TokenValidationRuleType = "RequiredClaim"
)

func PossibleValuesForTokenValidationRuleType() []string {
	return []string{
		string(TokenValidationRuleTypeRequiredClaim),
	}
}

func (s *TokenValidationRuleType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTokenValidationRuleType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTokenValidationRuleType(input string) (*TokenValidationRuleType, error) {
	vals := map[string]TokenValidationRuleType{
		"requiredclaim": TokenValidationRuleTypeRequiredClaim,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TokenValidationRuleType(input)
	return &out, nil
}

type UsernameClaimPrefixPolicy string

const (
	UsernameClaimPrefixPolicyNoPrefix UsernameClaimPrefixPolicy = "NoPrefix"
	UsernameClaimPrefixPolicyNone     UsernameClaimPrefixPolicy = "None"
	UsernameClaimPrefixPolicyPrefix   UsernameClaimPrefixPolicy = "Prefix"
)

func PossibleValuesForUsernameClaimPrefixPolicy() []string {
	return []string{
		string(UsernameClaimPrefixPolicyNoPrefix),
		string(UsernameClaimPrefixPolicyNone),
		string(UsernameClaimPrefixPolicyPrefix),
	}
}

func (s *UsernameClaimPrefixPolicy) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseUsernameClaimPrefixPolicy(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseUsernameClaimPrefixPolicy(input string) (*UsernameClaimPrefixPolicy, error) {
	vals := map[string]UsernameClaimPrefixPolicy{
		"noprefix": UsernameClaimPrefixPolicyNoPrefix,
		"none":     UsernameClaimPrefixPolicyNone,
		"prefix":   UsernameClaimPrefixPolicyPrefix,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UsernameClaimPrefixPolicy(input)
	return &out, nil
}
