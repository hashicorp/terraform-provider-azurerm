package managementlocks

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ResourceLockId{}

// ResourceLockId is a struct representing the Resource ID for a Resource Lock
type ResourceLockId struct {
	SubscriptionId     string
	ResourceGroupName  string
	ProviderName       string
	ParentResourcePath string
	ResourceType       string
	ResourceName       string
	LockName           string
}

// NewResourceLockID returns a new ResourceLockId struct
func NewResourceLockID(subscriptionId string, resourceGroupName string, providerName string, parentResourcePath string, resourceType string, resourceName string, lockName string) ResourceLockId {
	return ResourceLockId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		ProviderName:       providerName,
		ParentResourcePath: parentResourcePath,
		ResourceType:       resourceType,
		ResourceName:       resourceName,
		LockName:           lockName,
	}
}

// ParseResourceLockID parses 'input' into a ResourceLockId
func ParseResourceLockID(input string) (*ResourceLockId, error) {
	parser := resourceids.NewParserFromResourceIdType(ResourceLockId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ResourceLockId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ProviderName, ok = parsed.Parsed["providerName"]; !ok {
		return nil, fmt.Errorf("the segment 'providerName' was not found in the resource id %q", input)
	}

	if id.ParentResourcePath, ok = parsed.Parsed["parentResourcePath"]; !ok {
		return nil, fmt.Errorf("the segment 'parentResourcePath' was not found in the resource id %q", input)
	}

	if id.ResourceType, ok = parsed.Parsed["resourceType"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceType' was not found in the resource id %q", input)
	}

	if id.ResourceName, ok = parsed.Parsed["resourceName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceName' was not found in the resource id %q", input)
	}

	if id.LockName, ok = parsed.Parsed["lockName"]; !ok {
		return nil, fmt.Errorf("the segment 'lockName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseResourceLockIDInsensitively parses 'input' case-insensitively into a ResourceLockId
// note: this method should only be used for API response data and not user input
func ParseResourceLockIDInsensitively(input string) (*ResourceLockId, error) {
	parser := resourceids.NewParserFromResourceIdType(ResourceLockId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ResourceLockId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ProviderName, ok = parsed.Parsed["providerName"]; !ok {
		return nil, fmt.Errorf("the segment 'providerName' was not found in the resource id %q", input)
	}

	if id.ParentResourcePath, ok = parsed.Parsed["parentResourcePath"]; !ok {
		return nil, fmt.Errorf("the segment 'parentResourcePath' was not found in the resource id %q", input)
	}

	if id.ResourceType, ok = parsed.Parsed["resourceType"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceType' was not found in the resource id %q", input)
	}

	if id.ResourceName, ok = parsed.Parsed["resourceName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceName' was not found in the resource id %q", input)
	}

	if id.LockName, ok = parsed.Parsed["lockName"]; !ok {
		return nil, fmt.Errorf("the segment 'lockName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateResourceLockID checks that 'input' can be parsed as a Resource Lock ID
func ValidateResourceLockID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseResourceLockID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Resource Lock ID
func (id ResourceLockId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/%s/%s/%s/%s/providers/Microsoft.Authorization/locks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ProviderName, id.ParentResourcePath, id.ResourceType, id.ResourceName, id.LockName)
}

// Segments returns a slice of Resource ID Segments which comprise this Resource Lock ID
func (id ResourceLockId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.UserSpecifiedSegment("providerName", "providerValue"),
		resourceids.UserSpecifiedSegment("parentResourcePath", "parentResourcePathValue"),
		resourceids.UserSpecifiedSegment("resourceType", "resourceTypeValue"),
		resourceids.UserSpecifiedSegment("resourceName", "resourceValue"),
		resourceids.StaticSegment("staticProviders2", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAuthorization", "Microsoft.Authorization", "Microsoft.Authorization"),
		resourceids.StaticSegment("staticLocks", "locks", "locks"),
		resourceids.UserSpecifiedSegment("lockName", "lockValue"),
	}
}

// String returns a human-readable description of this Resource Lock ID
func (id ResourceLockId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Provider Name: %q", id.ProviderName),
		fmt.Sprintf("Parent Resource Path: %q", id.ParentResourcePath),
		fmt.Sprintf("Resource Type: %q", id.ResourceType),
		fmt.Sprintf("Resource Name: %q", id.ResourceName),
		fmt.Sprintf("Lock Name: %q", id.LockName),
	}
	return fmt.Sprintf("Resource Lock (%s)", strings.Join(components, "\n"))
}
