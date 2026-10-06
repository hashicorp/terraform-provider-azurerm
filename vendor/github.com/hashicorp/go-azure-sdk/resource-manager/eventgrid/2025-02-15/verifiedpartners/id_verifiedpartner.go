package verifiedpartners

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&VerifiedPartnerId{})
}

var _ resourceids.ResourceId = &VerifiedPartnerId{}

// VerifiedPartnerId is a struct representing the Resource ID for a Verified Partner
type VerifiedPartnerId struct {
	VerifiedPartnerName string
}

// NewVerifiedPartnerID returns a new VerifiedPartnerId struct
func NewVerifiedPartnerID(verifiedPartnerName string) VerifiedPartnerId {
	return VerifiedPartnerId{
		VerifiedPartnerName: verifiedPartnerName,
	}
}

// ParseVerifiedPartnerID parses 'input' into a VerifiedPartnerId
func ParseVerifiedPartnerID(input string) (*VerifiedPartnerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VerifiedPartnerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VerifiedPartnerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseVerifiedPartnerIDInsensitively parses 'input' case-insensitively into a VerifiedPartnerId
// note: this method should only be used for API response data and not user input
func ParseVerifiedPartnerIDInsensitively(input string) (*VerifiedPartnerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VerifiedPartnerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VerifiedPartnerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *VerifiedPartnerId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.VerifiedPartnerName, ok = input.Parsed["verifiedPartnerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "verifiedPartnerName", input)
	}

	return nil
}

// ValidateVerifiedPartnerID checks that 'input' can be parsed as a Verified Partner ID
func ValidateVerifiedPartnerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVerifiedPartnerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Verified Partner ID
func (id VerifiedPartnerId) ID() string {
	fmtString := "/providers/Microsoft.EventGrid/verifiedPartners/%s"
	return fmt.Sprintf(fmtString, id.VerifiedPartnerName)
}

// Segments returns a slice of Resource ID Segments which comprise this Verified Partner ID
func (id VerifiedPartnerId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftEventGrid", "Microsoft.EventGrid", "Microsoft.EventGrid"),
		resourceids.StaticSegment("staticVerifiedPartners", "verifiedPartners", "verifiedPartners"),
		resourceids.UserSpecifiedSegment("verifiedPartnerName", "verifiedPartnerName"),
	}
}

// String returns a human-readable description of this Verified Partner ID
func (id VerifiedPartnerId) String() string {
	components := []string{
		fmt.Sprintf("Verified Partner Name: %q", id.VerifiedPartnerName),
	}
	return fmt.Sprintf("Verified Partner (%s)", strings.Join(components, "\n"))
}
