package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ImageBuilderTemplateId struct {
	SubscriptionId    string
	ResourceGroup     string
	ImageTemplateName string
}

func NewImageBuilderTemplateID(subscriptionId, resourceGroup, imageTemplateName string) ImageBuilderTemplateId {
	return ImageBuilderTemplateId{
		SubscriptionId:    subscriptionId,
		ResourceGroup:     resourceGroup,
		ImageTemplateName: imageTemplateName,
	}
}

func (id ImageBuilderTemplateId) String() string {
	segments := []string{
		fmt.Sprintf("Image Template Name %q", id.ImageTemplateName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Image Builder Template", segmentsStr)
}

func (id ImageBuilderTemplateId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.VirtualMachineImages/imageTemplates/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ImageTemplateName)
}

// ImageBuilderTemplateID parses a ImageBuilderTemplate ID into an ImageBuilderTemplateId struct
func ImageBuilderTemplateID(input string) (*ImageBuilderTemplateId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ImageBuilderTemplateId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ImageTemplateName, err = id.PopSegment("imageTemplates"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
