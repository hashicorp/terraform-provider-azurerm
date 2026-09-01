package backend

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackendProtocol string

const (
	BackendProtocolHTTP BackendProtocol = "http"
	BackendProtocolSoap BackendProtocol = "soap"
)

func PossibleValuesForBackendProtocol() []string {
	return []string{
		string(BackendProtocolHTTP),
		string(BackendProtocolSoap),
	}
}

func (s *BackendProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBackendProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBackendProtocol(input string) (*BackendProtocol, error) {
	vals := map[string]BackendProtocol{
		"http": BackendProtocolHTTP,
		"soap": BackendProtocolSoap,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BackendProtocol(input)
	return &out, nil
}

type BackendType string

const (
	BackendTypePool   BackendType = "Pool"
	BackendTypeSingle BackendType = "Single"
)

func PossibleValuesForBackendType() []string {
	return []string{
		string(BackendTypePool),
		string(BackendTypeSingle),
	}
}

func (s *BackendType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBackendType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBackendType(input string) (*BackendType, error) {
	vals := map[string]BackendType{
		"pool":   BackendTypePool,
		"single": BackendTypeSingle,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BackendType(input)
	return &out, nil
}
