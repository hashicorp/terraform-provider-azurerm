package workspace

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

func init() {
	recaser.RegisterResourceId(&WorkspaceId{})
}

var _ resourceids.ResourceId = &WorkspaceId{}

type WorkspaceId struct {
	SubscriptionId string
	ResourceGroup  string
	ServiceName    string
	WorkspaceId    string
}

func NewWorkspaceID(subscriptionId, resourceGroup, serviceName, workspaceId string) WorkspaceId {
	return WorkspaceId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ServiceName:    serviceName,
		WorkspaceId:    workspaceId,
	}
}

func ParseWorkspaceID(input string) (*WorkspaceId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := WorkspaceId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.ServiceName, err = id.PopSegment("service"); err != nil {
		return nil, err
	}
	if resourceId.WorkspaceId, err = id.PopSegment("workspaces"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

func (id *WorkspaceId) FromParseResult(input resourceids.ParseResult) error {
	if subscriptionId, ok := input.Parsed["subscriptionId"]; ok {
		id.SubscriptionId = subscriptionId
	}
	if resourceGroup, ok := input.Parsed["resourceGroup"]; ok {
		id.ResourceGroup = resourceGroup
	}
	if serviceName, ok := input.Parsed["service"]; ok {
		id.ServiceName = serviceName
	}
	if workspaceId, ok := input.Parsed["workspaces"]; ok {
		id.WorkspaceId = workspaceId
	}
	return nil
}

func ValidateWorkspaceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseWorkspaceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

func (id WorkspaceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/workspaces/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ServiceName, id.WorkspaceId)
}

func (id WorkspaceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroup", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.StaticSegment("staticMicrosoftApiManagement", "Microsoft.ApiManagement", "Microsoft.ApiManagement"),
		resourceids.StaticSegment("staticService", "service", "service"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceValue"),
		resourceids.StaticSegment("staticWorkspaces", "workspaces", "workspaces"),
		resourceids.UserSpecifiedSegment("workspaceId", "workspaceValue"),
	}
}

func (id WorkspaceId) String() string {
	components := []string{
		fmt.Sprintf("SubscriptionId: %q", id.SubscriptionId),
		fmt.Sprintf("ResourceGroup: %q", id.ResourceGroup),
		fmt.Sprintf("ServiceName: %q", id.ServiceName),
		fmt.Sprintf("WorkspaceId: %q", id.WorkspaceId),
	}
	return fmt.Sprintf("Workspace (%s)", strings.Join(components, " / "))
}
