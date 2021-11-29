package applicationtypeversion

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = VersionId{}

// VersionId is a struct representing the Resource ID for a Version
type VersionId struct {
	SubscriptionId      string
	ResourceGroupName   string
	ClusterName         string
	ApplicationTypeName string
	Version             string
}

// NewVersionID returns a new VersionId struct
func NewVersionID(subscriptionId string, resourceGroupName string, clusterName string, applicationTypeName string, version string) VersionId {
	return VersionId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		ClusterName:         clusterName,
		ApplicationTypeName: applicationTypeName,
		Version:             version,
	}
}

// ParseVersionID parses 'input' into a VersionId
func ParseVersionID(input string) (*VersionId, error) {
	parser := resourceids.NewParserFromResourceIdType(VersionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VersionId{}

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

	if id.Version, ok = parsed.Parsed["version"]; !ok {
		return nil, fmt.Errorf("the segment 'version' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseVersionIDInsensitively parses 'input' case-insensitively into a VersionId
// note: this method should only be used for API response data and not user input
func ParseVersionIDInsensitively(input string) (*VersionId, error) {
	parser := resourceids.NewParserFromResourceIdType(VersionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VersionId{}

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

	if id.Version, ok = parsed.Parsed["version"]; !ok {
		return nil, fmt.Errorf("the segment 'version' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateVersionID checks that 'input' can be parsed as a Version ID
func ValidateVersionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVersionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Version ID
func (id VersionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ServiceFabric/managedClusters/%s/applicationTypes/%s/versions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ClusterName, id.ApplicationTypeName, id.Version)
}

// Segments returns a slice of Resource ID Segments which comprise this Version ID
func (id VersionId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("versions", "versions", "versions"),
		resourceids.UserSpecifiedSegment("version", "versionValue"),
	}
}

// String returns a human-readable description of this Version ID
func (id VersionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Cluster Name: %q", id.ClusterName),
		fmt.Sprintf("Application Type Name: %q", id.ApplicationTypeName),
		fmt.Sprintf("Version: %q", id.Version),
	}
	return fmt.Sprintf("Version (%s)", strings.Join(components, "\n"))
}
