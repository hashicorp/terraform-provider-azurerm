package parse

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type PortalTenantConfigurationId struct {
	Name string
}

func NewPortalTenantConfigurationID(name string) PortalTenantConfigurationId {
	return PortalTenantConfigurationId{
		Name: name,
	}
}

func (id PortalTenantConfigurationId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Portal Tenant Configuration", segmentsStr)
}

func (id PortalTenantConfigurationId) ID() string {
	fmtString := "/providers/Microsoft.Portal/tenantConfigurations/%s"
	return fmt.Sprintf(fmtString, id.Name)
}

// PortalTenantConfigurationID parses a PortalTenantConfiguration ID into an PortalTenantConfigurationId struct
func PortalTenantConfigurationID(input string) (*PortalTenantConfigurationId, error) {
	id, err := azure.ParseAzureResourceIDWithoutSubscription(input)
	if err != nil {
		return nil, err
	}

	resourceId := PortalTenantConfigurationId{}

	if resourceId.Name, err = id.PopSegment("tenantConfigurations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
