package usages

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UsageUnit string

const (
	UsageUnitCount UsageUnit = "Count"
)

func PossibleValuesForUsageUnit() []string {
	return []string{
		string(UsageUnitCount),
	}
}

func (s *UsageUnit) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseUsageUnit(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseUsageUnit(input string) (*UsageUnit, error) {
	vals := map[string]UsageUnit{
		"count": UsageUnitCount,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UsageUnit(input)
	return &out, nil
}
