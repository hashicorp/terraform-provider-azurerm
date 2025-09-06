package managedidentities

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IsolationScope string

const (
	IsolationScopeNone     IsolationScope = "None"
	IsolationScopeRegional IsolationScope = "Regional"
)

func PossibleValuesForIsolationScope() []string {
	return []string{
		string(IsolationScopeNone),
		string(IsolationScopeRegional),
	}
}

func (s *IsolationScope) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIsolationScope(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIsolationScope(input string) (*IsolationScope, error) {
	vals := map[string]IsolationScope{
		"none":     IsolationScopeNone,
		"regional": IsolationScopeRegional,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IsolationScope(input)
	return &out, nil
}
