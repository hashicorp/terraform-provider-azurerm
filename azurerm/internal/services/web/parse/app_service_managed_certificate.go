package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AppServiceManagedCertificateId struct {
	Name          string
	ResourceGroup string
}

func NewAppServiceManagedCertificateId(name, resourceGroup string) AppServiceManagedCertificateId {
	return AppServiceManagedCertificateId{
		Name:          name,
		ResourceGroup: resourceGroup,
	}
}

func (id AppServiceManagedCertificateId) ID(subscriptionId string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/certificates/%s"
	return fmt.Sprintf(fmtString, subscriptionId, id.ResourceGroup, id.Name)
}

func AppServiceManagedCertificateID(input string) (*AppServiceManagedCertificateId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing ID for App Service Managed Certificate %q: %+v", input, err)
	}

	appServiceManagedCertificateId := AppServiceManagedCertificateId{
		ResourceGroup: id.ResourceGroup,
	}

	if appServiceManagedCertificateId.Name, err = id.PopSegment("certificates"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &appServiceManagedCertificateId, nil
}
