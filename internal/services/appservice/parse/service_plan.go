package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type ServicePlanId struct {
	SubscriptionId string
	ResourceGroup  string
	ServerfarmName string
}

func NewServicePlanID(subscriptionId, resourceGroup, serverfarmName string) ServicePlanId {
	return ServicePlanId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ServerfarmName: serverfarmName,
	}
}

func (id ServicePlanId) String() string {
	segments := []string{
		fmt.Sprintf("Serverfarm Name %q", id.ServerfarmName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Service Plan", segmentsStr)
}

func (id ServicePlanId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/serverfarms/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ServerfarmName)
}

// ServicePlanID parses a ServicePlan ID into an ServicePlanId struct
func ServicePlanID(input string) (*ServicePlanId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ServicePlanId{
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
