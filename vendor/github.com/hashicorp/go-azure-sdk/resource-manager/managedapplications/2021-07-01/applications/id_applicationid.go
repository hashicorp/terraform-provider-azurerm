package applications

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ApplicationIdId{})
}

var _ resourceids.ResourceId = &ApplicationIdId{}

// ApplicationIdId is a struct representing the Resource ID for a Application Id
type ApplicationIdId struct {
	ApplicationId string
}

// NewApplicationIdID returns a new ApplicationIdId struct
func NewApplicationIdID(applicationId string) ApplicationIdId {
	return ApplicationIdId{
		ApplicationId: applicationId,
	}
}

// ParseApplicationIdID parses 'input' into a ApplicationIdId
func ParseApplicationIdID(input string) (*ApplicationIdId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ApplicationIdId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ApplicationIdId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseApplicationIdIDInsensitively parses 'input' case-insensitively into a ApplicationIdId
// note: this method should only be used for API response data and not user input
func ParseApplicationIdIDInsensitively(input string) (*ApplicationIdId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ApplicationIdId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ApplicationIdId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ApplicationIdId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.ApplicationId, ok = input.Parsed["applicationId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "applicationId", input)
	}

	return nil
}

// ValidateApplicationIdID checks that 'input' can be parsed as a Application Id ID
func ValidateApplicationIdID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseApplicationIdID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Application Id ID
func (id ApplicationIdId) ID() string {
	fmtString := "/%s"
	return fmt.Sprintf(fmtString, id.ApplicationId)
}

// Segments returns a slice of Resource ID Segments which comprise this Application Id ID
func (id ApplicationIdId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.UserSpecifiedSegment("applicationId", "applicationId"),
	}
}

// String returns a human-readable description of this Application Id ID
func (id ApplicationIdId) String() string {
	components := []string{
		fmt.Sprintf("Application: %q", id.ApplicationId),
	}
	return fmt.Sprintf("Application Id (%s)", strings.Join(components, "\n"))
}
