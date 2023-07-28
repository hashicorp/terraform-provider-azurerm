package v2workspaceconnectionresource

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectionAuthType string

const (
	ConnectionAuthTypeManagedIdentity  ConnectionAuthType = "ManagedIdentity"
	ConnectionAuthTypeNone             ConnectionAuthType = "None"
	ConnectionAuthTypePAT              ConnectionAuthType = "PAT"
	ConnectionAuthTypeSAS              ConnectionAuthType = "SAS"
	ConnectionAuthTypeUsernamePassword ConnectionAuthType = "UsernamePassword"
)

func PossibleValuesForConnectionAuthType() []string {
	return []string{
		string(ConnectionAuthTypeManagedIdentity),
		string(ConnectionAuthTypeNone),
		string(ConnectionAuthTypePAT),
		string(ConnectionAuthTypeSAS),
		string(ConnectionAuthTypeUsernamePassword),
	}
}

func (s *ConnectionAuthType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConnectionAuthType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConnectionAuthType(input string) (*ConnectionAuthType, error) {
	vals := map[string]ConnectionAuthType{
		"managedidentity":  ConnectionAuthTypeManagedIdentity,
		"none":             ConnectionAuthTypeNone,
		"pat":              ConnectionAuthTypePAT,
		"sas":              ConnectionAuthTypeSAS,
		"usernamepassword": ConnectionAuthTypeUsernamePassword,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConnectionAuthType(input)
	return &out, nil
}

type ConnectionCategory string

const (
	ConnectionCategoryContainerRegistry ConnectionCategory = "ContainerRegistry"
	ConnectionCategoryGit               ConnectionCategory = "Git"
	ConnectionCategoryPythonFeed        ConnectionCategory = "PythonFeed"
)

func PossibleValuesForConnectionCategory() []string {
	return []string{
		string(ConnectionCategoryContainerRegistry),
		string(ConnectionCategoryGit),
		string(ConnectionCategoryPythonFeed),
	}
}

func (s *ConnectionCategory) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConnectionCategory(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConnectionCategory(input string) (*ConnectionCategory, error) {
	vals := map[string]ConnectionCategory{
		"containerregistry": ConnectionCategoryContainerRegistry,
		"git":               ConnectionCategoryGit,
		"pythonfeed":        ConnectionCategoryPythonFeed,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConnectionCategory(input)
	return &out, nil
}

type ValueFormat string

const (
	ValueFormatJSON ValueFormat = "JSON"
)

func PossibleValuesForValueFormat() []string {
	return []string{
		string(ValueFormatJSON),
	}
}

func (s *ValueFormat) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseValueFormat(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseValueFormat(input string) (*ValueFormat, error) {
	vals := map[string]ValueFormat{
		"json": ValueFormatJSON,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ValueFormat(input)
	return &out, nil
}
