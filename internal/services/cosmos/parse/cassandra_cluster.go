package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type CassandraClusterId struct {
	SubscriptionId string
	ResourceGroup  string
	ClusterName    string
}

func NewCassandraClusterID(subscriptionId, resourceGroup, clusterName string) CassandraClusterId {
	return CassandraClusterId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ClusterName:    clusterName,
	}
}

func (id CassandraClusterId) String() string {
	segments := []string{
		fmt.Sprintf("Cluster Name %q", id.ClusterName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Cassandra Cluster", segmentsStr)
}

func (id CassandraClusterId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/cassandraClusters/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ClusterName)
}

// CassandraClusterID parses a Cluster ID into an CassandraClusterId struct
func CassandraClusterID(input string) (*CassandraClusterId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := CassandraClusterId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ClusterName, err = id.PopSegment("cassandraClusters"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
