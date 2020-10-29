package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ImageBuilderTemplateId struct {
	Name          string
	ResourceGroup string
}

func NewImageBuilderTemplateID(resourceGroup, name string) ImageBuilderTemplateId {
	return ImageBuilderTemplateId{
		ResourceGroup: resourceGroup,
		Name:          name,
	}
}

func ImageBuilderTemplateID(input string) (*ImageBuilderTemplateId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Image Builder Template ID %q: %+v", input, err)
	}

	imageBuilderTemplate := ImageBuilderTemplateId{
		ResourceGroup: id.ResourceGroup,
	}

	if imageBuilderTemplate.Name, err = id.PopSegment("imageTemplates"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &imageBuilderTemplate, nil
}

func (id ImageBuilderTemplateId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.VirtualMachineImages/imageTemplates/%s",
		subscriptionId, id.ResourceGroup, id.Name)
}
