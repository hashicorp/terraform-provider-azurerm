package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type OrderId struct {
	SubscriptionId        string
	ResourceGroup         string
	DataBoxEdgeDeviceName string
	Name                  string
}

func NewOrderID(subscriptionId, resourceGroup, dataBoxEdgeDeviceName, name string) OrderId {
	return OrderId{
		SubscriptionId:        subscriptionId,
		ResourceGroup:         resourceGroup,
		DataBoxEdgeDeviceName: dataBoxEdgeDeviceName,
		Name:                  name,
	}
}

func (id OrderId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Data Box Edge Device Name %q", id.DataBoxEdgeDeviceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Order", segmentsStr)
}

func (id OrderId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataBoxEdge/dataBoxEdgeDevices/%s/orders/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.DataBoxEdgeDeviceName, id.Name)
}

// OrderID parses a Order ID into an OrderId struct
func OrderID(input string) (*OrderId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := OrderId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.DataBoxEdgeDeviceName, err = id.PopSegment("dataBoxEdgeDevices"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("orders"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
