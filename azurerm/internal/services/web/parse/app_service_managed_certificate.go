package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AppServiceManagedCertificateId struct {
	SubscriptionId string
	ResourceGroup  string
	Name           string
}

func NewAppServiceManagedCertificateId(subscriptionId, resourceGroup, name string) AppServiceManagedCertificateId {
	return AppServiceManagedCertificateId{
		SubscriptionId: subscriptionId,
		Name:           name,
		ResourceGroup:  resourceGroup,
	}
}

func (id AppServiceManagedCertificateId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/certificates/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

func AppServiceManagedCertificateID(input string) (*AppServiceManagedCertificateId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing ID for App Service Managed Certificate %q: %+v", input, err)
	}

	appServiceManagedCertificateId := AppServiceManagedCertificateId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if appServiceManagedCertificateId.Name, err = id.PopSegment("certificates"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &appServiceManagedCertificateId, nil
}
