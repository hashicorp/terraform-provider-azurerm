package zones

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ZoneType string

const (
	ZoneTypePrivate ZoneType = "Private"
	ZoneTypePublic  ZoneType = "Public"
)

func PossibleValuesForZoneType() []string {
	return []string{
		string(ZoneTypePrivate),
		string(ZoneTypePublic),
	}
}

func (s *ZoneType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseZoneType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseZoneType(input string) (*ZoneType, error) {
	vals := map[string]ZoneType{
		"private": ZoneTypePrivate,
		"public":  ZoneTypePublic,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ZoneType(input)
	return &out, nil
}
