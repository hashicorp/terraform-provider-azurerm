package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type DatadogMonitorId struct {
	SubscriptionId string
	ResourceGroup  string
	MonitorName    string
}

func NewDatadogMonitorID(subscriptionId, resourceGroup, monitorName string) DatadogMonitorId {
	return DatadogMonitorId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		MonitorName:    monitorName,
	}
}

func (id DatadogMonitorId) String() string {
	segments := []string{
		fmt.Sprintf("Monitor Name %q", id.MonitorName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Datadog Monitor", segmentsStr)
}

func (id DatadogMonitorId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Datadog/monitors/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.MonitorName)
}

// DatadogMonitorID parses a DatadogMonitor ID into an DatadogMonitorId struct
func DatadogMonitorID(input string) (*DatadogMonitorId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := DatadogMonitorId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.MonitorName, err = id.PopSegment("monitors"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
