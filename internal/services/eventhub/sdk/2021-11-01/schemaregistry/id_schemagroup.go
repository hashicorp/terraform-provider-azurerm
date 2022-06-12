package schemaregistry

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = SchemagroupId{}

// SchemagroupId is a struct representing the Resource ID for a Schemagroup
type SchemagroupId struct {
	SubscriptionId    string
	ResourceGroupName string
	NamespaceName     string
	SchemaGroupName   string
}

// NewSchemagroupID returns a new SchemagroupId struct
func NewSchemagroupID(subscriptionId string, resourceGroupName string, namespaceName string, schemaGroupName string) SchemagroupId {
	return SchemagroupId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		NamespaceName:     namespaceName,
		SchemaGroupName:   schemaGroupName,
	}
}

// ParseSchemagroupID parses 'input' into a SchemagroupId
func ParseSchemagroupID(input string) (*SchemagroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(SchemagroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SchemagroupId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.NamespaceName, ok = parsed.Parsed["namespaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'namespaceName' was not found in the resource id %q", input)
	}

	if id.SchemaGroupName, ok = parsed.Parsed["schemaGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'schemaGroupName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseSchemagroupIDInsensitively parses 'input' case-insensitively into a SchemagroupId
// note: this method should only be used for API response data and not user input
func ParseSchemagroupIDInsensitively(input string) (*SchemagroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(SchemagroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SchemagroupId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.NamespaceName, ok = parsed.Parsed["namespaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'namespaceName' was not found in the resource id %q", input)
	}

	if id.SchemaGroupName, ok = parsed.Parsed["schemaGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'schemaGroupName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateSchemagroupID checks that 'input' can be parsed as a Schemagroup ID
func ValidateSchemagroupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSchemagroupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Schemagroup ID
func (id SchemagroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventHub/namespaces/%s/schemagroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.SchemaGroupName)
}

// Segments returns a slice of Resource ID Segments which comprise this Schemagroup ID
func (id SchemagroupId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftEventHub", "Microsoft.EventHub", "Microsoft.EventHub"),
		resourceids.StaticSegment("staticNamespaces", "namespaces", "namespaces"),
		resourceids.UserSpecifiedSegment("namespaceName", "namespaceValue"),
		resourceids.StaticSegment("staticSchemagroups", "schemagroups", "schemagroups"),
		resourceids.UserSpecifiedSegment("schemaGroupName", "schemaGroupValue"),
	}
}

// String returns a human-readable description of this Schemagroup ID
func (id SchemagroupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Namespace Name: %q", id.NamespaceName),
		fmt.Sprintf("Schema Group Name: %q", id.SchemaGroupName),
	}
	return fmt.Sprintf("Schemagroup (%s)", strings.Join(components, "\n"))
}
