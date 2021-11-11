package cognitiveservicesaccounts

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type DeletedAccountId struct {
	SubscriptionId string
	LocationName   string
	ResourceGroup  string
	Name           string
}

func NewDeletedAccountID(subscriptionId, locationName, resourceGroup, name string) DeletedAccountId {
	return DeletedAccountId{
		SubscriptionId: subscriptionId,
		LocationName:   locationName,
		ResourceGroup:  resourceGroup,
		Name:           name,
	}
}

func (id DeletedAccountId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
		fmt.Sprintf("Location Name %q", id.LocationName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Deleted Account", segmentsStr)
}

func (id DeletedAccountId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.CognitiveServices/locations/%s/resourceGroups/%s/deletedAccounts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.ResourceGroup, id.Name)
}

// ParseDeletedAccountID parses a DeletedAccount ID into an DeletedAccountId struct
func ParseDeletedAccountID(input string) (*DeletedAccountId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := DeletedAccountId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.LocationName, err = id.PopSegment("locations"); err != nil {
		return nil, err
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.Name, err = id.PopSegment("deletedAccounts"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ParseDeletedAccountIDInsensitively parses an DeletedAccount ID into an DeletedAccountId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the ParseDeletedAccountID method should be used instead for validation etc.
func ParseDeletedAccountIDInsensitively(input string) (*DeletedAccountId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := DeletedAccountId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	// find the correct casing for the 'locations' segment
	locationsKey := "locations"
	for key := range id.Path {
		if strings.EqualFold(key, locationsKey) {
			locationsKey = key
			break
		}
	}
	if resourceId.LocationName, err = id.PopSegment(locationsKey); err != nil {
		return nil, err
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'deletedAccounts' segment
	deletedAccountsKey := "deletedAccounts"
	for key := range id.Path {
		if strings.EqualFold(key, deletedAccountsKey) {
			deletedAccountsKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(deletedAccountsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
