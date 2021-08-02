package privateclouds

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type PrivateCloudId struct {
	SubscriptionId string
	ResourceGroup  string
	Name           string
}

func NewPrivateCloudID(subscriptionId, resourceGroup, name string) PrivateCloudId {
	return PrivateCloudId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Name:           name,
	}
}

func (id PrivateCloudId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Private Cloud", segmentsStr)
}

func (id PrivateCloudId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AVS/privateClouds/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

// ParsePrivateCloudID parses a PrivateCloud ID into an PrivateCloudId struct
func ParsePrivateCloudID(input string) (*PrivateCloudId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := PrivateCloudId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.Name, err = id.PopSegment("privateClouds"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ParsePrivateCloudIDInsensitively parses an PrivateCloud ID into an PrivateCloudId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the ParsePrivateCloudID method should be used instead for validation etc.
func ParsePrivateCloudIDInsensitively(input string) (*PrivateCloudId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := PrivateCloudId{
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
	if resourceId.Name, err = id.PopSegment(privateCloudsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
