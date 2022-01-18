package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ElasticMonitorId struct {
	SubscriptionId string
	ResourceGroup  string
	MonitorName    string
}

func NewElasticMonitorID(subscriptionId, resourceGroup, monitorName string) ElasticMonitorId {
	return ElasticMonitorId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		MonitorName:    monitorName,
	}
}

func (id ElasticMonitorId) String() string {
	segments := []string{
		fmt.Sprintf("Monitor Name %q", id.MonitorName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Elastic Monitor", segmentsStr)
}

func (id ElasticMonitorId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Elastic/monitors/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.MonitorName)
}

// ElasticMonitorID parses a ElasticMonitor ID into an ElasticMonitorId struct
func ElasticMonitorID(input string) (*ElasticMonitorId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ElasticMonitorId{
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
