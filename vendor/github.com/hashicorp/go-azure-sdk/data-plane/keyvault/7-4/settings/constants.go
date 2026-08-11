package settings

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SettingTypeEnum string

const (
	SettingTypeEnumBoolean SettingTypeEnum = "boolean"
)

func PossibleValuesForSettingTypeEnum() []string {
	return []string{
		string(SettingTypeEnumBoolean),
	}
}

func (s *SettingTypeEnum) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSettingTypeEnum(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSettingTypeEnum(input string) (*SettingTypeEnum, error) {
	vals := map[string]SettingTypeEnum{
		"boolean": SettingTypeEnumBoolean,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SettingTypeEnum(input)
	return &out, nil
}
