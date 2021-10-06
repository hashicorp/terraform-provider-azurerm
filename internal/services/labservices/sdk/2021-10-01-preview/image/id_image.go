package image

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ImageId struct {
	SubscriptionId string
	ResourceGroup  string
	LabPlanName    string
	Name           string
}

func NewImageID(subscriptionId, resourceGroup, labPlanName, name string) ImageId {
	return ImageId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		LabPlanName:    labPlanName,
		Name:           name,
	}
}

func (id ImageId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Lab Plan Name %q", id.LabPlanName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Image", segmentsStr)
}

func (id ImageId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.LabServices/labPlans/%s/images/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.LabPlanName, id.Name)
}

// ParseImageID parses a Image ID into an ImageId struct
func ParseImageID(input string) (*ImageId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ImageId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.LabPlanName, err = id.PopSegment("labPlans"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("images"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ParseImageIDInsensitively parses an Image ID into an ImageId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the ParseImageID method should be used instead for validation etc.
func ParseImageIDInsensitively(input string) (*ImageId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ImageId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'labPlans' segment
	labPlansKey := "labPlans"
	for key := range id.Path {
		if strings.EqualFold(key, labPlansKey) {
			labPlansKey = key
			break
		}
	}
	if resourceId.LabPlanName, err = id.PopSegment(labPlansKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'images' segment
	imagesKey := "images"
	for key := range id.Path {
		if strings.EqualFold(key, imagesKey) {
			imagesKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(imagesKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
