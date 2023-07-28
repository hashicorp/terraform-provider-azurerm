package modelcontainer

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = RegistryModelId{}

// RegistryModelId is a struct representing the Resource ID for a Registry Model
type RegistryModelId struct {
	SubscriptionId    string
	ResourceGroupName string
	RegistryName      string
	ModelName         string
}

// NewRegistryModelID returns a new RegistryModelId struct
func NewRegistryModelID(subscriptionId string, resourceGroupName string, registryName string, modelName string) RegistryModelId {
	return RegistryModelId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		RegistryName:      registryName,
		ModelName:         modelName,
	}
}

// ParseRegistryModelID parses 'input' into a RegistryModelId
func ParseRegistryModelID(input string) (*RegistryModelId, error) {
	parser := resourceids.NewParserFromResourceIdType(RegistryModelId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RegistryModelId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.RegistryName, ok = parsed.Parsed["registryName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "registryName", *parsed)
	}

	if id.ModelName, ok = parsed.Parsed["modelName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "modelName", *parsed)
	}

	return &id, nil
}

// ParseRegistryModelIDInsensitively parses 'input' case-insensitively into a RegistryModelId
// note: this method should only be used for API response data and not user input
func ParseRegistryModelIDInsensitively(input string) (*RegistryModelId, error) {
	parser := resourceids.NewParserFromResourceIdType(RegistryModelId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RegistryModelId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.RegistryName, ok = parsed.Parsed["registryName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "registryName", *parsed)
	}

	if id.ModelName, ok = parsed.Parsed["modelName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "modelName", *parsed)
	}

	return &id, nil
}

// ValidateRegistryModelID checks that 'input' can be parsed as a Registry Model ID
func ValidateRegistryModelID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRegistryModelID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Registry Model ID
func (id RegistryModelId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.MachineLearningServices/registries/%s/models/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.RegistryName, id.ModelName)
}

// Segments returns a slice of Resource ID Segments which comprise this Registry Model ID
func (id RegistryModelId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMachineLearningServices", "Microsoft.MachineLearningServices", "Microsoft.MachineLearningServices"),
		resourceids.StaticSegment("staticRegistries", "registries", "registries"),
		resourceids.UserSpecifiedSegment("registryName", "registryValue"),
		resourceids.StaticSegment("staticModels", "models", "models"),
		resourceids.UserSpecifiedSegment("modelName", "modelValue"),
	}
}

// String returns a human-readable description of this Registry Model ID
func (id RegistryModelId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Registry Name: %q", id.RegistryName),
		fmt.Sprintf("Model Name: %q", id.ModelName),
	}
	return fmt.Sprintf("Registry Model (%s)", strings.Join(components, "\n"))
}
