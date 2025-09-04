package backups

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Origin string

const (
	OriginCustomerOnNegativeDemand Origin = "Customer On-Demand"
	OriginFull                     Origin = "Full"
)

func PossibleValuesForOrigin() []string {
	return []string{
		string(OriginCustomerOnNegativeDemand),
		string(OriginFull),
	}
}

func (s *Origin) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOrigin(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOrigin(input string) (*Origin, error) {
	vals := map[string]Origin{
		"customer on-demand": OriginCustomerOnNegativeDemand,
		"full":               OriginFull,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Origin(input)
	return &out, nil
}
