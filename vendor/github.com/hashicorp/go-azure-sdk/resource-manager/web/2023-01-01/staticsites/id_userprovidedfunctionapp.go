package staticsites

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&UserProvidedFunctionAppId{})
}

var _ resourceids.ResourceId = &UserProvidedFunctionAppId{}

// UserProvidedFunctionAppId is a struct representing the Resource ID for a User Provided Function App
type UserProvidedFunctionAppId struct {
	SubscriptionId              string
	ResourceGroupName           string
	StaticSiteName              string
	UserProvidedFunctionAppName string
}

// NewUserProvidedFunctionAppID returns a new UserProvidedFunctionAppId struct
func NewUserProvidedFunctionAppID(subscriptionId string, resourceGroupName string, staticSiteName string, userProvidedFunctionAppName string) UserProvidedFunctionAppId {
	return UserProvidedFunctionAppId{
		SubscriptionId:              subscriptionId,
		ResourceGroupName:           resourceGroupName,
		StaticSiteName:              staticSiteName,
		UserProvidedFunctionAppName: userProvidedFunctionAppName,
	}
}

// ParseUserProvidedFunctionAppID parses 'input' into a UserProvidedFunctionAppId
func ParseUserProvidedFunctionAppID(input string) (*UserProvidedFunctionAppId, error) {
	parser := resourceids.NewParserFromResourceIdType(&UserProvidedFunctionAppId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := UserProvidedFunctionAppId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseUserProvidedFunctionAppIDInsensitively parses 'input' case-insensitively into a UserProvidedFunctionAppId
// note: this method should only be used for API response data and not user input
func ParseUserProvidedFunctionAppIDInsensitively(input string) (*UserProvidedFunctionAppId, error) {
	parser := resourceids.NewParserFromResourceIdType(&UserProvidedFunctionAppId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := UserProvidedFunctionAppId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *UserProvidedFunctionAppId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.StaticSiteName, ok = input.Parsed["staticSiteName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "staticSiteName", input)
	}

	if id.UserProvidedFunctionAppName, ok = input.Parsed["userProvidedFunctionAppName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "userProvidedFunctionAppName", input)
	}

	return nil
}

// ValidateUserProvidedFunctionAppID checks that 'input' can be parsed as a User Provided Function App ID
func ValidateUserProvidedFunctionAppID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseUserProvidedFunctionAppID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted User Provided Function App ID
func (id UserProvidedFunctionAppId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/staticSites/%s/userProvidedFunctionApps/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.StaticSiteName, id.UserProvidedFunctionAppName)
}

// Segments returns a slice of Resource ID Segments which comprise this User Provided Function App ID
func (id UserProvidedFunctionAppId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticStaticSites", "staticSites", "staticSites"),
		resourceids.UserSpecifiedSegment("staticSiteName", "staticSiteName"),
		resourceids.StaticSegment("staticUserProvidedFunctionApps", "userProvidedFunctionApps", "userProvidedFunctionApps"),
		resourceids.UserSpecifiedSegment("userProvidedFunctionAppName", "userProvidedFunctionAppName"),
	}
}

// String returns a human-readable description of this User Provided Function App ID
func (id UserProvidedFunctionAppId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Static Site Name: %q", id.StaticSiteName),
		fmt.Sprintf("User Provided Function App Name: %q", id.UserProvidedFunctionAppName),
	}
	return fmt.Sprintf("User Provided Function App (%s)", strings.Join(components, "\n"))
}
