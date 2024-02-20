package functions

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UdfType string

const (
	UdfTypeScalar UdfType = "Scalar"
)

func PossibleValuesForUdfType() []string {
	return []string{
		string(UdfTypeScalar),
	}
}

func (s *UdfType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseUdfType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseUdfType(input string) (*UdfType, error) {
	vals := map[string]UdfType{
		"scalar": UdfTypeScalar,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UdfType(input)
	return &out, nil
}
