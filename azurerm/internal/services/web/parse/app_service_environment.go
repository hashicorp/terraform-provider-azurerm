package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AppServiceEnvironmentId struct {
	SubscriptionId         string
	ResourceGroup          string
	HostingEnvironmentName string
}

func NewAppServiceEnvironmentID(subscriptionId, resourceGroup, hostingEnvironmentName string) AppServiceEnvironmentId {
	return AppServiceEnvironmentId{
		SubscriptionId:         subscriptionId,
		ResourceGroup:          resourceGroup,
		HostingEnvironmentName: hostingEnvironmentName,
	}
}

func (id AppServiceEnvironmentId) String() string {
	segments := []string{
		fmt.Sprintf("Hosting Environment Name %q", id.HostingEnvironmentName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "App Service Environment", segmentsStr)
}

func (id AppServiceEnvironmentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/hostingEnvironments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.HostingEnvironmentName)
}

// AppServiceEnvironmentID parses a AppServiceEnvironment ID into an AppServiceEnvironmentId struct
func AppServiceEnvironmentID(input string) (*AppServiceEnvironmentId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := AppServiceEnvironmentId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.HostingEnvironmentName, err = id.PopSegment("hostingEnvironments"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
