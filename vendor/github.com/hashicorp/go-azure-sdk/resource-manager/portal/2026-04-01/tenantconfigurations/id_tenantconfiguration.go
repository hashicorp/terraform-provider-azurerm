package tenantconfigurations

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&TenantConfigurationId{})
}

var _ resourceids.ResourceId = &TenantConfigurationId{}

// TenantConfigurationId is a struct representing the Resource ID for a Tenant Configuration
type TenantConfigurationId struct {
	TenantConfigurationName string
}

// NewTenantConfigurationID returns a new TenantConfigurationId struct
func NewTenantConfigurationID(tenantConfigurationName string) TenantConfigurationId {
	return TenantConfigurationId{
		TenantConfigurationName: tenantConfigurationName,
	}
}

// ParseTenantConfigurationID parses 'input' into a TenantConfigurationId
func ParseTenantConfigurationID(input string) (*TenantConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TenantConfigurationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TenantConfigurationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseTenantConfigurationIDInsensitively parses 'input' case-insensitively into a TenantConfigurationId
// note: this method should only be used for API response data and not user input
func ParseTenantConfigurationIDInsensitively(input string) (*TenantConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TenantConfigurationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TenantConfigurationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *TenantConfigurationId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.TenantConfigurationName, ok = input.Parsed["tenantConfigurationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "tenantConfigurationName", input)
	}

	return nil
}

// ValidateTenantConfigurationID checks that 'input' can be parsed as a Tenant Configuration ID
func ValidateTenantConfigurationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseTenantConfigurationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Tenant Configuration ID
func (id TenantConfigurationId) ID() string {
	fmtString := "/providers/Microsoft.Portal/tenantConfigurations/%s"
	return fmt.Sprintf(fmtString, id.TenantConfigurationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Tenant Configuration ID
func (id TenantConfigurationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftPortal", "Microsoft.Portal", "Microsoft.Portal"),
		resourceids.StaticSegment("staticTenantConfigurations", "tenantConfigurations", "tenantConfigurations"),
		resourceids.UserSpecifiedSegment("tenantConfigurationName", "tenantConfigurationName"),
	}
}

// String returns a human-readable description of this Tenant Configuration ID
func (id TenantConfigurationId) String() string {
	components := []string{
		fmt.Sprintf("Tenant Configuration Name: %q", id.TenantConfigurationName),
	}
	return fmt.Sprintf("Tenant Configuration (%s)", strings.Join(components, "\n"))
}
