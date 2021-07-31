package authorizations

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type AuthorizationId struct {
	SubscriptionId   string
	ResourceGroup    string
	PrivateCloudName string
	Name             string
}

func NewAuthorizationID(subscriptionId, resourceGroup, privateCloudName, name string) AuthorizationId {
	return AuthorizationId{
		SubscriptionId:   subscriptionId,
		ResourceGroup:    resourceGroup,
		PrivateCloudName: privateCloudName,
		Name:             name,
	}
}

func (id AuthorizationId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Private Cloud Name %q", id.PrivateCloudName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Authorization", segmentsStr)
}

func (id AuthorizationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AVS/privateClouds/%s/authorizations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.PrivateCloudName, id.Name)
}

// ParseAuthorizationID parses a Authorization ID into an AuthorizationId struct
func ParseAuthorizationID(input string) (*AuthorizationId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := AuthorizationId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.PrivateCloudName, err = id.PopSegment("privateClouds"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("authorizations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ParseAuthorizationIDInsensitively parses an Authorization ID into an AuthorizationId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the ParseAuthorizationID method should be used instead for validation etc.
func ParseAuthorizationIDInsensitively(input string) (*AuthorizationId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := AuthorizationId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'privateClouds' segment
	privateCloudsKey := "privateClouds"
	for key := range id.Path {
		if strings.EqualFold(key, privateCloudsKey) {
			privateCloudsKey = key
			break
		}
	}
	if resourceId.PrivateCloudName, err = id.PopSegment(privateCloudsKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'authorizations' segment
	authorizationsKey := "authorizations"
	for key := range id.Path {
		if strings.EqualFold(key, authorizationsKey) {
			authorizationsKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(authorizationsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
