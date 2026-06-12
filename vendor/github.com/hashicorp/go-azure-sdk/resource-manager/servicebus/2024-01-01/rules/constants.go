package rules

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FilterType string

const (
	FilterTypeCorrelationFilter FilterType = "CorrelationFilter"
	FilterTypeSqlFilter         FilterType = "SqlFilter"
)

func PossibleValuesForFilterType() []string {
	return []string{
		string(FilterTypeCorrelationFilter),
		string(FilterTypeSqlFilter),
	}
}

func (s *FilterType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFilterType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFilterType(input string) (*FilterType, error) {
	vals := map[string]FilterType{
		"correlationfilter": FilterTypeCorrelationFilter,
		"sqlfilter":         FilterTypeSqlFilter,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FilterType(input)
	return &out, nil
}
