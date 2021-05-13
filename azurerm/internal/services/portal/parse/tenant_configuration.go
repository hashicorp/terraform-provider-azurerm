package parse

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type TenantConfigurationId struct {
	Name string
}

func NewTenantConfigurationID(name string) TenantConfigurationId {
	return TenantConfigurationId{
		Name: name,
	}
}

func (id TenantConfigurationId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Tenant Configuration", segmentsStr)
}

func (id TenantConfigurationId) ID() string {
	fmtString := "/providers/Microsoft.Portal/tenantConfigurations/%s"
	return fmt.Sprintf(fmtString, id.Name)
}

// TenantConfigurationID parses a TenantConfiguration ID into an TenantConfigurationId struct
func TenantConfigurationID(input string) (*TenantConfigurationId, error) {
	id, err := azure.ParseAzureResourceIDWithoutSubscription(input)
	if err != nil {
		return nil, err
	}

	resourceId := TenantConfigurationId{}

	if resourceId.Name, err = id.PopSegment("tenantConfigurations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
