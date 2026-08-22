package webapps

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&FunctionId{})
}

var _ resourceids.ResourceId = &FunctionId{}

// FunctionId is a struct representing the Resource ID for a Function
type FunctionId struct {
	SubscriptionId    string
	ResourceGroupName string
	SiteName          string
	FunctionName      string
}

// NewFunctionID returns a new FunctionId struct
func NewFunctionID(subscriptionId string, resourceGroupName string, siteName string, functionName string) FunctionId {
	return FunctionId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SiteName:          siteName,
		FunctionName:      functionName,
	}
}

// ParseFunctionID parses 'input' into a FunctionId
func ParseFunctionID(input string) (*FunctionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&FunctionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := FunctionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseFunctionIDInsensitively parses 'input' case-insensitively into a FunctionId
// note: this method should only be used for API response data and not user input
func ParseFunctionIDInsensitively(input string) (*FunctionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&FunctionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := FunctionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *FunctionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SiteName, ok = input.Parsed["siteName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "siteName", input)
	}

	if id.FunctionName, ok = input.Parsed["functionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "functionName", input)
	}

	return nil
}

// ValidateFunctionID checks that 'input' can be parsed as a Function ID
func ValidateFunctionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseFunctionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Function ID
func (id FunctionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/functions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SiteName, id.FunctionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Function ID
func (id FunctionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticSites", "sites", "sites"),
		resourceids.UserSpecifiedSegment("siteName", "siteName"),
		resourceids.StaticSegment("staticFunctions", "functions", "functions"),
		resourceids.UserSpecifiedSegment("functionName", "functionName"),
	}
}

// String returns a human-readable description of this Function ID
func (id FunctionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
		fmt.Sprintf("Function Name: %q", id.FunctionName),
	}
	return fmt.Sprintf("Function (%s)", strings.Join(components, "\n"))
}
