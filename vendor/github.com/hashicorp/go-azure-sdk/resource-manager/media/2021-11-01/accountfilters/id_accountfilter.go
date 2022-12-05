package accountfilters

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = AccountFilterId{}

// AccountFilterId is a struct representing the Resource ID for a Account Filter
type AccountFilterId struct {
	SubscriptionId    string
	ResourceGroupName string
	AccountName       string
	FilterName        string
}

// NewAccountFilterID returns a new AccountFilterId struct
func NewAccountFilterID(subscriptionId string, resourceGroupName string, accountName string, filterName string) AccountFilterId {
	return AccountFilterId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		AccountName:       accountName,
		FilterName:        filterName,
	}
}

// ParseAccountFilterID parses 'input' into a AccountFilterId
func ParseAccountFilterID(input string) (*AccountFilterId, error) {
	parser := resourceids.NewParserFromResourceIdType(AccountFilterId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AccountFilterId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.FilterName, ok = parsed.Parsed["filterName"]; !ok {
		return nil, fmt.Errorf("the segment 'filterName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseAccountFilterIDInsensitively parses 'input' case-insensitively into a AccountFilterId
// note: this method should only be used for API response data and not user input
func ParseAccountFilterIDInsensitively(input string) (*AccountFilterId, error) {
	parser := resourceids.NewParserFromResourceIdType(AccountFilterId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AccountFilterId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.FilterName, ok = parsed.Parsed["filterName"]; !ok {
		return nil, fmt.Errorf("the segment 'filterName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateAccountFilterID checks that 'input' can be parsed as a Account Filter ID
func ValidateAccountFilterID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAccountFilterID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Account Filter ID
func (id AccountFilterId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Media/mediaServices/%s/accountFilters/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.FilterName)
}

// Segments returns a slice of Resource ID Segments which comprise this Account Filter ID
func (id AccountFilterId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMedia", "Microsoft.Media", "Microsoft.Media"),
		resourceids.StaticSegment("staticMediaServices", "mediaServices", "mediaServices"),
		resourceids.UserSpecifiedSegment("accountName", "accountValue"),
		resourceids.StaticSegment("staticAccountFilters", "accountFilters", "accountFilters"),
		resourceids.UserSpecifiedSegment("filterName", "filterValue"),
	}
}

// String returns a human-readable description of this Account Filter ID
func (id AccountFilterId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Account Name: %q", id.AccountName),
		fmt.Sprintf("Filter Name: %q", id.FilterName),
	}
	return fmt.Sprintf("Account Filter (%s)", strings.Join(components, "\n"))
}
