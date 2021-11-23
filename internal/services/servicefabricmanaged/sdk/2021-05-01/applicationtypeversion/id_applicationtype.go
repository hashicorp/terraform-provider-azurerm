package applicationtypeversion

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ApplicationTypeId{}

// ApplicationTypeId is a struct representing the Resource ID for a Application Type
type ApplicationTypeId struct {
	SubscriptionId      string
	ResourceGroupName   string
	ClusterName         string
	ApplicationTypeName string
}

// NewApplicationTypeID returns a new ApplicationTypeId struct
func NewApplicationTypeID(subscriptionId string, resourceGroupName string, clusterName string, applicationTypeName string) ApplicationTypeId {
	return ApplicationTypeId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		ClusterName:         clusterName,
		ApplicationTypeName: applicationTypeName,
	}
}

// ParseApplicationTypeID parses 'input' into a ApplicationTypeId
func ParseApplicationTypeID(input string) (*ApplicationTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(ApplicationTypeId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ApplicationTypeId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ClusterName, ok = parsed.Parsed["clusterName"]; !ok {
		return nil, fmt.Errorf("the segment 'clusterName' was not found in the resource id %q", input)
	}

	if id.ApplicationTypeName, ok = parsed.Parsed["applicationTypeName"]; !ok {
		return nil, fmt.Errorf("the segment 'applicationTypeName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseApplicationTypeIDInsensitively parses 'input' case-insensitively into a ApplicationTypeId
// note: this method should only be used for API response data and not user input
func ParseApplicationTypeIDInsensitively(input string) (*ApplicationTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(ApplicationTypeId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ApplicationTypeId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ClusterName, ok = parsed.Parsed["clusterName"]; !ok {
		return nil, fmt.Errorf("the segment 'clusterName' was not found in the resource id %q", input)
	}

	if id.ApplicationTypeName, ok = parsed.Parsed["applicationTypeName"]; !ok {
		return nil, fmt.Errorf("the segment 'applicationTypeName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateApplicationTypeID checks that 'input' can be parsed as a Application Type ID
func ValidateApplicationTypeID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseApplicationTypeID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Application Type ID
func (id ApplicationTypeId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ServiceFabric/managedClusters/%s/applicationTypes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ClusterName, id.ApplicationTypeName)
}

// Segments returns a slice of Resource ID Segments which comprise this Application Type ID
func (id ApplicationTypeId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("subscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("resourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("providers", "providers", "providers"),
		resourceids.ResourceProviderSegment("microsoftServiceFabric", "Microsoft.ServiceFabric", "Microsoft.ServiceFabric"),
		resourceids.StaticSegment("managedClusters", "managedClusters", "managedClusters"),
		resourceids.UserSpecifiedSegment("clusterName", "clusterValue"),
		resourceids.StaticSegment("applicationTypes", "applicationTypes", "applicationTypes"),
		resourceids.UserSpecifiedSegment("applicationTypeName", "applicationTypeValue"),
	}
}

// String returns a human-readable description of this Application Type ID
func (id ApplicationTypeId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Cluster Name: %q", id.ClusterName),
		fmt.Sprintf("Application Type Name: %q", id.ApplicationTypeName),
	}
	return fmt.Sprintf("Application Type (%s)", strings.Join(components, "\n"))
}
