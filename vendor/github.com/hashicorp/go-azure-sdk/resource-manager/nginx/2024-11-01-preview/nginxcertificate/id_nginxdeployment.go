package nginxcertificate

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&NginxDeploymentId{})
}

var _ resourceids.ResourceId = &NginxDeploymentId{}

// NginxDeploymentId is a struct representing the Resource ID for a Nginx Deployment
type NginxDeploymentId struct {
	SubscriptionId      string
	ResourceGroupName   string
	NginxDeploymentName string
}

// NewNginxDeploymentID returns a new NginxDeploymentId struct
func NewNginxDeploymentID(subscriptionId string, resourceGroupName string, nginxDeploymentName string) NginxDeploymentId {
	return NginxDeploymentId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		NginxDeploymentName: nginxDeploymentName,
	}
}

// ParseNginxDeploymentID parses 'input' into a NginxDeploymentId
func ParseNginxDeploymentID(input string) (*NginxDeploymentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NginxDeploymentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NginxDeploymentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseNginxDeploymentIDInsensitively parses 'input' case-insensitively into a NginxDeploymentId
// note: this method should only be used for API response data and not user input
func ParseNginxDeploymentIDInsensitively(input string) (*NginxDeploymentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NginxDeploymentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NginxDeploymentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *NginxDeploymentId) FromParseResult(input resourceids.ParseResult) error {
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

	return nil
}

// ValidateNginxDeploymentID checks that 'input' can be parsed as a Nginx Deployment ID
func ValidateNginxDeploymentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseNginxDeploymentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Nginx Deployment ID
func (id NginxDeploymentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Nginx.NginxPlus/nginxDeployments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NginxDeploymentName)
}

// Segments returns a slice of Resource ID Segments which comprise this Nginx Deployment ID
func (id NginxDeploymentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticNginxNginxPlus", "Nginx.NginxPlus", "Nginx.NginxPlus"),
		resourceids.StaticSegment("staticNginxDeployments", "nginxDeployments", "nginxDeployments"),
		resourceids.UserSpecifiedSegment("nginxDeploymentName", "nginxDeploymentName"),
	}
}

// String returns a human-readable description of this Nginx Deployment ID
func (id NginxDeploymentId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Nginx Deployment Name: %q", id.NginxDeploymentName),
	}
	return fmt.Sprintf("Nginx Deployment (%s)", strings.Join(components, "\n"))
}
