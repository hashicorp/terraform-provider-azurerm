package managedclusterversion

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ManagedClusterVersionId{}

// ManagedClusterVersionId is a struct representing the Resource ID for a Managed Cluster Version
type ManagedClusterVersionId struct {
	SubscriptionId string
	Location       string
	ClusterVersion string
}

// NewManagedClusterVersionID returns a new ManagedClusterVersionId struct
func NewManagedClusterVersionID(subscriptionId string, location string, clusterVersion string) ManagedClusterVersionId {
	return ManagedClusterVersionId{
		SubscriptionId: subscriptionId,
		Location:       location,
		ClusterVersion: clusterVersion,
	}
}

// ParseManagedClusterVersionID parses 'input' into a ManagedClusterVersionId
func ParseManagedClusterVersionID(input string) (*ManagedClusterVersionId, error) {
	parser := resourceids.NewParserFromResourceIdType(ManagedClusterVersionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ManagedClusterVersionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.Location, ok = parsed.Parsed["location"]; !ok {
		return nil, fmt.Errorf("the segment 'location' was not found in the resource id %q", input)
	}

	if id.ClusterVersion, ok = parsed.Parsed["clusterVersion"]; !ok {
		return nil, fmt.Errorf("the segment 'clusterVersion' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseManagedClusterVersionIDInsensitively parses 'input' case-insensitively into a ManagedClusterVersionId
// note: this method should only be used for API response data and not user input
func ParseManagedClusterVersionIDInsensitively(input string) (*ManagedClusterVersionId, error) {
	parser := resourceids.NewParserFromResourceIdType(ManagedClusterVersionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ManagedClusterVersionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.Location, ok = parsed.Parsed["location"]; !ok {
		return nil, fmt.Errorf("the segment 'location' was not found in the resource id %q", input)
	}

	if id.ClusterVersion, ok = parsed.Parsed["clusterVersion"]; !ok {
		return nil, fmt.Errorf("the segment 'clusterVersion' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateManagedClusterVersionID checks that 'input' can be parsed as a Managed Cluster Version ID
func ValidateManagedClusterVersionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseManagedClusterVersionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Managed Cluster Version ID
func (id ManagedClusterVersionId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.ServiceFabric/locations/%s/managedClusterVersions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.Location, id.ClusterVersion)
}

// Segments returns a slice of Resource ID Segments which comprise this Managed Cluster Version ID
func (id ManagedClusterVersionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("subscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("providers", "providers", "providers"),
		resourceids.ResourceProviderSegment("microsoftServiceFabric", "Microsoft.ServiceFabric", "Microsoft.ServiceFabric"),
		resourceids.StaticSegment("locations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("location", "locationValue"),
		resourceids.StaticSegment("managedClusterVersions", "managedClusterVersions", "managedClusterVersions"),
		resourceids.UserSpecifiedSegment("clusterVersion", "clusterVersionValue"),
	}
}

// String returns a human-readable description of this Managed Cluster Version ID
func (id ManagedClusterVersionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location: %q", id.Location),
		fmt.Sprintf("Cluster Version: %q", id.ClusterVersion),
	}
	return fmt.Sprintf("Managed Cluster Version (%s)", strings.Join(components, "\n"))
}
