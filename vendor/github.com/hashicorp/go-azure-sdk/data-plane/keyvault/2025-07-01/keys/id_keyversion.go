package keys

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &KeyversionId{}

// KeyversionId is a struct representing the Resource ID for a Keyversion
type KeyversionId struct {
	BaseURI    string
	KeyName    string
	Keyversion string
}

// NewKeyversionID returns a new KeyversionId struct
func NewKeyversionID(baseURI string, keyName string, keyversion string) KeyversionId {
	return KeyversionId{
		BaseURI:    strings.TrimSuffix(baseURI, "/"),
		KeyName:    keyName,
		Keyversion: keyversion,
	}
}

// ParseKeyversionID parses 'input' into a KeyversionId
func ParseKeyversionID(input string) (*KeyversionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&KeyversionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := KeyversionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseKeyversionIDInsensitively parses 'input' case-insensitively into a KeyversionId
// note: this method should only be used for API response data and not user input
func ParseKeyversionIDInsensitively(input string) (*KeyversionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&KeyversionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := KeyversionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *KeyversionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.BaseURI, ok = input.Parsed["baseURI"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "baseURI", input)
	}

	if id.KeyName, ok = input.Parsed["keyName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "keyName", input)
	}

	if id.Keyversion, ok = input.Parsed["keyversion"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "keyversion", input)
	}

	return nil
}

// ValidateKeyversionID checks that 'input' can be parsed as a Keyversion ID
func ValidateKeyversionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseKeyversionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Keyversion ID
func (id KeyversionId) ID() string {
	fmtString := "%s/keys/%s/%s"
	return fmt.Sprintf(fmtString, id.BaseURI, id.KeyName, id.Keyversion)
}

// Path returns the formatted Keyversion ID without the BaseURI
func (id KeyversionId) Path() string {
	fmtString := "/keys/%s/%s"
	return fmt.Sprintf(fmtString, id.KeyName, id.Keyversion)
}

// PathElements returns the values of Keyversion ID Segments without the BaseURI
func (id KeyversionId) PathElements() []any {
	return []any{id.KeyName, id.Keyversion}
}

// Segments returns a slice of Resource ID Segments which comprise this Keyversion ID
func (id KeyversionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.DataPlaneBaseURISegment("baseURI", "https://endpoint_url"),
		resourceids.StaticSegment("staticKeys", "keys", "keys"),
		resourceids.UserSpecifiedSegment("keyName", "keyName"),
		resourceids.UserSpecifiedSegment("keyversion", "keyversion"),
	}
}

// String returns a human-readable description of this Keyversion ID
func (id KeyversionId) String() string {
	components := []string{
		fmt.Sprintf("Base URI: %q", id.BaseURI),
		fmt.Sprintf("Key Name: %q", id.KeyName),
		fmt.Sprintf("Keyversion: %q", id.Keyversion),
	}
	return fmt.Sprintf("Keyversion (%s)", strings.Join(components, "\n"))
}
