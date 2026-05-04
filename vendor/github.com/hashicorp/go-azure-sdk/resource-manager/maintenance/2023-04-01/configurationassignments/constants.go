package configurationassignments

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TagOperators string

const (
	TagOperatorsAll TagOperators = "All"
	TagOperatorsAny TagOperators = "Any"
)

func PossibleValuesForTagOperators() []string {
	return []string{
		string(TagOperatorsAll),
		string(TagOperatorsAny),
	}
}

func (s *TagOperators) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTagOperators(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTagOperators(input string) (*TagOperators, error) {
	vals := map[string]TagOperators{
		"all": TagOperatorsAll,
		"any": TagOperatorsAny,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TagOperators(input)
	return &out, nil
}
