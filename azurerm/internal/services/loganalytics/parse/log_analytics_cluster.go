package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type LogAnalyticsClusterId struct {
	SubscriptionId string
	ResourceGroup  string
	ClusterName    string
}

func NewLogAnalyticsClusterID(subscriptionId, resourceGroup, clusterName string) LogAnalyticsClusterId {
	return LogAnalyticsClusterId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ClusterName:    clusterName,
	}
}

func (id LogAnalyticsClusterId) String() string {
	segments := []string{
		fmt.Sprintf("Cluster Name %q", id.ClusterName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Log Analytics Cluster", segmentsStr)
}

func (id LogAnalyticsClusterId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/microsoft.operationalinsights/clusters/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ClusterName)
}

// LogAnalyticsClusterID parses a LogAnalyticsCluster ID into an LogAnalyticsClusterId struct
func LogAnalyticsClusterID(input string) (*LogAnalyticsClusterId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := LogAnalyticsClusterId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ClusterName, err = id.PopSegment("clusters"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
