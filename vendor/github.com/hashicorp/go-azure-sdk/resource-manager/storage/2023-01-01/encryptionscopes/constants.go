package encryptionscopes

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EncryptionScopeSource string

const (
	EncryptionScopeSourceMicrosoftPointKeyVault EncryptionScopeSource = "Microsoft.KeyVault"
	EncryptionScopeSourceMicrosoftPointStorage  EncryptionScopeSource = "Microsoft.Storage"
)

func PossibleValuesForEncryptionScopeSource() []string {
	return []string{
		string(EncryptionScopeSourceMicrosoftPointKeyVault),
		string(EncryptionScopeSourceMicrosoftPointStorage),
	}
}

func (s *EncryptionScopeSource) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEncryptionScopeSource(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEncryptionScopeSource(input string) (*EncryptionScopeSource, error) {
	vals := map[string]EncryptionScopeSource{
		"microsoft.keyvault": EncryptionScopeSourceMicrosoftPointKeyVault,
		"microsoft.storage":  EncryptionScopeSourceMicrosoftPointStorage,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EncryptionScopeSource(input)
	return &out, nil
}

type EncryptionScopeState string

const (
	EncryptionScopeStateDisabled EncryptionScopeState = "Disabled"
	EncryptionScopeStateEnabled  EncryptionScopeState = "Enabled"
)

func PossibleValuesForEncryptionScopeState() []string {
	return []string{
		string(EncryptionScopeStateDisabled),
		string(EncryptionScopeStateEnabled),
	}
}

func (s *EncryptionScopeState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEncryptionScopeState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEncryptionScopeState(input string) (*EncryptionScopeState, error) {
	vals := map[string]EncryptionScopeState{
		"disabled": EncryptionScopeStateDisabled,
		"enabled":  EncryptionScopeStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EncryptionScopeState(input)
	return &out, nil
}

type ListEncryptionScopesInclude string

const (
	ListEncryptionScopesIncludeAll      ListEncryptionScopesInclude = "All"
	ListEncryptionScopesIncludeDisabled ListEncryptionScopesInclude = "Disabled"
	ListEncryptionScopesIncludeEnabled  ListEncryptionScopesInclude = "Enabled"
)

func PossibleValuesForListEncryptionScopesInclude() []string {
	return []string{
		string(ListEncryptionScopesIncludeAll),
		string(ListEncryptionScopesIncludeDisabled),
		string(ListEncryptionScopesIncludeEnabled),
	}
}

func (s *ListEncryptionScopesInclude) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseListEncryptionScopesInclude(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseListEncryptionScopesInclude(input string) (*ListEncryptionScopesInclude, error) {
	vals := map[string]ListEncryptionScopesInclude{
		"all":      ListEncryptionScopesIncludeAll,
		"disabled": ListEncryptionScopesIncludeDisabled,
		"enabled":  ListEncryptionScopesIncludeEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ListEncryptionScopesInclude(input)
	return &out, nil
}
