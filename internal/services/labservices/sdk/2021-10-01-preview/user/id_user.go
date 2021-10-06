package user

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type UserId struct {
	SubscriptionId string
	ResourceGroup  string
	LabName        string
	Name           string
}

func NewUserID(subscriptionId, resourceGroup, labName, name string) UserId {
	return UserId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		LabName:        labName,
		Name:           name,
	}
}

func (id UserId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Lab Name %q", id.LabName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "User", segmentsStr)
}

func (id UserId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.LabServices/labs/%s/users/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.LabName, id.Name)
}

// ParseUserID parses a User ID into an UserId struct
func ParseUserID(input string) (*UserId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := UserId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.LabName, err = id.PopSegment("labs"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("users"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ParseUserIDInsensitively parses an User ID into an UserId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the ParseUserID method should be used instead for validation etc.
func ParseUserIDInsensitively(input string) (*UserId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := UserId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'labs' segment
	labsKey := "labs"
	for key := range id.Path {
		if strings.EqualFold(key, labsKey) {
			labsKey = key
			break
		}
	}
	if resourceId.LabName, err = id.PopSegment(labsKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'users' segment
	usersKey := "users"
	for key := range id.Path {
		if strings.EqualFold(key, usersKey) {
			usersKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(usersKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
