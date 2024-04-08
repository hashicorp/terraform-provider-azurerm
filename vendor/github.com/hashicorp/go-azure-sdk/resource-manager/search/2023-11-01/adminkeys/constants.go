package adminkeys

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AdminKeyKind string

const (
	AdminKeyKindPrimary   AdminKeyKind = "primary"
	AdminKeyKindSecondary AdminKeyKind = "secondary"
)

func PossibleValuesForAdminKeyKind() []string {
	return []string{
		string(AdminKeyKindPrimary),
		string(AdminKeyKindSecondary),
	}
}

func (s *AdminKeyKind) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAdminKeyKind(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAdminKeyKind(input string) (*AdminKeyKind, error) {
	vals := map[string]AdminKeyKind{
		"primary":   AdminKeyKindPrimary,
		"secondary": AdminKeyKindSecondary,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AdminKeyKind(input)
	return &out, nil
}
