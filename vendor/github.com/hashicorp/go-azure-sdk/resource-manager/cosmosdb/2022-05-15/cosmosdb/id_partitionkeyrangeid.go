package cosmosdb

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = PartitionKeyRangeIdId{}

// PartitionKeyRangeIdId is a struct representing the Resource ID for a Partition Key Range Id
type PartitionKeyRangeIdId struct {
	SubscriptionId      string
	ResourceGroupName   string
	AccountName         string
	DatabaseRid         string
	CollectionRid       string
	PartitionKeyRangeId string
}

// NewPartitionKeyRangeIdID returns a new PartitionKeyRangeIdId struct
func NewPartitionKeyRangeIdID(subscriptionId string, resourceGroupName string, accountName string, databaseRid string, collectionRid string, partitionKeyRangeId string) PartitionKeyRangeIdId {
	return PartitionKeyRangeIdId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		AccountName:         accountName,
		DatabaseRid:         databaseRid,
		CollectionRid:       collectionRid,
		PartitionKeyRangeId: partitionKeyRangeId,
	}
}

// ParsePartitionKeyRangeIdID parses 'input' into a PartitionKeyRangeIdId
func ParsePartitionKeyRangeIdID(input string) (*PartitionKeyRangeIdId, error) {
	parser := resourceids.NewParserFromResourceIdType(PartitionKeyRangeIdId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PartitionKeyRangeIdId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.DatabaseRid, ok = parsed.Parsed["databaseRid"]; !ok {
		return nil, fmt.Errorf("the segment 'databaseRid' was not found in the resource id %q", input)
	}

	if id.CollectionRid, ok = parsed.Parsed["collectionRid"]; !ok {
		return nil, fmt.Errorf("the segment 'collectionRid' was not found in the resource id %q", input)
	}

	if id.PartitionKeyRangeId, ok = parsed.Parsed["partitionKeyRangeId"]; !ok {
		return nil, fmt.Errorf("the segment 'partitionKeyRangeId' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParsePartitionKeyRangeIdIDInsensitively parses 'input' case-insensitively into a PartitionKeyRangeIdId
// note: this method should only be used for API response data and not user input
func ParsePartitionKeyRangeIdIDInsensitively(input string) (*PartitionKeyRangeIdId, error) {
	parser := resourceids.NewParserFromResourceIdType(PartitionKeyRangeIdId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PartitionKeyRangeIdId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.DatabaseRid, ok = parsed.Parsed["databaseRid"]; !ok {
		return nil, fmt.Errorf("the segment 'databaseRid' was not found in the resource id %q", input)
	}

	if id.CollectionRid, ok = parsed.Parsed["collectionRid"]; !ok {
		return nil, fmt.Errorf("the segment 'collectionRid' was not found in the resource id %q", input)
	}

	if id.PartitionKeyRangeId, ok = parsed.Parsed["partitionKeyRangeId"]; !ok {
		return nil, fmt.Errorf("the segment 'partitionKeyRangeId' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidatePartitionKeyRangeIdID checks that 'input' can be parsed as a Partition Key Range Id ID
func ValidatePartitionKeyRangeIdID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePartitionKeyRangeIdID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Partition Key Range Id ID
func (id PartitionKeyRangeIdId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/databaseAccounts/%s/databases/%s/collections/%s/partitionKeyRangeId/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.DatabaseRid, id.CollectionRid, id.PartitionKeyRangeId)
}

// Segments returns a slice of Resource ID Segments which comprise this Partition Key Range Id ID
func (id PartitionKeyRangeIdId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDocumentDB", "Microsoft.DocumentDB", "Microsoft.DocumentDB"),
		resourceids.StaticSegment("staticDatabaseAccounts", "databaseAccounts", "databaseAccounts"),
		resourceids.UserSpecifiedSegment("accountName", "accountValue"),
		resourceids.StaticSegment("staticDatabases", "databases", "databases"),
		resourceids.UserSpecifiedSegment("databaseRid", "databaseRidValue"),
		resourceids.StaticSegment("staticCollections", "collections", "collections"),
		resourceids.UserSpecifiedSegment("collectionRid", "collectionRidValue"),
		resourceids.StaticSegment("staticPartitionKeyRangeId", "partitionKeyRangeId", "partitionKeyRangeId"),
		resourceids.UserSpecifiedSegment("partitionKeyRangeId", "partitionKeyRangeIdValue"),
	}
}

// String returns a human-readable description of this Partition Key Range Id ID
func (id PartitionKeyRangeIdId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Account Name: %q", id.AccountName),
		fmt.Sprintf("Database Rid: %q", id.DatabaseRid),
		fmt.Sprintf("Collection Rid: %q", id.CollectionRid),
		fmt.Sprintf("Partition Key Range: %q", id.PartitionKeyRangeId),
	}
	return fmt.Sprintf("Partition Key Range Id (%s)", strings.Join(components, "\n"))
}
