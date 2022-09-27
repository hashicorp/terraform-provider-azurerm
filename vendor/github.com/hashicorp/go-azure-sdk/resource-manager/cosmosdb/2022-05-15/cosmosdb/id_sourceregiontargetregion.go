package cosmosdb

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = SourceRegionTargetRegionId{}

// SourceRegionTargetRegionId is a struct representing the Resource ID for a Source Region Target Region
type SourceRegionTargetRegionId struct {
	SubscriptionId    string
	ResourceGroupName string
	AccountName       string
	SourceRegion      string
	TargetRegion      string
}

// NewSourceRegionTargetRegionID returns a new SourceRegionTargetRegionId struct
func NewSourceRegionTargetRegionID(subscriptionId string, resourceGroupName string, accountName string, sourceRegion string, targetRegion string) SourceRegionTargetRegionId {
	return SourceRegionTargetRegionId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		AccountName:       accountName,
		SourceRegion:      sourceRegion,
		TargetRegion:      targetRegion,
	}
}

// ParseSourceRegionTargetRegionID parses 'input' into a SourceRegionTargetRegionId
func ParseSourceRegionTargetRegionID(input string) (*SourceRegionTargetRegionId, error) {
	parser := resourceids.NewParserFromResourceIdType(SourceRegionTargetRegionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SourceRegionTargetRegionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.SourceRegion, ok = parsed.Parsed["sourceRegion"]; !ok {
		return nil, fmt.Errorf("the segment 'sourceRegion' was not found in the resource id %q", input)
	}

	if id.TargetRegion, ok = parsed.Parsed["targetRegion"]; !ok {
		return nil, fmt.Errorf("the segment 'targetRegion' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseSourceRegionTargetRegionIDInsensitively parses 'input' case-insensitively into a SourceRegionTargetRegionId
// note: this method should only be used for API response data and not user input
func ParseSourceRegionTargetRegionIDInsensitively(input string) (*SourceRegionTargetRegionId, error) {
	parser := resourceids.NewParserFromResourceIdType(SourceRegionTargetRegionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SourceRegionTargetRegionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.SourceRegion, ok = parsed.Parsed["sourceRegion"]; !ok {
		return nil, fmt.Errorf("the segment 'sourceRegion' was not found in the resource id %q", input)
	}

	if id.TargetRegion, ok = parsed.Parsed["targetRegion"]; !ok {
		return nil, fmt.Errorf("the segment 'targetRegion' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateSourceRegionTargetRegionID checks that 'input' can be parsed as a Source Region Target Region ID
func ValidateSourceRegionTargetRegionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSourceRegionTargetRegionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Source Region Target Region ID
func (id SourceRegionTargetRegionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/databaseAccounts/%s/sourceRegion/%s/targetRegion/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.SourceRegion, id.TargetRegion)
}

// Segments returns a slice of Resource ID Segments which comprise this Source Region Target Region ID
func (id SourceRegionTargetRegionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDocumentDB", "Microsoft.DocumentDB", "Microsoft.DocumentDB"),
		resourceids.StaticSegment("staticDatabaseAccounts", "databaseAccounts", "databaseAccounts"),
		resourceids.UserSpecifiedSegment("accountName", "accountValue"),
		resourceids.StaticSegment("staticSourceRegion", "sourceRegion", "sourceRegion"),
		resourceids.UserSpecifiedSegment("sourceRegion", "sourceRegionValue"),
		resourceids.StaticSegment("staticTargetRegion", "targetRegion", "targetRegion"),
		resourceids.UserSpecifiedSegment("targetRegion", "targetRegionValue"),
	}
}

// String returns a human-readable description of this Source Region Target Region ID
func (id SourceRegionTargetRegionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Account Name: %q", id.AccountName),
		fmt.Sprintf("Source Region: %q", id.SourceRegion),
		fmt.Sprintf("Target Region: %q", id.TargetRegion),
	}
	return fmt.Sprintf("Source Region Target Region (%s)", strings.Join(components, "\n"))
}
