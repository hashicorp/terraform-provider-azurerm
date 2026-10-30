package graphquery

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResultKind string

const (
	ResultKindBasic ResultKind = "basic"
)

func PossibleValuesForResultKind() []string {
	return []string{
		string(ResultKindBasic),
	}
}

func (s *ResultKind) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseResultKind(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseResultKind(input string) (*ResultKind, error) {
	vals := map[string]ResultKind{
		"basic": ResultKindBasic,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResultKind(input)
	return &out, nil
}
