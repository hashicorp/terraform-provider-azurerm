package schemaregistry

import "strings"

type CreatedByType string

const (
	CreatedByTypeApplication     CreatedByType = "Application"
	CreatedByTypeKey             CreatedByType = "Key"
	CreatedByTypeManagedIdentity CreatedByType = "ManagedIdentity"
	CreatedByTypeUser            CreatedByType = "User"
)

func PossibleValuesForCreatedByType() []string {
	return []string{
		string(CreatedByTypeApplication),
		string(CreatedByTypeKey),
		string(CreatedByTypeManagedIdentity),
		string(CreatedByTypeUser),
	}
}

func parseCreatedByType(input string) (*CreatedByType, error) {
	vals := map[string]CreatedByType{
		"application":     CreatedByTypeApplication,
		"key":             CreatedByTypeKey,
		"managedidentity": CreatedByTypeManagedIdentity,
		"user":            CreatedByTypeUser,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CreatedByType(input)
	return &out, nil
}

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
