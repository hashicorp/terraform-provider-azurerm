package datalakestoreaccounts

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = DataLakeStoreAccountId{}

// DataLakeStoreAccountId is a struct representing the Resource ID for a Data Lake Store Account
type DataLakeStoreAccountId struct {
	SubscriptionId           string
	ResourceGroupName        string
	AccountName              string
	DataLakeStoreAccountName string
}

// NewDataLakeStoreAccountID returns a new DataLakeStoreAccountId struct
func NewDataLakeStoreAccountID(subscriptionId string, resourceGroupName string, accountName string, dataLakeStoreAccountName string) DataLakeStoreAccountId {
	return DataLakeStoreAccountId{
		SubscriptionId:           subscriptionId,
		ResourceGroupName:        resourceGroupName,
		AccountName:              accountName,
		DataLakeStoreAccountName: dataLakeStoreAccountName,
	}
}

// ParseDataLakeStoreAccountID parses 'input' into a DataLakeStoreAccountId
func ParseDataLakeStoreAccountID(input string) (*DataLakeStoreAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(DataLakeStoreAccountId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DataLakeStoreAccountId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.DataLakeStoreAccountName, ok = parsed.Parsed["dataLakeStoreAccountName"]; !ok {
		return nil, fmt.Errorf("the segment 'dataLakeStoreAccountName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseDataLakeStoreAccountIDInsensitively parses 'input' case-insensitively into a DataLakeStoreAccountId
// note: this method should only be used for API response data and not user input
func ParseDataLakeStoreAccountIDInsensitively(input string) (*DataLakeStoreAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(DataLakeStoreAccountId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DataLakeStoreAccountId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.DataLakeStoreAccountName, ok = parsed.Parsed["dataLakeStoreAccountName"]; !ok {
		return nil, fmt.Errorf("the segment 'dataLakeStoreAccountName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateDataLakeStoreAccountID checks that 'input' can be parsed as a Data Lake Store Account ID
func ValidateDataLakeStoreAccountID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDataLakeStoreAccountID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Data Lake Store Account ID
func (id DataLakeStoreAccountId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataLakeAnalytics/accounts/%s/dataLakeStoreAccounts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.DataLakeStoreAccountName)
}

// Segments returns a slice of Resource ID Segments which comprise this Data Lake Store Account ID
func (id DataLakeStoreAccountId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDataLakeAnalytics", "Microsoft.DataLakeAnalytics", "Microsoft.DataLakeAnalytics"),
		resourceids.StaticSegment("staticAccounts", "accounts", "accounts"),
		resourceids.UserSpecifiedSegment("accountName", "accountValue"),
		resourceids.StaticSegment("staticDataLakeStoreAccounts", "dataLakeStoreAccounts", "dataLakeStoreAccounts"),
		resourceids.UserSpecifiedSegment("dataLakeStoreAccountName", "dataLakeStoreAccountValue"),
	}
}

// String returns a human-readable description of this Data Lake Store Account ID
func (id DataLakeStoreAccountId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Account Name: %q", id.AccountName),
		fmt.Sprintf("Data Lake Store Account Name: %q", id.DataLakeStoreAccountName),
	}
	return fmt.Sprintf("Data Lake Store Account (%s)", strings.Join(components, "\n"))
}
