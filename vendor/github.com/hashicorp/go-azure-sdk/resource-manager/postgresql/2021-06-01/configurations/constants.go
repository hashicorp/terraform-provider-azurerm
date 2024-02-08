package configurations

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigurationDataType string

const (
	ConfigurationDataTypeBoolean     ConfigurationDataType = "Boolean"
	ConfigurationDataTypeEnumeration ConfigurationDataType = "Enumeration"
	ConfigurationDataTypeInteger     ConfigurationDataType = "Integer"
	ConfigurationDataTypeNumeric     ConfigurationDataType = "Numeric"
)

func PossibleValuesForConfigurationDataType() []string {
	return []string{
		string(ConfigurationDataTypeBoolean),
		string(ConfigurationDataTypeEnumeration),
		string(ConfigurationDataTypeInteger),
		string(ConfigurationDataTypeNumeric),
	}
}

func (s *ConfigurationDataType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConfigurationDataType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConfigurationDataType(input string) (*ConfigurationDataType, error) {
	vals := map[string]ConfigurationDataType{
		"boolean":     ConfigurationDataTypeBoolean,
		"enumeration": ConfigurationDataTypeEnumeration,
		"integer":     ConfigurationDataTypeInteger,
		"numeric":     ConfigurationDataTypeNumeric,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConfigurationDataType(input)
	return &out, nil
}
