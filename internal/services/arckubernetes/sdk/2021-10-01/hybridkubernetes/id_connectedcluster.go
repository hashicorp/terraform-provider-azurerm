package hybridkubernetes

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ConnectedClusterId{}

// ConnectedClusterId is a struct representing the Resource ID for a Connected Cluster
type ConnectedClusterId struct {
	SubscriptionId    string
	ResourceGroupName string
	ClusterName       string
}

// NewConnectedClusterID returns a new ConnectedClusterId struct
func NewConnectedClusterID(subscriptionId string, resourceGroupName string, clusterName string) ConnectedClusterId {
	return ConnectedClusterId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ClusterName:       clusterName,
	}
}

// ParseConnectedClusterID parses 'input' into a ConnectedClusterId
func ParseConnectedClusterID(input string) (*ConnectedClusterId, error) {
	parser := resourceids.NewParserFromResourceIdType(ConnectedClusterId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ConnectedClusterId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ClusterName, ok = parsed.Parsed["clusterName"]; !ok {
		return nil, fmt.Errorf("the segment 'clusterName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseConnectedClusterIDInsensitively parses 'input' case-insensitively into a ConnectedClusterId
// note: this method should only be used for API response data and not user input
func ParseConnectedClusterIDInsensitively(input string) (*ConnectedClusterId, error) {
	parser := resourceids.NewParserFromResourceIdType(ConnectedClusterId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ConnectedClusterId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ClusterName, ok = parsed.Parsed["clusterName"]; !ok {
		return nil, fmt.Errorf("the segment 'clusterName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateConnectedClusterID checks that 'input' can be parsed as a Connected Cluster ID
func ValidateConnectedClusterID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseConnectedClusterID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Connected Cluster ID
func (id ConnectedClusterId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Kubernetes/connectedClusters/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ClusterName)
}

// Segments returns a slice of Resource ID Segments which comprise this Connected Cluster ID
func (id ConnectedClusterId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftKubernetes", "Microsoft.Kubernetes", "Microsoft.Kubernetes"),
		resourceids.StaticSegment("staticConnectedClusters", "connectedClusters", "connectedClusters"),
		resourceids.UserSpecifiedSegment("clusterName", "clusterValue"),
	}
}

// String returns a human-readable description of this Connected Cluster ID
func (id ConnectedClusterId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Cluster Name: %q", id.ClusterName),
	}
	return fmt.Sprintf("Connected Cluster (%s)", strings.Join(components, "\n"))
}
