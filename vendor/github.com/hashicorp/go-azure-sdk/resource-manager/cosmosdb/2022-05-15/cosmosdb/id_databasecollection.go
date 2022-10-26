package cosmosdb

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = DatabaseCollectionId{}

// DatabaseCollectionId is a struct representing the Resource ID for a Database Collection
type DatabaseCollectionId struct {
	SubscriptionId    string
	ResourceGroupName string
	AccountName       string
	Region            string
	DatabaseRid       string
	CollectionRid     string
}

// NewDatabaseCollectionID returns a new DatabaseCollectionId struct
func NewDatabaseCollectionID(subscriptionId string, resourceGroupName string, accountName string, region string, databaseRid string, collectionRid string) DatabaseCollectionId {
	return DatabaseCollectionId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		AccountName:       accountName,
		Region:            region,
		DatabaseRid:       databaseRid,
		CollectionRid:     collectionRid,
	}
}

// ParseDatabaseCollectionID parses 'input' into a DatabaseCollectionId
func ParseDatabaseCollectionID(input string) (*DatabaseCollectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(DatabaseCollectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DatabaseCollectionId{}

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

	return &id, nil
}

// ParseDatabaseCollectionIDInsensitively parses 'input' case-insensitively into a DatabaseCollectionId
// note: this method should only be used for API response data and not user input
func ParseDatabaseCollectionIDInsensitively(input string) (*DatabaseCollectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(DatabaseCollectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DatabaseCollectionId{}

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

	return &id, nil
}

// ValidateDatabaseCollectionID checks that 'input' can be parsed as a Database Collection ID
func ValidateDatabaseCollectionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDatabaseCollectionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Database Collection ID
func (id DatabaseCollectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/databaseAccounts/%s/region/%s/databases/%s/collections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.Region, id.DatabaseRid, id.CollectionRid)
}

// Segments returns a slice of Resource ID Segments which comprise this Database Collection ID
func (id DatabaseCollectionId) Segments() []resourceids.Segment {
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
	}
}

// String returns a human-readable description of this Database Collection ID
func (id DatabaseCollectionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Account Name: %q", id.AccountName),
		fmt.Sprintf("Region: %q", id.Region),
		fmt.Sprintf("Database Rid: %q", id.DatabaseRid),
		fmt.Sprintf("Collection Rid: %q", id.CollectionRid),
	}
	return fmt.Sprintf("Database Collection (%s)", strings.Join(components, "\n"))
}
