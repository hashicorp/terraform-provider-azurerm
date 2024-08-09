package agentversions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&AgentVersionId{})
}

var _ resourceids.ResourceId = &AgentVersionId{}

// AgentVersionId is a struct representing the Resource ID for a Agent Version
type AgentVersionId struct {
	OsTypeName       string
	AgentVersionName string
}

// NewAgentVersionID returns a new AgentVersionId struct
func NewAgentVersionID(osTypeName string, agentVersionName string) AgentVersionId {
	return AgentVersionId{
		OsTypeName:       osTypeName,
		AgentVersionName: agentVersionName,
	}
}

// ParseAgentVersionID parses 'input' into a AgentVersionId
func ParseAgentVersionID(input string) (*AgentVersionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AgentVersionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AgentVersionId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAgentVersionIDInsensitively parses 'input' case-insensitively into a AgentVersionId
// note: this method should only be used for API response data and not user input
func ParseAgentVersionIDInsensitively(input string) (*AgentVersionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AgentVersionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AgentVersionId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AgentVersionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.OsTypeName, ok = input.Parsed["osTypeName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "osTypeName", input)
	}

	if id.AgentVersionName, ok = input.Parsed["agentVersionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "agentVersionName", input)
	}

	return nil
}

// ValidateAgentVersionID checks that 'input' can be parsed as a Agent Version ID
func ValidateAgentVersionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAgentVersionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Agent Version ID
func (id AgentVersionId) ID() string {
	fmtString := "/providers/Microsoft.HybridCompute/osType/%s/agentVersions/%s"
	return fmt.Sprintf(fmtString, id.OsTypeName, id.AgentVersionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Agent Version ID
func (id AgentVersionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftHybridCompute", "Microsoft.HybridCompute", "Microsoft.HybridCompute"),
		resourceids.StaticSegment("staticOsType", "osType", "osType"),
		resourceids.UserSpecifiedSegment("osTypeName", "osTypeValue"),
		resourceids.StaticSegment("staticAgentVersions", "agentVersions", "agentVersions"),
		resourceids.UserSpecifiedSegment("agentVersionName", "agentVersionValue"),
	}
}

// String returns a human-readable description of this Agent Version ID
func (id AgentVersionId) String() string {
	components := []string{
		fmt.Sprintf("Os Type Name: %q", id.OsTypeName),
		fmt.Sprintf("Agent Version Name: %q", id.AgentVersionName),
	}
	return fmt.Sprintf("Agent Version (%s)", strings.Join(components, "\n"))
}
