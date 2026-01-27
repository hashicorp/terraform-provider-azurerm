package settings

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &SettingId{}

// SettingId is a struct representing the Resource ID for a Setting
type SettingId struct {
	BaseURI     string
	SettingName string
}

// NewSettingID returns a new SettingId struct
func NewSettingID(baseURI string, settingName string) SettingId {
	return SettingId{
		BaseURI:     strings.TrimSuffix(baseURI, "/"),
		SettingName: settingName,
	}
}

// ParseSettingID parses 'input' into a SettingId
func ParseSettingID(input string) (*SettingId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SettingId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SettingId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSettingIDInsensitively parses 'input' case-insensitively into a SettingId
// note: this method should only be used for API response data and not user input
func ParseSettingIDInsensitively(input string) (*SettingId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SettingId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SettingId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SettingId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.BaseURI, ok = input.Parsed["baseURI"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "baseURI", input)
	}

	if id.SettingName, ok = input.Parsed["settingName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "settingName", input)
	}

	return nil
}

// ValidateSettingID checks that 'input' can be parsed as a Setting ID
func ValidateSettingID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSettingID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Setting ID
func (id SettingId) ID() string {
	fmtString := "%s/settings/%s"
	return fmt.Sprintf(fmtString, id.BaseURI, id.SettingName)
}

// Path returns the formatted Setting ID without the BaseURI
func (id SettingId) Path() string {
	fmtString := "/settings/%s"
	return fmt.Sprintf(fmtString, id.SettingName)
}

// PathElements returns the values of Setting ID Segments without the BaseURI
func (id SettingId) PathElements() []any {
	return []any{id.SettingName}
}

// Segments returns a slice of Resource ID Segments which comprise this Setting ID
func (id SettingId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.DataPlaneBaseURISegment("baseURI", "https://endpoint_url"),
		resourceids.StaticSegment("staticSettings", "settings", "settings"),
		resourceids.UserSpecifiedSegment("settingName", "settingName"),
	}
}

// String returns a human-readable description of this Setting ID
func (id SettingId) String() string {
	components := []string{
		fmt.Sprintf("Base URI: %q", id.BaseURI),
		fmt.Sprintf("Setting Name: %q", id.SettingName),
	}
	return fmt.Sprintf("Setting (%s)", strings.Join(components, "\n"))
}
