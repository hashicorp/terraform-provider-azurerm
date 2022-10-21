package cosmosdb

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = CollectionPartitionKeyRangeIdId{}

// CollectionPartitionKeyRangeIdId is a struct representing the Resource ID for a Collection Partition Key Range Id
type CollectionPartitionKeyRangeIdId struct {
	SubscriptionId      string
	ResourceGroupName   string
	AccountName         string
	Region              string
	DatabaseRid         string
	CollectionRid       string
	PartitionKeyRangeId string
}

// NewCollectionPartitionKeyRangeIdID returns a new CollectionPartitionKeyRangeIdId struct
func NewCollectionPartitionKeyRangeIdID(subscriptionId string, resourceGroupName string, accountName string, region string, databaseRid string, collectionRid string, partitionKeyRangeId string) CollectionPartitionKeyRangeIdId {
	return CollectionPartitionKeyRangeIdId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		AccountName:         accountName,
		Region:              region,
		DatabaseRid:         databaseRid,
		CollectionRid:       collectionRid,
		PartitionKeyRangeId: partitionKeyRangeId,
	}
}

// ParseCollectionPartitionKeyRangeIdID parses 'input' into a CollectionPartitionKeyRangeIdId
func ParseCollectionPartitionKeyRangeIdID(input string) (*CollectionPartitionKeyRangeIdId, error) {
	parser := resourceids.NewParserFromResourceIdType(CollectionPartitionKeyRangeIdId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CollectionPartitionKeyRangeIdId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.Region, ok = parsed.Parsed["region"]; !ok {
		return nil, fmt.Errorf("the segment 'region' was not found in the resource id %q", input)
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

// ParseCollectionPartitionKeyRangeIdIDInsensitively parses 'input' case-insensitively into a CollectionPartitionKeyRangeIdId
// note: this method should only be used for API response data and not user input
func ParseCollectionPartitionKeyRangeIdIDInsensitively(input string) (*CollectionPartitionKeyRangeIdId, error) {
	parser := resourceids.NewParserFromResourceIdType(CollectionPartitionKeyRangeIdId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CollectionPartitionKeyRangeIdId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.Region, ok = parsed.Parsed["region"]; !ok {
		return nil, fmt.Errorf("the segment 'region' was not found in the resource id %q", input)
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

// ValidateCollectionPartitionKeyRangeIdID checks that 'input' can be parsed as a Collection Partition Key Range Id ID
func ValidateCollectionPartitionKeyRangeIdID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCollectionPartitionKeyRangeIdID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Collection Partition Key Range Id ID
func (id CollectionPartitionKeyRangeIdId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/databaseAccounts/%s/region/%s/databases/%s/collections/%s/partitionKeyRangeId/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.Region, id.DatabaseRid, id.CollectionRid, id.PartitionKeyRangeId)
}

// Segments returns a slice of Resource ID Segments which comprise this Collection Partition Key Range Id ID
func (id CollectionPartitionKeyRangeIdId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDocumentDB", "Microsoft.DocumentDB", "Microsoft.DocumentDB"),
		resourceids.StaticSegment("staticDatabaseAccounts", "databaseAccounts", "databaseAccounts"),
		resourceids.UserSpecifiedSegment("accountName", "accountValue"),
		resourceids.StaticSegment("staticRegion", "region", "region"),
		resourceids.UserSpecifiedSegment("region", "regionValue"),
		resourceids.StaticSegment("staticDatabases", "databases", "databases"),
		resourceids.UserSpecifiedSegment("databaseRid", "databaseRidValue"),
		resourceids.StaticSegment("staticCollections", "collections", "collections"),
		resourceids.UserSpecifiedSegment("collectionRid", "collectionRidValue"),
		resourceids.StaticSegment("staticPartitionKeyRangeId", "partitionKeyRangeId", "partitionKeyRangeId"),
		resourceids.UserSpecifiedSegment("partitionKeyRangeId", "partitionKeyRangeIdValue"),
	}
}

// String returns a human-readable description of this Collection Partition Key Range Id ID
func (id CollectionPartitionKeyRangeIdId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Account Name: %q", id.AccountName),
		fmt.Sprintf("Region: %q", id.Region),
		fmt.Sprintf("Database Rid: %q", id.DatabaseRid),
		fmt.Sprintf("Collection Rid: %q", id.CollectionRid),
		fmt.Sprintf("Partition Key Range: %q", id.PartitionKeyRangeId),
	}
	return fmt.Sprintf("Collection Partition Key Range Id (%s)", strings.Join(components, "\n"))
}
