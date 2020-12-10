package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AppServicePlanId struct {
	SubscriptionId string
	ResourceGroup  string
	ServerfarmName string
}

func NewAppServicePlanID(subscriptionId, resourceGroup, serverfarmName string) AppServicePlanId {
	return AppServicePlanId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ServerfarmName: serverfarmName,
	}
}

func (id AppServicePlanId) String() string {
	segments := []string{
		fmt.Sprintf("Serverfarm Name %q", id.ServerfarmName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "App Service Plan", segmentsStr)
}

func (id AppServicePlanId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/serverfarms/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ServerfarmName)
}

// AppServicePlanID parses a AppServicePlan ID into an AppServicePlanId struct
func AppServicePlanID(input string) (*AppServicePlanId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := AppServicePlanId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ServerfarmName, err = id.PopSegment("serverfarms"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
