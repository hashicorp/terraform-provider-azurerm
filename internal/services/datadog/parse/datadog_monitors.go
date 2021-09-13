package parse

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DatadogMonitorId struct {
    SubscriptionId string
    ResourceGroup string
    Name string
}

func NewDatadogMonitorID(subscriptionId, resourcegroup, name string) 
DatadogMonitorId {
    return DatadogMonitorId{
        SubscriptionId: subscriptionId,
        ResourceGroup: resourcegroup,
        Name: name,
    }
}

func (id DatadogMonitorId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Datadog", segmentsStr)
}

func (id DatadogMonitorId) ID() string {
    fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Datadog/monitors/%s"
    return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

// DatadogMonitorID parses a DatadogMonitor ID into an DatadogMonitorId struct
func DatadogMonitorID(input string) (*DatadogMonitorId, error) {
    id, err := azure.ParseAzureResourceID(input)
    if err != nil {
        return nil, err
    }

    resourceId := DatadogMonitorId{
        SubscriptionId: id.SubscriptionID,
        ResourceGroup: id.ResourceGroup,
    }

    if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

    if resourceId.Name, err = id.PopSegment("DatadogMonitors"); err != nil {
        return nil, err
    }

    if err := id.ValidateNoEmptySegments(input); err != nil {
        return nil, err
    }

    return &resourceId, nil
}