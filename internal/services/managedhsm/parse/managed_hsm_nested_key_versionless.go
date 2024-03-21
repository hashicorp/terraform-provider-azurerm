// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// example: https://hsm-name.managedhsm.azure.net/keys/key-name
var _ resourceids.ResourceId = &ManagedHSMNestedkeyVersionless{}

type ManagedHSMNestedkeyVersionless struct {
	BaseURI string
	KeyName string
}

func NewManagedHSMNestedkeyVersionless(baseURI, keyName string) ManagedHSMNestedkeyVersionless {
	return ManagedHSMNestedkeyVersionless{
		BaseURI: baseURI,
		KeyName: keyName,
	}
}

// FromParseResult implements resourceids.ResourceId.
func (id *ManagedHSMNestedkeyVersionless) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.BaseURI, ok = input.Parsed["baseURI"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "baseURI", input)
	}

	if id.KeyName, ok = input.Parsed["keyName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "keyName", input)
	}

	return nil
}

// Segments implements resourceids.ResourceId.
func (id *ManagedHSMNestedkeyVersionless) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.DataPlaneHostSegment("baseURI", "https://hsm-name.managedhsm.azure.net"),
		resourceids.StaticSegment("staticKeys", "keys", "keys"),
		resourceids.UserSpecifiedSegment("keyName", "keyValue"),
	}
}

func (id ManagedHSMNestedkeyVersionless) String() string {
	return fmt.Sprintf("%s: (%s)", "Managed H S M Nested Key Versionless", id.ID())
}

func (id ManagedHSMNestedkeyVersionless) ID() string {
	return fmt.Sprintf("%s/keys/%s", id.BaseURI, id.KeyName)
}

// ManagedHSMNestedkeyVersionless parses a ManagedHSMNestedKeyWithVersion ID into an ManagedHSMNestedkeyVersionless struct
func ParseManagedHSMNestedkeyVersionless(input string) (*ManagedHSMNestedkeyVersionless, error) {
	parser := resourceids.NewParserFromResourceIdType(&ManagedHSMNestedkeyVersionless{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ManagedHSMNestedkeyVersionless{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func ParseManagedHSMNestedkeyVersionlessInsensitively(input string) (*ManagedHSMNestedkeyVersionless, error) {
	parser := resourceids.NewParserFromResourceIdType(&ManagedHSMNestedkeyVersionless{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ManagedHSMNestedkeyVersionless{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func ValidateManagedHSMNestedKeyVersionlessID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseManagedHSMNestedkeyVersionless(v); err != nil {
		errors = append(errors, err)
	}

	return
}
