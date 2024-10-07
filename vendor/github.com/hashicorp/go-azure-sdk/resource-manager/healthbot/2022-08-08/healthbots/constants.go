package healthbots

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SkuName string

const (
	SkuNameCZero SkuName = "C0"
	SkuNameFZero SkuName = "F0"
	SkuNameSOne  SkuName = "S1"
)

func PossibleValuesForSkuName() []string {
	return []string{
		string(SkuNameCZero),
		string(SkuNameFZero),
		string(SkuNameSOne),
	}
}

func (s *SkuName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSkuName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSkuName(input string) (*SkuName, error) {
	vals := map[string]SkuName{
		"c0": SkuNameCZero,
		"f0": SkuNameFZero,
		"s1": SkuNameSOne,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuName(input)
	return &out, nil
}
