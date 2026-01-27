package secrets

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &SecretId{}

// SecretId is a struct representing the Resource ID for a Secret
type SecretId struct {
	BaseURI    string
	SecretName string
}

// NewSecretID returns a new SecretId struct
func NewSecretID(baseURI string, secretName string) SecretId {
	return SecretId{
		BaseURI:    strings.TrimSuffix(baseURI, "/"),
		SecretName: secretName,
	}
}

// ParseSecretID parses 'input' into a SecretId
func ParseSecretID(input string) (*SecretId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SecretId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SecretId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSecretIDInsensitively parses 'input' case-insensitively into a SecretId
// note: this method should only be used for API response data and not user input
func ParseSecretIDInsensitively(input string) (*SecretId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SecretId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SecretId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SecretId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.BaseURI, ok = input.Parsed["baseURI"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "baseURI", input)
	}

	if id.SecretName, ok = input.Parsed["secretName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "secretName", input)
	}

	return nil
}

// ValidateSecretID checks that 'input' can be parsed as a Secret ID
func ValidateSecretID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSecretID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Secret ID
func (id SecretId) ID() string {
	fmtString := "%s/secrets/%s"
	return fmt.Sprintf(fmtString, id.BaseURI, id.SecretName)
}

// Path returns the formatted Secret ID without the BaseURI
func (id SecretId) Path() string {
	fmtString := "/secrets/%s"
	return fmt.Sprintf(fmtString, id.SecretName)
}

// PathElements returns the values of Secret ID Segments without the BaseURI
func (id SecretId) PathElements() []any {
	return []any{id.SecretName}
}

// Segments returns a slice of Resource ID Segments which comprise this Secret ID
func (id SecretId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.DataPlaneBaseURISegment("baseURI", "https://endpoint_url"),
		resourceids.StaticSegment("staticSecrets", "secrets", "secrets"),
		resourceids.UserSpecifiedSegment("secretName", "secretName"),
	}
}

// String returns a human-readable description of this Secret ID
func (id SecretId) String() string {
	components := []string{
		fmt.Sprintf("Base URI: %q", id.BaseURI),
		fmt.Sprintf("Secret Name: %q", id.SecretName),
	}
	return fmt.Sprintf("Secret (%s)", strings.Join(components, "\n"))
}
