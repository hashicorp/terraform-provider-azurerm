package trustedidproviders

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = TrustedIdProviderId{}

// TrustedIdProviderId is a struct representing the Resource ID for a Trusted Id Provider
type TrustedIdProviderId struct {
	SubscriptionId        string
	ResourceGroupName     string
	AccountName           string
	TrustedIdProviderName string
}

// NewTrustedIdProviderID returns a new TrustedIdProviderId struct
func NewTrustedIdProviderID(subscriptionId string, resourceGroupName string, accountName string, trustedIdProviderName string) TrustedIdProviderId {
	return TrustedIdProviderId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		AccountName:           accountName,
		TrustedIdProviderName: trustedIdProviderName,
	}
}

// ParseTrustedIdProviderID parses 'input' into a TrustedIdProviderId
func ParseTrustedIdProviderID(input string) (*TrustedIdProviderId, error) {
	parser := resourceids.NewParserFromResourceIdType(TrustedIdProviderId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := TrustedIdProviderId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.TrustedIdProviderName, ok = parsed.Parsed["trustedIdProviderName"]; !ok {
		return nil, fmt.Errorf("the segment 'trustedIdProviderName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseTrustedIdProviderIDInsensitively parses 'input' case-insensitively into a TrustedIdProviderId
// note: this method should only be used for API response data and not user input
func ParseTrustedIdProviderIDInsensitively(input string) (*TrustedIdProviderId, error) {
	parser := resourceids.NewParserFromResourceIdType(TrustedIdProviderId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := TrustedIdProviderId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.TrustedIdProviderName, ok = parsed.Parsed["trustedIdProviderName"]; !ok {
		return nil, fmt.Errorf("the segment 'trustedIdProviderName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateTrustedIdProviderID checks that 'input' can be parsed as a Trusted Id Provider ID
func ValidateTrustedIdProviderID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseTrustedIdProviderID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Trusted Id Provider ID
func (id TrustedIdProviderId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataLakeStore/accounts/%s/trustedIdProviders/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.TrustedIdProviderName)
}

// Segments returns a slice of Resource ID Segments which comprise this Trusted Id Provider ID
func (id TrustedIdProviderId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDataLakeStore", "Microsoft.DataLakeStore", "Microsoft.DataLakeStore"),
		resourceids.StaticSegment("staticAccounts", "accounts", "accounts"),
		resourceids.UserSpecifiedSegment("accountName", "accountValue"),
		resourceids.StaticSegment("staticTrustedIdProviders", "trustedIdProviders", "trustedIdProviders"),
		resourceids.UserSpecifiedSegment("trustedIdProviderName", "trustedIdProviderValue"),
	}
}

// String returns a human-readable description of this Trusted Id Provider ID
func (id TrustedIdProviderId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Account Name: %q", id.AccountName),
		fmt.Sprintf("Trusted Id Provider Name: %q", id.TrustedIdProviderName),
	}
	return fmt.Sprintf("Trusted Id Provider (%s)", strings.Join(components, "\n"))
}
