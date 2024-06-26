package accounts

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IdentityType string

const (
	IdentityTypeDelegatedResourceIdentity IdentityType = "delegatedResourceIdentity"
	IdentityTypeSystemAssignedIdentity    IdentityType = "systemAssignedIdentity"
	IdentityTypeUserAssignedIdentity      IdentityType = "userAssignedIdentity"
)

func PossibleValuesForIdentityType() []string {
	return []string{
		string(IdentityTypeDelegatedResourceIdentity),
		string(IdentityTypeSystemAssignedIdentity),
		string(IdentityTypeUserAssignedIdentity),
	}
}

func (s *IdentityType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIdentityType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIdentityType(input string) (*IdentityType, error) {
	vals := map[string]IdentityType{
		"delegatedresourceidentity": IdentityTypeDelegatedResourceIdentity,
		"systemassignedidentity":    IdentityTypeSystemAssignedIdentity,
		"userassignedidentity":      IdentityTypeUserAssignedIdentity,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IdentityType(input)
	return &out, nil
}

type InfrastructureEncryption string

const (
	InfrastructureEncryptionDisabled InfrastructureEncryption = "disabled"
	InfrastructureEncryptionEnabled  InfrastructureEncryption = "enabled"
)

func PossibleValuesForInfrastructureEncryption() []string {
	return []string{
		string(InfrastructureEncryptionDisabled),
		string(InfrastructureEncryptionEnabled),
	}
}

func (s *InfrastructureEncryption) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseInfrastructureEncryption(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseInfrastructureEncryption(input string) (*InfrastructureEncryption, error) {
	vals := map[string]InfrastructureEncryption{
		"disabled": InfrastructureEncryptionDisabled,
		"enabled":  InfrastructureEncryptionEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := InfrastructureEncryption(input)
	return &out, nil
}

type KeyType string

const (
	KeyTypePrimary   KeyType = "primary"
	KeyTypeSecondary KeyType = "secondary"
)

func PossibleValuesForKeyType() []string {
	return []string{
		string(KeyTypePrimary),
		string(KeyTypeSecondary),
	}
}

func (s *KeyType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKeyType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKeyType(input string) (*KeyType, error) {
	vals := map[string]KeyType{
		"primary":   KeyTypePrimary,
		"secondary": KeyTypeSecondary,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KeyType(input)
	return &out, nil
}

type Kind string

const (
	KindGenOne Kind = "Gen1"
	KindGenTwo Kind = "Gen2"
)

func PossibleValuesForKind() []string {
	return []string{
		string(KindGenOne),
		string(KindGenTwo),
	}
}

func (s *Kind) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKind(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKind(input string) (*Kind, error) {
	vals := map[string]Kind{
		"gen1": KindGenOne,
		"gen2": KindGenTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Kind(input)
	return &out, nil
}

type Name string

const (
	NameGTwo  Name = "G2"
	NameSOne  Name = "S1"
	NameSZero Name = "S0"
)

func PossibleValuesForName() []string {
	return []string{
		string(NameGTwo),
		string(NameSOne),
		string(NameSZero),
	}
}

func (s *Name) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseName(input string) (*Name, error) {
	vals := map[string]Name{
		"g2": NameGTwo,
		"s1": NameSOne,
		"s0": NameSZero,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Name(input)
	return &out, nil
}

type SigningKey string

const (
	SigningKeyManagedIdentity SigningKey = "managedIdentity"
	SigningKeyPrimaryKey      SigningKey = "primaryKey"
	SigningKeySecondaryKey    SigningKey = "secondaryKey"
)

func PossibleValuesForSigningKey() []string {
	return []string{
		string(SigningKeyManagedIdentity),
		string(SigningKeyPrimaryKey),
		string(SigningKeySecondaryKey),
	}
}

func (s *SigningKey) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSigningKey(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSigningKey(input string) (*SigningKey, error) {
	vals := map[string]SigningKey{
		"managedidentity": SigningKeyManagedIdentity,
		"primarykey":      SigningKeyPrimaryKey,
		"secondarykey":    SigningKeySecondaryKey,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SigningKey(input)
	return &out, nil
}
