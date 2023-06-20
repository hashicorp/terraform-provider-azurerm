package integrationaccountschemas

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KeyType string

const (
	KeyTypeNotSpecified KeyType = "NotSpecified"
	KeyTypePrimary      KeyType = "Primary"
	KeyTypeSecondary    KeyType = "Secondary"
)

func PossibleValuesForKeyType() []string {
	return []string{
		string(KeyTypeNotSpecified),
		string(KeyTypePrimary),
		string(KeyTypeSecondary),
	}
}

func parseKeyType(input string) (*KeyType, error) {
	vals := map[string]KeyType{
		"notspecified": KeyTypeNotSpecified,
		"primary":      KeyTypePrimary,
		"secondary":    KeyTypeSecondary,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KeyType(input)
	return &out, nil
}

type SchemaType string

const (
	SchemaTypeNotSpecified SchemaType = "NotSpecified"
	SchemaTypeXml          SchemaType = "Xml"
)

func PossibleValuesForSchemaType() []string {
	return []string{
		string(SchemaTypeNotSpecified),
		string(SchemaTypeXml),
	}
}

func parseSchemaType(input string) (*SchemaType, error) {
	vals := map[string]SchemaType{
		"notspecified": SchemaTypeNotSpecified,
		"xml":          SchemaTypeXml,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SchemaType(input)
	return &out, nil
}
