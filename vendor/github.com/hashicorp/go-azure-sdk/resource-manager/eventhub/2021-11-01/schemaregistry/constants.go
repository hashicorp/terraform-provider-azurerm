package schemaregistry

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SchemaCompatibility string

const (
	SchemaCompatibilityBackward SchemaCompatibility = "Backward"
	SchemaCompatibilityForward  SchemaCompatibility = "Forward"
	SchemaCompatibilityNone     SchemaCompatibility = "None"
)

func PossibleValuesForSchemaCompatibility() []string {
	return []string{
		string(SchemaCompatibilityBackward),
		string(SchemaCompatibilityForward),
		string(SchemaCompatibilityNone),
	}
}

func parseSchemaCompatibility(input string) (*SchemaCompatibility, error) {
	vals := map[string]SchemaCompatibility{
		"backward": SchemaCompatibilityBackward,
		"forward":  SchemaCompatibilityForward,
		"none":     SchemaCompatibilityNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SchemaCompatibility(input)
	return &out, nil
}

type SchemaType string

const (
	SchemaTypeAvro    SchemaType = "Avro"
	SchemaTypeUnknown SchemaType = "Unknown"
)

func PossibleValuesForSchemaType() []string {
	return []string{
		string(SchemaTypeAvro),
		string(SchemaTypeUnknown),
	}
}

func parseSchemaType(input string) (*SchemaType, error) {
	vals := map[string]SchemaType{
		"avro":    SchemaTypeAvro,
		"unknown": SchemaTypeUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SchemaType(input)
	return &out, nil
}
