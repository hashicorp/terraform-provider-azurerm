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
	recaser.RegisterResourceId(&KeyId{})
}

var _ resourceids.ResourceId = &KeyId{}

// KeyId is a struct representing the Resource ID for a Key
type KeyId struct {
	SubscriptionId    string
	ResourceGroupName string
	SiteName          string
	FunctionName      string
	KeyName           string
}

// NewKeyID returns a new KeyId struct
func NewKeyID(subscriptionId string, resourceGroupName string, siteName string, functionName string, keyName string) KeyId {
	return KeyId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SiteName:          siteName,
		FunctionName:      functionName,
		KeyName:           keyName,
	}
}

// ParseKeyID parses 'input' into a KeyId
func ParseKeyID(input string) (*KeyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&KeyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := KeyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseKeyIDInsensitively parses 'input' case-insensitively into a KeyId
// note: this method should only be used for API response data and not user input
func ParseKeyIDInsensitively(input string) (*KeyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&KeyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := KeyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *KeyId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.KeyName, ok = input.Parsed["keyName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "keyName", input)
	}

	return nil
}

// ValidateKeyID checks that 'input' can be parsed as a Key ID
func ValidateKeyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseKeyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Key ID
func (id KeyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/functions/%s/keys/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SiteName, id.FunctionName, id.KeyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Key ID
func (id KeyId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticKeys", "keys", "keys"),
		resourceids.UserSpecifiedSegment("keyName", "keyName"),
	}
}

// String returns a human-readable description of this Key ID
func (id KeyId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
		fmt.Sprintf("Function Name: %q", id.FunctionName),
		fmt.Sprintf("Key Name: %q", id.KeyName),
	}
	return fmt.Sprintf("Key (%s)", strings.Join(components, "\n"))
}
