package schema

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SchemaType string

const (
	SchemaTypeJson SchemaType = "json"
	SchemaTypeXml  SchemaType = "xml"
)

func PossibleValuesForSchemaType() []string {
	return []string{
		string(SchemaTypeJson),
		string(SchemaTypeXml),
	}
}

func (s *SchemaType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSchemaType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSchemaType(input string) (*SchemaType, error) {
	vals := map[string]SchemaType{
		"json": SchemaTypeJson,
		"xml":  SchemaTypeXml,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SchemaType(input)
	return &out, nil
}
