package nginxdeploymentapikeyresponses

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ApiKeyId{})
}

var _ resourceids.ResourceId = &ApiKeyId{}

// ApiKeyId is a struct representing the Resource ID for a Api Key
type ApiKeyId struct {
	SubscriptionId      string
	ResourceGroupName   string
	NginxDeploymentName string
	ApiKeyName          string
}

// NewApiKeyID returns a new ApiKeyId struct
func NewApiKeyID(subscriptionId string, resourceGroupName string, nginxDeploymentName string, apiKeyName string) ApiKeyId {
	return ApiKeyId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		NginxDeploymentName: nginxDeploymentName,
		ApiKeyName:          apiKeyName,
	}
}

// ParseApiKeyID parses 'input' into a ApiKeyId
func ParseApiKeyID(input string) (*ApiKeyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ApiKeyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ApiKeyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseApiKeyIDInsensitively parses 'input' case-insensitively into a ApiKeyId
// note: this method should only be used for API response data and not user input
func ParseApiKeyIDInsensitively(input string) (*ApiKeyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ApiKeyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ApiKeyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ApiKeyId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.NginxDeploymentName, ok = input.Parsed["nginxDeploymentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "nginxDeploymentName", input)
	}

	if id.ApiKeyName, ok = input.Parsed["apiKeyName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "apiKeyName", input)
	}

	return nil
}

// ValidateApiKeyID checks that 'input' can be parsed as a Api Key ID
func ValidateApiKeyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseApiKeyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Api Key ID
func (id ApiKeyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Nginx.NginxPlus/nginxDeployments/%s/apiKeys/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NginxDeploymentName, id.ApiKeyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Api Key ID
func (id ApiKeyId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticNginxNginxPlus", "Nginx.NginxPlus", "Nginx.NginxPlus"),
		resourceids.StaticSegment("staticNginxDeployments", "nginxDeployments", "nginxDeployments"),
		resourceids.UserSpecifiedSegment("nginxDeploymentName", "nginxDeploymentName"),
		resourceids.StaticSegment("staticApiKeys", "apiKeys", "apiKeys"),
		resourceids.UserSpecifiedSegment("apiKeyName", "apiKeyName"),
	}
}

// String returns a human-readable description of this Api Key ID
func (id ApiKeyId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Nginx Deployment Name: %q", id.NginxDeploymentName),
		fmt.Sprintf("Api Key Name: %q", id.ApiKeyName),
	}
	return fmt.Sprintf("Api Key (%s)", strings.Join(components, "\n"))
}
