package assetsandassetfilters

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = AssetId{}

// AssetId is a struct representing the Resource ID for a Asset
type AssetId struct {
	SubscriptionId    string
	ResourceGroupName string
	AccountName       string
	AssetName         string
}

// NewAssetID returns a new AssetId struct
func NewAssetID(subscriptionId string, resourceGroupName string, accountName string, assetName string) AssetId {
	return AssetId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		AccountName:       accountName,
		AssetName:         assetName,
	}
}

// ParseAssetID parses 'input' into a AssetId
func ParseAssetID(input string) (*AssetId, error) {
	parser := resourceids.NewParserFromResourceIdType(AssetId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AssetId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.AssetName, ok = parsed.Parsed["assetName"]; !ok {
		return nil, fmt.Errorf("the segment 'assetName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseAssetIDInsensitively parses 'input' case-insensitively into a AssetId
// note: this method should only be used for API response data and not user input
func ParseAssetIDInsensitively(input string) (*AssetId, error) {
	parser := resourceids.NewParserFromResourceIdType(AssetId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AssetId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.AssetName, ok = parsed.Parsed["assetName"]; !ok {
		return nil, fmt.Errorf("the segment 'assetName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateAssetID checks that 'input' can be parsed as a Asset ID
func ValidateAssetID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAssetID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Asset ID
func (id AssetId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Media/mediaServices/%s/assets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.AssetName)
}

// Segments returns a slice of Resource ID Segments which comprise this Asset ID
func (id AssetId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMedia", "Microsoft.Media", "Microsoft.Media"),
		resourceids.StaticSegment("staticMediaServices", "mediaServices", "mediaServices"),
		resourceids.UserSpecifiedSegment("accountName", "accountValue"),
		resourceids.StaticSegment("staticAssets", "assets", "assets"),
		resourceids.UserSpecifiedSegment("assetName", "assetValue"),
	}
}

// String returns a human-readable description of this Asset ID
func (id AssetId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Account Name: %q", id.AccountName),
		fmt.Sprintf("Asset Name: %q", id.AssetName),
	}
	return fmt.Sprintf("Asset (%s)", strings.Join(components, "\n"))
}
