package encryptionscopes

import "strings"

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
