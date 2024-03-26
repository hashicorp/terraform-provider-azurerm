// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// example: https://hsm-name.managedhsm.azure.net/keys/key-name/version-name
var _ resourceids.ResourceId = &ManagedHSMNestedKeyWithVersionId{}

type ManagedHSMNestedKeyWithVersionId struct {
	BaseURI string
	KeyName string
	Version string
}

func NewManagedHSMNestedKeyWithVersionID(baseURI, keyName, versionName string) ManagedHSMNestedKeyWithVersionId {
	return ManagedHSMNestedKeyWithVersionId{
		BaseURI: baseURI,
		KeyName: keyName,
		Version: versionName,
	}
}

// FromParseResult implements resourceids.ResourceId.
func (id *ManagedHSMNestedKeyWithVersionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.BaseURI, ok = input.Parsed["baseURI"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "baseURI", input)
	}

	if id.KeyName, ok = input.Parsed["keyName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "keyName", input)
	}

	if id.Version, ok = input.Parsed["versionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "versionName", input)
	}

	return nil
}

// Segments implements resourceids.ResourceId.
func (id *ManagedHSMNestedKeyWithVersionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.DataPlaneHostSegment("baseURI", "https://hsm-name.managedhsm.azure.net"),
		resourceids.StaticSegment("staticKeys", "keys", "keys"),
		resourceids.UserSpecifiedSegment("keyName", "keyValue"),
		resourceids.UserSpecifiedSegment("versionName", "versionValue"),
	}
}

func (id ManagedHSMNestedKeyWithVersionId) String() string {
	return fmt.Sprintf("%s: (%s)", "Managed H S M Nested Key With Version", id.ID())
}

func (id ManagedHSMNestedKeyWithVersionId) ID() string {
	fmtString := "%s/keys/%s/%s"
	return fmt.Sprintf(fmtString, id.BaseURI, id.KeyName, id.Version)
}

// ManagedHSMNestedKeyWithVersionID parses a ManagedHSMNestedKeyWithVersion ID into an ManagedHSMNestedKeyWithVersionId struct
func ParseManagedHSMNestedKeyWithVersionID(input string) (*ManagedHSMNestedKeyWithVersionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ManagedHSMNestedKeyWithVersionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ManagedHSMNestedKeyWithVersionId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func ParseManagedHSMNestedKeyWithVersionIDInsensitively(input string) (*ManagedHSMNestedKeyWithVersionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ManagedHSMNestedKeyWithVersionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ManagedHSMNestedKeyWithVersionId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func ValidateManagedHSMNestedKeyWithVersionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseManagedHSMNestedKeyWithVersionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}
