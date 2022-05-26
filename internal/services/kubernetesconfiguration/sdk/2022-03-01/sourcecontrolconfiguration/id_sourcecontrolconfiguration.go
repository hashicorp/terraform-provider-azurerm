package sourcecontrolconfiguration

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = SourceControlConfigurationId{}

// SourceControlConfigurationId is a struct representing the Resource ID for a Source Control Configuration
type SourceControlConfigurationId struct {
	SubscriptionId                 string
	ResourceGroupName              string
	ClusterRp                      string
	ClusterResourceName            string
	ClusterName                    string
	SourceControlConfigurationName string
}

// NewSourceControlConfigurationID returns a new SourceControlConfigurationId struct
func NewSourceControlConfigurationID(subscriptionId string, resourceGroupName string, clusterRp string, clusterResourceName string, clusterName string, sourceControlConfigurationName string) SourceControlConfigurationId {
	return SourceControlConfigurationId{
		SubscriptionId:                 subscriptionId,
		ResourceGroupName:              resourceGroupName,
		ClusterRp:                      clusterRp,
		ClusterResourceName:            clusterResourceName,
		ClusterName:                    clusterName,
		SourceControlConfigurationName: sourceControlConfigurationName,
	}
}

// ParseSourceControlConfigurationID parses 'input' into a SourceControlConfigurationId
func ParseSourceControlConfigurationID(input string) (*SourceControlConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(SourceControlConfigurationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SourceControlConfigurationId{}

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

	if id.SourceControlConfigurationName, ok = parsed.Parsed["sourceControlConfigurationName"]; !ok {
		return nil, fmt.Errorf("the segment 'sourceControlConfigurationName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseSourceControlConfigurationIDInsensitively parses 'input' case-insensitively into a SourceControlConfigurationId
// note: this method should only be used for API response data and not user input
func ParseSourceControlConfigurationIDInsensitively(input string) (*SourceControlConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(SourceControlConfigurationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SourceControlConfigurationId{}

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

	if id.SourceControlConfigurationName, ok = parsed.Parsed["sourceControlConfigurationName"]; !ok {
		return nil, fmt.Errorf("the segment 'sourceControlConfigurationName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateSourceControlConfigurationID checks that 'input' can be parsed as a Source Control Configuration ID
func ValidateSourceControlConfigurationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSourceControlConfigurationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Source Control Configuration ID
func (id SourceControlConfigurationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/%s/%s/%s/providers/Microsoft.KubernetesConfiguration/sourceControlConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ClusterRp, id.ClusterResourceName, id.ClusterName, id.SourceControlConfigurationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Source Control Configuration ID
func (id SourceControlConfigurationId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticSourceControlConfigurations", "sourceControlConfigurations", "sourceControlConfigurations"),
		resourceids.UserSpecifiedSegment("sourceControlConfigurationName", "sourceControlConfigurationValue"),
	}
}

// String returns a human-readable description of this Source Control Configuration ID
func (id SourceControlConfigurationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Cluster Rp: %q", id.ClusterRp),
		fmt.Sprintf("Cluster Resource Name: %q", id.ClusterResourceName),
		fmt.Sprintf("Cluster Name: %q", id.ClusterName),
		fmt.Sprintf("Source Control Configuration Name: %q", id.SourceControlConfigurationName),
	}
	return fmt.Sprintf("Source Control Configuration (%s)", strings.Join(components, "\n"))
}
