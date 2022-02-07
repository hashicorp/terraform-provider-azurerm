package heatmaps

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = HeatMapTypeId{}

// HeatMapTypeId is a struct representing the Resource ID for a Heat Map Type
type HeatMapTypeId struct {
	SubscriptionId    string
	ResourceGroupName string
	ProfileName       string
	HeatMapType       HeatMapType
}

// NewHeatMapTypeID returns a new HeatMapTypeId struct
func NewHeatMapTypeID(subscriptionId string, resourceGroupName string, profileName string, heatMapType HeatMapType) HeatMapTypeId {
	return HeatMapTypeId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ProfileName:       profileName,
		HeatMapType:       heatMapType,
	}
}

// ParseHeatMapTypeID parses 'input' into a HeatMapTypeId
func ParseHeatMapTypeID(input string) (*HeatMapTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(HeatMapTypeId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := HeatMapTypeId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ProfileName, ok = parsed.Parsed["profileName"]; !ok {
		return nil, fmt.Errorf("the segment 'profileName' was not found in the resource id %q", input)
	}

	if v, constFound := parsed.Parsed["heatMapType"]; true {
		if !constFound {
			return nil, fmt.Errorf("the segment 'heatMapType' was not found in the resource id %q", input)
		}

		heatMapType, err := parseHeatMapType(v)
		if err != nil {
			return nil, fmt.Errorf("parsing %q: %+v", v, err)
		}
		id.HeatMapType = *heatMapType
	}

	return &id, nil
}

// ParseHeatMapTypeIDInsensitively parses 'input' case-insensitively into a HeatMapTypeId
// note: this method should only be used for API response data and not user input
func ParseHeatMapTypeIDInsensitively(input string) (*HeatMapTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(HeatMapTypeId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := HeatMapTypeId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ProfileName, ok = parsed.Parsed["profileName"]; !ok {
		return nil, fmt.Errorf("the segment 'profileName' was not found in the resource id %q", input)
	}

	if v, constFound := parsed.Parsed["heatMapType"]; true {
		if !constFound {
			return nil, fmt.Errorf("the segment 'heatMapType' was not found in the resource id %q", input)
		}

		heatMapType, err := parseHeatMapType(v)
		if err != nil {
			return nil, fmt.Errorf("parsing %q: %+v", v, err)
		}
		id.HeatMapType = *heatMapType
	}

	return &id, nil
}

// ValidateHeatMapTypeID checks that 'input' can be parsed as a Heat Map Type ID
func ValidateHeatMapTypeID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseHeatMapTypeID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Heat Map Type ID
func (id HeatMapTypeId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/trafficManagerProfiles/%s/heatMaps/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ProfileName, string(id.HeatMapType))
}

// Segments returns a slice of Resource ID Segments which comprise this Heat Map Type ID
func (id HeatMapTypeId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticTrafficManagerProfiles", "trafficManagerProfiles", "trafficManagerProfiles"),
		resourceids.UserSpecifiedSegment("profileName", "profileValue"),
		resourceids.StaticSegment("staticHeatMaps", "heatMaps", "heatMaps"),
		resourceids.ConstantSegment("heatMapType", PossibleValuesForHeatMapType(), "default"),
	}
}

// String returns a human-readable description of this Heat Map Type ID
func (id HeatMapTypeId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Profile Name: %q", id.ProfileName),
		fmt.Sprintf("Heat Map Type: %q", string(id.HeatMapType)),
	}
	return fmt.Sprintf("Heat Map Type (%s)", strings.Join(components, "\n"))
}
