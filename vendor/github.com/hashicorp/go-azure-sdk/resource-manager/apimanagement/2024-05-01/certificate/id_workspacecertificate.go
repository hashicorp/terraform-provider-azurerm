package certificate

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&WorkspaceCertificateId{})
}

var _ resourceids.ResourceId = &WorkspaceCertificateId{}

// WorkspaceCertificateId is a struct representing the Resource ID for a Workspace Certificate
type WorkspaceCertificateId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	WorkspaceId       string
	CertificateId     string
}

// NewWorkspaceCertificateID returns a new WorkspaceCertificateId struct
func NewWorkspaceCertificateID(subscriptionId string, resourceGroupName string, serviceName string, workspaceId string, certificateId string) WorkspaceCertificateId {
	return WorkspaceCertificateId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		WorkspaceId:       workspaceId,
		CertificateId:     certificateId,
	}
}

// ParseWorkspaceCertificateID parses 'input' into a WorkspaceCertificateId
func ParseWorkspaceCertificateID(input string) (*WorkspaceCertificateId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WorkspaceCertificateId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WorkspaceCertificateId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseWorkspaceCertificateIDInsensitively parses 'input' case-insensitively into a WorkspaceCertificateId
// note: this method should only be used for API response data and not user input
func ParseWorkspaceCertificateIDInsensitively(input string) (*WorkspaceCertificateId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WorkspaceCertificateId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WorkspaceCertificateId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *WorkspaceCertificateId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ServiceName, ok = input.Parsed["serviceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "serviceName", input)
	}

	if id.WorkspaceId, ok = input.Parsed["workspaceId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "workspaceId", input)
	}

	if id.CertificateId, ok = input.Parsed["certificateId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "certificateId", input)
	}

	return nil
}

// ValidateWorkspaceCertificateID checks that 'input' can be parsed as a Workspace Certificate ID
func ValidateWorkspaceCertificateID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseWorkspaceCertificateID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Workspace Certificate ID
func (id WorkspaceCertificateId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/workspaces/%s/certificates/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.WorkspaceId, id.CertificateId)
}

// Segments returns a slice of Resource ID Segments which comprise this Workspace Certificate ID
func (id WorkspaceCertificateId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiManagement", "Microsoft.ApiManagement", "Microsoft.ApiManagement"),
		resourceids.StaticSegment("staticService", "service", "service"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceName"),
		resourceids.StaticSegment("staticWorkspaces", "workspaces", "workspaces"),
		resourceids.UserSpecifiedSegment("workspaceId", "workspaceId"),
		resourceids.StaticSegment("staticCertificates", "certificates", "certificates"),
		resourceids.UserSpecifiedSegment("certificateId", "certificateId"),
	}
}

// String returns a human-readable description of this Workspace Certificate ID
func (id WorkspaceCertificateId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Workspace: %q", id.WorkspaceId),
		fmt.Sprintf("Certificate: %q", id.CertificateId),
	}
	return fmt.Sprintf("Workspace Certificate (%s)", strings.Join(components, "\n"))
}
