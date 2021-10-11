package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type AnalyticsUserItemId struct {
	SubscriptionId      string
	ResourceGroup       string
	ComponentName       string
	MyanalyticsItemName string
}

func NewAnalyticsUserItemID(subscriptionId, resourceGroup, componentName, myanalyticsItemName string) AnalyticsUserItemId {
	return AnalyticsUserItemId{
		SubscriptionId:      subscriptionId,
		ResourceGroup:       resourceGroup,
		ComponentName:       componentName,
		MyanalyticsItemName: myanalyticsItemName,
	}
}

func (id AnalyticsUserItemId) String() string {
	segments := []string{
		fmt.Sprintf("Myanalytics Item Name %q", id.MyanalyticsItemName),
		fmt.Sprintf("Component Name %q", id.ComponentName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Analytics User Item", segmentsStr)
}

func (id AnalyticsUserItemId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/microsoft.insights/components/%s/myanalyticsItems/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ComponentName, id.MyanalyticsItemName)
}

// AnalyticsUserItemID parses a AnalyticsUserItem ID into an AnalyticsUserItemId struct
func AnalyticsUserItemID(input string) (*AnalyticsUserItemId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := AnalyticsUserItemId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ComponentName, err = id.PopSegment("components"); err != nil {
		return nil, err
	}
	if resourceId.MyanalyticsItemName, err = id.PopSegment("myanalyticsItems"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
