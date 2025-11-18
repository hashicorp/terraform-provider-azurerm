package netappaccounts

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActiveDirectoryStatus string

const (
	ActiveDirectoryStatusCreated  ActiveDirectoryStatus = "Created"
	ActiveDirectoryStatusDeleted  ActiveDirectoryStatus = "Deleted"
	ActiveDirectoryStatusError    ActiveDirectoryStatus = "Error"
	ActiveDirectoryStatusInUse    ActiveDirectoryStatus = "InUse"
	ActiveDirectoryStatusUpdating ActiveDirectoryStatus = "Updating"
)

func PossibleValuesForActiveDirectoryStatus() []string {
	return []string{
		string(ActiveDirectoryStatusCreated),
		string(ActiveDirectoryStatusDeleted),
		string(ActiveDirectoryStatusError),
		string(ActiveDirectoryStatusInUse),
		string(ActiveDirectoryStatusUpdating),
	}
}

func (s *ActiveDirectoryStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseActiveDirectoryStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseActiveDirectoryStatus(input string) (*ActiveDirectoryStatus, error) {
	vals := map[string]ActiveDirectoryStatus{
		"created":  ActiveDirectoryStatusCreated,
		"deleted":  ActiveDirectoryStatusDeleted,
		"error":    ActiveDirectoryStatusError,
		"inuse":    ActiveDirectoryStatusInUse,
		"updating": ActiveDirectoryStatusUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ActiveDirectoryStatus(input)
	return &out, nil
}

type KeySource string

const (
	KeySourceMicrosoftPointKeyVault KeySource = "Microsoft.KeyVault"
	KeySourceMicrosoftPointNetApp   KeySource = "Microsoft.NetApp"
)

func PossibleValuesForKeySource() []string {
	return []string{
		string(KeySourceMicrosoftPointKeyVault),
		string(KeySourceMicrosoftPointNetApp),
	}
}

func (s *KeySource) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKeySource(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKeySource(input string) (*KeySource, error) {
	vals := map[string]KeySource{
		"microsoft.keyvault": KeySourceMicrosoftPointKeyVault,
		"microsoft.netapp":   KeySourceMicrosoftPointNetApp,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KeySource(input)
	return &out, nil
}

type KeyVaultStatus string

const (
	KeyVaultStatusCreated  KeyVaultStatus = "Created"
	KeyVaultStatusDeleted  KeyVaultStatus = "Deleted"
	KeyVaultStatusError    KeyVaultStatus = "Error"
	KeyVaultStatusInUse    KeyVaultStatus = "InUse"
	KeyVaultStatusUpdating KeyVaultStatus = "Updating"
)

func PossibleValuesForKeyVaultStatus() []string {
	return []string{
		string(KeyVaultStatusCreated),
		string(KeyVaultStatusDeleted),
		string(KeyVaultStatusError),
		string(KeyVaultStatusInUse),
		string(KeyVaultStatusUpdating),
	}
}

func (s *KeyVaultStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKeyVaultStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKeyVaultStatus(input string) (*KeyVaultStatus, error) {
	vals := map[string]KeyVaultStatus{
		"created":  KeyVaultStatusCreated,
		"deleted":  KeyVaultStatusDeleted,
		"error":    KeyVaultStatusError,
		"inuse":    KeyVaultStatusInUse,
		"updating": KeyVaultStatusUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KeyVaultStatus(input)
	return &out, nil
}

type MultiAdStatus string

const (
	MultiAdStatusDisabled MultiAdStatus = "Disabled"
	MultiAdStatusEnabled  MultiAdStatus = "Enabled"
)

func PossibleValuesForMultiAdStatus() []string {
	return []string{
		string(MultiAdStatusDisabled),
		string(MultiAdStatusEnabled),
	}
}

func (s *MultiAdStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMultiAdStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMultiAdStatus(input string) (*MultiAdStatus, error) {
	vals := map[string]MultiAdStatus{
		"disabled": MultiAdStatusDisabled,
		"enabled":  MultiAdStatusEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MultiAdStatus(input)
	return &out, nil
}
