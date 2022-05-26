package extensionoperationstatus

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = OperationId{}

// OperationId is a struct representing the Resource ID for a Operation
type OperationId struct {
	SubscriptionId      string
	ResourceGroupName   string
	ClusterRp           string
	ClusterResourceName string
	ClusterName         string
	ExtensionName       string
	OperationId         string
}

// NewOperationID returns a new OperationId struct
func NewOperationID(subscriptionId string, resourceGroupName string, clusterRp string, clusterResourceName string, clusterName string, extensionName string, operationId string) OperationId {
	return OperationId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		ClusterRp:           clusterRp,
		ClusterResourceName: clusterResourceName,
		ClusterName:         clusterName,
		ExtensionName:       extensionName,
		OperationId:         operationId,
	}
}

// ParseOperationID parses 'input' into a OperationId
func ParseOperationID(input string) (*OperationId, error) {
	parser := resourceids.NewParserFromResourceIdType(OperationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := OperationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ClusterRp, ok = parsed.Parsed["clusterRp"]; !ok {
		return nil, fmt.Errorf("the segment 'clusterRp' was not found in the resource id %q", input)
	}

	if id.ClusterResourceName, ok = parsed.Parsed["clusterResourceName"]; !ok {
		return nil, fmt.Errorf("the segment 'clusterResourceName' was not found in the resource id %q", input)
	}

	if id.ClusterName, ok = parsed.Parsed["clusterName"]; !ok {
		return nil, fmt.Errorf("the segment 'clusterName' was not found in the resource id %q", input)
	}

	if id.ExtensionName, ok = parsed.Parsed["extensionName"]; !ok {
		return nil, fmt.Errorf("the segment 'extensionName' was not found in the resource id %q", input)
	}

	if id.OperationId, ok = parsed.Parsed["operationId"]; !ok {
		return nil, fmt.Errorf("the segment 'operationId' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseOperationIDInsensitively parses 'input' case-insensitively into a OperationId
// note: this method should only be used for API response data and not user input
func ParseOperationIDInsensitively(input string) (*OperationId, error) {
	parser := resourceids.NewParserFromResourceIdType(OperationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := OperationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ClusterRp, ok = parsed.Parsed["clusterRp"]; !ok {
		return nil, fmt.Errorf("the segment 'clusterRp' was not found in the resource id %q", input)
	}

	if id.ClusterResourceName, ok = parsed.Parsed["clusterResourceName"]; !ok {
		return nil, fmt.Errorf("the segment 'clusterResourceName' was not found in the resource id %q", input)
	}

	if id.ClusterName, ok = parsed.Parsed["clusterName"]; !ok {
		return nil, fmt.Errorf("the segment 'clusterName' was not found in the resource id %q", input)
	}

	if id.ExtensionName, ok = parsed.Parsed["extensionName"]; !ok {
		return nil, fmt.Errorf("the segment 'extensionName' was not found in the resource id %q", input)
	}

	if id.OperationId, ok = parsed.Parsed["operationId"]; !ok {
		return nil, fmt.Errorf("the segment 'operationId' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateOperationID checks that 'input' can be parsed as a Operation ID
func ValidateOperationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseOperationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Operation ID
func (id OperationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/%s/%s/%s/providers/Microsoft.KubernetesConfiguration/extensions/%s/operations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ClusterRp, id.ClusterResourceName, id.ClusterName, id.ExtensionName, id.OperationId)
}

// Segments returns a slice of Resource ID Segments which comprise this Operation ID
func (id OperationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.UserSpecifiedSegment("clusterRp", "clusterRpValue"),
		resourceids.UserSpecifiedSegment("clusterResourceName", "clusterResourceValue"),
		resourceids.UserSpecifiedSegment("clusterName", "clusterValue"),
		resourceids.StaticSegment("staticProviders2", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftKubernetesConfiguration", "Microsoft.KubernetesConfiguration", "Microsoft.KubernetesConfiguration"),
		resourceids.StaticSegment("staticExtensions", "extensions", "extensions"),
		resourceids.UserSpecifiedSegment("extensionName", "extensionValue"),
		resourceids.StaticSegment("staticOperations", "operations", "operations"),
		resourceids.UserSpecifiedSegment("operationId", "operationIdValue"),
	}
}

// String returns a human-readable description of this Operation ID
func (id OperationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Cluster Rp: %q", id.ClusterRp),
		fmt.Sprintf("Cluster Resource Name: %q", id.ClusterResourceName),
		fmt.Sprintf("Cluster Name: %q", id.ClusterName),
		fmt.Sprintf("Extension Name: %q", id.ExtensionName),
		fmt.Sprintf("Operation: %q", id.OperationId),
	}
	return fmt.Sprintf("Operation (%s)", strings.Join(components, "\n"))
}
