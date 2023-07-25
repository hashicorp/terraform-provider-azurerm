package cloudendpointresource

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ChangeDetectionMode string

const (
	ChangeDetectionModeDefault   ChangeDetectionMode = "Default"
	ChangeDetectionModeRecursive ChangeDetectionMode = "Recursive"
)

func PossibleValuesForChangeDetectionMode() []string {
	return []string{
		string(ChangeDetectionModeDefault),
		string(ChangeDetectionModeRecursive),
	}
}

func (s *ChangeDetectionMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseChangeDetectionMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseChangeDetectionMode(input string) (*ChangeDetectionMode, error) {
	vals := map[string]ChangeDetectionMode{
		"default":   ChangeDetectionModeDefault,
		"recursive": ChangeDetectionModeRecursive,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ChangeDetectionMode(input)
	return &out, nil
}
