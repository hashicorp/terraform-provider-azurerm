package resource

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = SpatialAnchorsAccountId{}

// SpatialAnchorsAccountId is a struct representing the Resource ID for a Spatial Anchors Account
type SpatialAnchorsAccountId struct {
	SubscriptionId    string
	ResourceGroupName string
	AccountName       string
}

// NewSpatialAnchorsAccountID returns a new SpatialAnchorsAccountId struct
func NewSpatialAnchorsAccountID(subscriptionId string, resourceGroupName string, accountName string) SpatialAnchorsAccountId {
	return SpatialAnchorsAccountId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		AccountName:       accountName,
	}
}

// ParseSpatialAnchorsAccountID parses 'input' into a SpatialAnchorsAccountId
func ParseSpatialAnchorsAccountID(input string) (*SpatialAnchorsAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(SpatialAnchorsAccountId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SpatialAnchorsAccountId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseSpatialAnchorsAccountIDInsensitively parses 'input' case-insensitively into a SpatialAnchorsAccountId
// note: this method should only be used for API response data and not user input
func ParseSpatialAnchorsAccountIDInsensitively(input string) (*SpatialAnchorsAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(SpatialAnchorsAccountId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SpatialAnchorsAccountId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateSpatialAnchorsAccountID checks that 'input' can be parsed as a Spatial Anchors Account ID
func ValidateSpatialAnchorsAccountID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSpatialAnchorsAccountID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Spatial Anchors Account ID
func (id SpatialAnchorsAccountId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.MixedReality/spatialAnchorsAccounts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AccountName)
}

// Segments returns a slice of Resource ID Segments which comprise this Spatial Anchors Account ID
func (id SpatialAnchorsAccountId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMixedReality", "Microsoft.MixedReality", "Microsoft.MixedReality"),
		resourceids.StaticSegment("staticSpatialAnchorsAccounts", "spatialAnchorsAccounts", "spatialAnchorsAccounts"),
		resourceids.UserSpecifiedSegment("accountName", "accountValue"),
	}
}

// String returns a human-readable description of this Spatial Anchors Account ID
func (id SpatialAnchorsAccountId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Account Name: %q", id.AccountName),
	}
	return fmt.Sprintf("Spatial Anchors Account (%s)", strings.Join(components, "\n"))
}
