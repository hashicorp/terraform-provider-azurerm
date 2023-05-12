package contentkeypolicies

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ContentKeyPolicyId{}

// ContentKeyPolicyId is a struct representing the Resource ID for a Content Key Policy
type ContentKeyPolicyId struct {
	SubscriptionId       string
	ResourceGroupName    string
	MediaServiceName     string
	ContentKeyPolicyName string
}

// NewContentKeyPolicyID returns a new ContentKeyPolicyId struct
func NewContentKeyPolicyID(subscriptionId string, resourceGroupName string, mediaServiceName string, contentKeyPolicyName string) ContentKeyPolicyId {
	return ContentKeyPolicyId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		MediaServiceName:     mediaServiceName,
		ContentKeyPolicyName: contentKeyPolicyName,
	}
}

// ParseContentKeyPolicyID parses 'input' into a ContentKeyPolicyId
func ParseContentKeyPolicyID(input string) (*ContentKeyPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(ContentKeyPolicyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ContentKeyPolicyId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.MediaServiceName, ok = parsed.Parsed["mediaServiceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "mediaServiceName", *parsed)
	}

	if id.ContentKeyPolicyName, ok = parsed.Parsed["contentKeyPolicyName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "contentKeyPolicyName", *parsed)
	}

	return &id, nil
}

// ParseContentKeyPolicyIDInsensitively parses 'input' case-insensitively into a ContentKeyPolicyId
// note: this method should only be used for API response data and not user input
func ParseContentKeyPolicyIDInsensitively(input string) (*ContentKeyPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(ContentKeyPolicyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ContentKeyPolicyId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.MediaServiceName, ok = parsed.Parsed["mediaServiceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "mediaServiceName", *parsed)
	}

	if id.ContentKeyPolicyName, ok = parsed.Parsed["contentKeyPolicyName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "contentKeyPolicyName", *parsed)
	}

	return &id, nil
}

// ValidateContentKeyPolicyID checks that 'input' can be parsed as a Content Key Policy ID
func ValidateContentKeyPolicyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseContentKeyPolicyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Content Key Policy ID
func (id ContentKeyPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Media/mediaServices/%s/contentKeyPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.MediaServiceName, id.ContentKeyPolicyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Content Key Policy ID
func (id ContentKeyPolicyId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMedia", "Microsoft.Media", "Microsoft.Media"),
		resourceids.StaticSegment("staticMediaServices", "mediaServices", "mediaServices"),
		resourceids.UserSpecifiedSegment("mediaServiceName", "mediaServiceValue"),
		resourceids.StaticSegment("staticContentKeyPolicies", "contentKeyPolicies", "contentKeyPolicies"),
		resourceids.UserSpecifiedSegment("contentKeyPolicyName", "contentKeyPolicyValue"),
	}
}

// String returns a human-readable description of this Content Key Policy ID
func (id ContentKeyPolicyId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Media Service Name: %q", id.MediaServiceName),
		fmt.Sprintf("Content Key Policy Name: %q", id.ContentKeyPolicyName),
	}
	return fmt.Sprintf("Content Key Policy (%s)", strings.Join(components, "\n"))
}
