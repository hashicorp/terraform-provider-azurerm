package secrets

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &SecretversionId{}

// SecretversionId is a struct representing the Resource ID for a Secretversion
type SecretversionId struct {
	BaseURI       string
	SecretName    string
	Secretversion string
}

// NewSecretversionID returns a new SecretversionId struct
func NewSecretversionID(baseURI string, secretName string, secretversion string) SecretversionId {
	return SecretversionId{
		BaseURI:       strings.TrimSuffix(baseURI, "/"),
		SecretName:    secretName,
		Secretversion: secretversion,
	}
}

// ParseSecretversionID parses 'input' into a SecretversionId
func ParseSecretversionID(input string) (*SecretversionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SecretversionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SecretversionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSecretversionIDInsensitively parses 'input' case-insensitively into a SecretversionId
// note: this method should only be used for API response data and not user input
func ParseSecretversionIDInsensitively(input string) (*SecretversionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SecretversionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SecretversionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SecretversionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.BaseURI, ok = input.Parsed["baseURI"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "baseURI", input)
	}

	if id.SecretName, ok = input.Parsed["secretName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "secretName", input)
	}

	if id.Secretversion, ok = input.Parsed["secretversion"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "secretversion", input)
	}

	return nil
}

// ValidateSecretversionID checks that 'input' can be parsed as a Secretversion ID
func ValidateSecretversionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSecretversionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Secretversion ID
func (id SecretversionId) ID() string {
	fmtString := "%s/secrets/%s/%s"
	return fmt.Sprintf(fmtString, strings.TrimSuffix(id.BaseURI, "/"), id.SecretName, id.Secretversion)
}

// Path returns the formatted Secretversion ID without the BaseURI
func (id SecretversionId) Path() string {
	fmtString := "/secrets/%s/%s"
	return fmt.Sprintf(fmtString, id.SecretName, id.Secretversion)
}

// PathElements returns the values of Secretversion ID Segments without the BaseURI
func (id SecretversionId) PathElements() []any {
	return []any{id.SecretName, id.Secretversion}
}

// Segments returns a slice of Resource ID Segments which comprise this Secretversion ID
func (id SecretversionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.DataPlaneBaseURISegment("baseURI", "https://endpoint-url.example.com"),
		resourceids.StaticSegment("staticSecrets", "secrets", "secrets"),
		resourceids.UserSpecifiedSegment("secretName", "secretName"),
		resourceids.UserSpecifiedSegment("secretversion", "secretversion"),
	}
}

// String returns a human-readable description of this Secretversion ID
func (id SecretversionId) String() string {
	components := []string{
		fmt.Sprintf("Base URI: %q", id.BaseURI),
		fmt.Sprintf("Secret Name: %q", id.SecretName),
		fmt.Sprintf("Secretversion: %q", id.Secretversion),
	}
	return fmt.Sprintf("Secretversion (%s)", strings.Join(components, "\n"))
}
