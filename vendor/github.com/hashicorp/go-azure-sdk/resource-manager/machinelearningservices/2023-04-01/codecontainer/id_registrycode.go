package codecontainer

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = RegistryCodeId{}

// RegistryCodeId is a struct representing the Resource ID for a Registry Code
type RegistryCodeId struct {
	SubscriptionId    string
	ResourceGroupName string
	RegistryName      string
	CodeName          string
}

// NewRegistryCodeID returns a new RegistryCodeId struct
func NewRegistryCodeID(subscriptionId string, resourceGroupName string, registryName string, codeName string) RegistryCodeId {
	return RegistryCodeId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		RegistryName:      registryName,
		CodeName:          codeName,
	}
}

// ParseRegistryCodeID parses 'input' into a RegistryCodeId
func ParseRegistryCodeID(input string) (*RegistryCodeId, error) {
	parser := resourceids.NewParserFromResourceIdType(RegistryCodeId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RegistryCodeId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.RegistryName, ok = parsed.Parsed["registryName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "registryName", *parsed)
	}

	if id.CodeName, ok = parsed.Parsed["codeName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "codeName", *parsed)
	}

	return &id, nil
}

// ParseRegistryCodeIDInsensitively parses 'input' case-insensitively into a RegistryCodeId
// note: this method should only be used for API response data and not user input
func ParseRegistryCodeIDInsensitively(input string) (*RegistryCodeId, error) {
	parser := resourceids.NewParserFromResourceIdType(RegistryCodeId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RegistryCodeId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.RegistryName, ok = parsed.Parsed["registryName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "registryName", *parsed)
	}

	if id.CodeName, ok = parsed.Parsed["codeName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "codeName", *parsed)
	}

	return &id, nil
}

// ValidateRegistryCodeID checks that 'input' can be parsed as a Registry Code ID
func ValidateRegistryCodeID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRegistryCodeID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Registry Code ID
func (id RegistryCodeId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.MachineLearningServices/registries/%s/codes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.RegistryName, id.CodeName)
}

// Segments returns a slice of Resource ID Segments which comprise this Registry Code ID
func (id RegistryCodeId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMachineLearningServices", "Microsoft.MachineLearningServices", "Microsoft.MachineLearningServices"),
		resourceids.StaticSegment("staticRegistries", "registries", "registries"),
		resourceids.UserSpecifiedSegment("registryName", "registryValue"),
		resourceids.StaticSegment("staticCodes", "codes", "codes"),
		resourceids.UserSpecifiedSegment("codeName", "codeValue"),
	}
}

// String returns a human-readable description of this Registry Code ID
func (id RegistryCodeId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Registry Name: %q", id.RegistryName),
		fmt.Sprintf("Code Name: %q", id.CodeName),
	}
	return fmt.Sprintf("Registry Code (%s)", strings.Join(components, "\n"))
}
