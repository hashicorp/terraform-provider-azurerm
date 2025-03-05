package issue

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type State string

const (
	StateClosed   State = "closed"
	StateOpen     State = "open"
	StateProposed State = "proposed"
	StateRemoved  State = "removed"
	StateResolved State = "resolved"
)

func PossibleValuesForState() []string {
	return []string{
		string(StateClosed),
		string(StateOpen),
		string(StateProposed),
		string(StateRemoved),
		string(StateResolved),
	}
}

func (s *State) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseState(input string) (*State, error) {
	vals := map[string]State{
		"closed":   StateClosed,
		"open":     StateOpen,
		"proposed": StateProposed,
		"removed":  StateRemoved,
		"resolved": StateResolved,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := State(input)
	return &out, nil
}
