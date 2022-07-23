package deletedconfigurationstores

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = DeletedConfigurationStoreId{}

// DeletedConfigurationStoreId is a struct representing the Resource ID for a Deleted Configuration Store
type DeletedConfigurationStoreId struct {
	SubscriptionId  string
	Location        string
	ConfigStoreName string
}

// NewDeletedConfigurationStoreID returns a new DeletedConfigurationStoreId struct
func NewDeletedConfigurationStoreID(subscriptionId string, location string, configStoreName string) DeletedConfigurationStoreId {
	return DeletedConfigurationStoreId{
		SubscriptionId:  subscriptionId,
		Location:        location,
		ConfigStoreName: configStoreName,
	}
}

// ParseDeletedConfigurationStoreID parses 'input' into a DeletedConfigurationStoreId
func ParseDeletedConfigurationStoreID(input string) (*DeletedConfigurationStoreId, error) {
	parser := resourceids.NewParserFromResourceIdType(DeletedConfigurationStoreId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DeletedConfigurationStoreId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.Location, ok = parsed.Parsed["location"]; !ok {
		return nil, fmt.Errorf("the segment 'location' was not found in the resource id %q", input)
	}

	if id.ConfigStoreName, ok = parsed.Parsed["configStoreName"]; !ok {
		return nil, fmt.Errorf("the segment 'configStoreName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseDeletedConfigurationStoreIDInsensitively parses 'input' case-insensitively into a DeletedConfigurationStoreId
// note: this method should only be used for API response data and not user input
func ParseDeletedConfigurationStoreIDInsensitively(input string) (*DeletedConfigurationStoreId, error) {
	parser := resourceids.NewParserFromResourceIdType(DeletedConfigurationStoreId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DeletedConfigurationStoreId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.Location, ok = parsed.Parsed["location"]; !ok {
		return nil, fmt.Errorf("the segment 'location' was not found in the resource id %q", input)
	}

	if id.ConfigStoreName, ok = parsed.Parsed["configStoreName"]; !ok {
		return nil, fmt.Errorf("the segment 'configStoreName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateDeletedConfigurationStoreID checks that 'input' can be parsed as a Deleted Configuration Store ID
func ValidateDeletedConfigurationStoreID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDeletedConfigurationStoreID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Deleted Configuration Store ID
func (id DeletedConfigurationStoreId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.AppConfiguration/locations/%s/deletedConfigurationStores/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.Location, id.ConfigStoreName)
}

// Segments returns a slice of Resource ID Segments which comprise this Deleted Configuration Store ID
func (id DeletedConfigurationStoreId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAppConfiguration", "Microsoft.AppConfiguration", "Microsoft.AppConfiguration"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("location", "locationValue"),
		resourceids.StaticSegment("staticDeletedConfigurationStores", "deletedConfigurationStores", "deletedConfigurationStores"),
		resourceids.UserSpecifiedSegment("configStoreName", "configStoreValue"),
	}
}

// String returns a human-readable description of this Deleted Configuration Store ID
func (id DeletedConfigurationStoreId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location: %q", id.Location),
		fmt.Sprintf("Config Store Name: %q", id.ConfigStoreName),
	}
	return fmt.Sprintf("Deleted Configuration Store (%s)", strings.Join(components, "\n"))
}
