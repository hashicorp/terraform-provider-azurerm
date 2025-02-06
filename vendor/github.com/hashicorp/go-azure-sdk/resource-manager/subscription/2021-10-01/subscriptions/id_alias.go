package subscriptions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&AliasId{})
}

var _ resourceids.ResourceId = &AliasId{}

// AliasId is a struct representing the Resource ID for a Alias
type AliasId struct {
	AliasName string
}

// NewAliasID returns a new AliasId struct
func NewAliasID(aliasName string) AliasId {
	return AliasId{
		AliasName: aliasName,
	}
}

// ParseAliasID parses 'input' into a AliasId
func ParseAliasID(input string) (*AliasId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AliasId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AliasId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAliasIDInsensitively parses 'input' case-insensitively into a AliasId
// note: this method should only be used for API response data and not user input
func ParseAliasIDInsensitively(input string) (*AliasId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AliasId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AliasId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AliasId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.AliasName, ok = input.Parsed["aliasName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "aliasName", input)
	}

	return nil
}

// ValidateAliasID checks that 'input' can be parsed as a Alias ID
func ValidateAliasID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAliasID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Alias ID
func (id AliasId) ID() string {
	fmtString := "/providers/Microsoft.Subscription/aliases/%s"
	return fmt.Sprintf(fmtString, id.AliasName)
}

// Segments returns a slice of Resource ID Segments which comprise this Alias ID
func (id AliasId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSubscription", "Microsoft.Subscription", "Microsoft.Subscription"),
		resourceids.StaticSegment("staticAliases", "aliases", "aliases"),
		resourceids.UserSpecifiedSegment("aliasName", "aliasName"),
	}
}

// String returns a human-readable description of this Alias ID
func (id AliasId) String() string {
	components := []string{
		fmt.Sprintf("Alias Name: %q", id.AliasName),
	}
	return fmt.Sprintf("Alias (%s)", strings.Join(components, "\n"))
}
