package settings

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ScopedSettingId{})
}

var _ resourceids.ResourceId = &ScopedSettingId{}

// ScopedSettingId is a struct representing the Resource ID for a Scoped Setting
type ScopedSettingId struct {
	Scope       string
	SettingName string
}

// NewScopedSettingID returns a new ScopedSettingId struct
func NewScopedSettingID(scope string, settingName string) ScopedSettingId {
	return ScopedSettingId{
		Scope:       scope,
		SettingName: settingName,
	}
}

// ParseScopedSettingID parses 'input' into a ScopedSettingId
func ParseScopedSettingID(input string) (*ScopedSettingId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedSettingId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedSettingId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseScopedSettingIDInsensitively parses 'input' case-insensitively into a ScopedSettingId
// note: this method should only be used for API response data and not user input
func ParseScopedSettingIDInsensitively(input string) (*ScopedSettingId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedSettingId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedSettingId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ScopedSettingId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.Scope, ok = input.Parsed["scope"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "scope", input)
	}

	if id.SettingName, ok = input.Parsed["settingName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "settingName", input)
	}

	return nil
}

// ValidateScopedSettingID checks that 'input' can be parsed as a Scoped Setting ID
func ValidateScopedSettingID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedSettingID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Setting ID
func (id ScopedSettingId) ID() string {
	fmtString := "/%s/providers/Microsoft.HybridCompute/settings/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.Scope, "/"), id.SettingName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Setting ID
func (id ScopedSettingId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("scope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftHybridCompute", "Microsoft.HybridCompute", "Microsoft.HybridCompute"),
		resourceids.StaticSegment("staticSettings", "settings", "settings"),
		resourceids.UserSpecifiedSegment("settingName", "settingValue"),
	}
}

// String returns a human-readable description of this Scoped Setting ID
func (id ScopedSettingId) String() string {
	components := []string{
		fmt.Sprintf("Scope: %q", id.Scope),
		fmt.Sprintf("Setting Name: %q", id.SettingName),
	}
	return fmt.Sprintf("Scoped Setting (%s)", strings.Join(components, "\n"))
}
