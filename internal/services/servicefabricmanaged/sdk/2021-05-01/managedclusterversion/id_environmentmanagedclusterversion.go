package managedclusterversion

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = EnvironmentManagedClusterVersionId{}

// EnvironmentManagedClusterVersionId is a struct representing the Resource ID for a Environment Managed Cluster Version
type EnvironmentManagedClusterVersionId struct {
	SubscriptionId string
	Location       string
	Environment    Environment
	ClusterVersion string
}

// NewEnvironmentManagedClusterVersionID returns a new EnvironmentManagedClusterVersionId struct
func NewEnvironmentManagedClusterVersionID(subscriptionId string, location string, environment Environment, clusterVersion string) EnvironmentManagedClusterVersionId {
	return EnvironmentManagedClusterVersionId{
		SubscriptionId: subscriptionId,
		Location:       location,
		Environment:    environment,
		ClusterVersion: clusterVersion,
	}
}

// ParseEnvironmentManagedClusterVersionID parses 'input' into a EnvironmentManagedClusterVersionId
func ParseEnvironmentManagedClusterVersionID(input string) (*EnvironmentManagedClusterVersionId, error) {
	parser := resourceids.NewParserFromResourceIdType(EnvironmentManagedClusterVersionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := EnvironmentManagedClusterVersionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.Location, ok = parsed.Parsed["location"]; !ok {
		return nil, fmt.Errorf("the segment 'location' was not found in the resource id %q", input)
	}

	if v, constFound := parsed.Parsed["environment"]; true {
		if !constFound {
			return nil, fmt.Errorf("the segment 'environment' was not found in the resource id %q", input)
		}

		environment, err := parseEnvironment(v)
		if err != nil {
			return nil, fmt.Errorf("parsing %q: %+v", v, err)
		}
		id.Environment = *environment
	}

	if id.ClusterVersion, ok = parsed.Parsed["clusterVersion"]; !ok {
		return nil, fmt.Errorf("the segment 'clusterVersion' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseEnvironmentManagedClusterVersionIDInsensitively parses 'input' case-insensitively into a EnvironmentManagedClusterVersionId
// note: this method should only be used for API response data and not user input
func ParseEnvironmentManagedClusterVersionIDInsensitively(input string) (*EnvironmentManagedClusterVersionId, error) {
	parser := resourceids.NewParserFromResourceIdType(EnvironmentManagedClusterVersionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := EnvironmentManagedClusterVersionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.Location, ok = parsed.Parsed["location"]; !ok {
		return nil, fmt.Errorf("the segment 'location' was not found in the resource id %q", input)
	}

	if v, constFound := parsed.Parsed["environment"]; true {
		if !constFound {
			return nil, fmt.Errorf("the segment 'environment' was not found in the resource id %q", input)
		}

		environment, err := parseEnvironment(v)
		if err != nil {
			return nil, fmt.Errorf("parsing %q: %+v", v, err)
		}
		id.Environment = *environment
	}

	if id.ClusterVersion, ok = parsed.Parsed["clusterVersion"]; !ok {
		return nil, fmt.Errorf("the segment 'clusterVersion' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateEnvironmentManagedClusterVersionID checks that 'input' can be parsed as a Environment Managed Cluster Version ID
func ValidateEnvironmentManagedClusterVersionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseEnvironmentManagedClusterVersionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Environment Managed Cluster Version ID
func (id EnvironmentManagedClusterVersionId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.ServiceFabric/locations/%s/environments/%s/managedClusterVersions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.Location, string(id.Environment), id.ClusterVersion)
}

// Segments returns a slice of Resource ID Segments which comprise this Environment Managed Cluster Version ID
func (id EnvironmentManagedClusterVersionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("subscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("providers", "providers", "providers"),
		resourceids.ResourceProviderSegment("microsoftServiceFabric", "Microsoft.ServiceFabric", "Microsoft.ServiceFabric"),
		resourceids.StaticSegment("locations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("location", "locationValue"),
		resourceids.StaticSegment("environments", "environments", "environments"),
		resourceids.ConstantSegment("environment", PossibleValuesForEnvironment(), "Windows"),
		resourceids.StaticSegment("managedClusterVersions", "managedClusterVersions", "managedClusterVersions"),
		resourceids.UserSpecifiedSegment("clusterVersion", "clusterVersionValue"),
	}
}

// String returns a human-readable description of this Environment Managed Cluster Version ID
func (id EnvironmentManagedClusterVersionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location: %q", id.Location),
		fmt.Sprintf("Environment: %q", string(id.Environment)),
		fmt.Sprintf("Cluster Version: %q", id.ClusterVersion),
	}
	return fmt.Sprintf("Environment Managed Cluster Version (%s)", strings.Join(components, "\n"))
}
