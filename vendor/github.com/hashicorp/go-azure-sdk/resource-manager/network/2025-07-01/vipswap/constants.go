package vipswap

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SlotType string

const (
	SlotTypeProduction SlotType = "Production"
	SlotTypeStaging    SlotType = "Staging"
)

func PossibleValuesForSlotType() []string {
	return []string{
		string(SlotTypeProduction),
		string(SlotTypeStaging),
	}
}

func (s *SlotType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSlotType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSlotType(input string) (*SlotType, error) {
	vals := map[string]SlotType{
		"production": SlotTypeProduction,
		"staging":    SlotTypeStaging,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SlotType(input)
	return &out, nil
}
