package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type CassandraDatacenterId struct {
	SubscriptionId string
	ResourceGroup  string
	ClusterName    string
	DatacenterName string
}

func NewCassandraDatacenterID(subscriptionId, resourceGroup, clusterName, DatacenterName string) CassandraDatacenterId {
	return CassandraDatacenterId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ClusterName:    clusterName,
		DatacenterName: DatacenterName,
	}
}

func (id CassandraDatacenterId) String() string {
	segments := []string{
		fmt.Sprintf("Datacenter Name %q", id.DatacenterName),
		fmt.Sprintf("Cluster Name %q", id.ClusterName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	log.Println("************** segmentsStr: " + segmentsStr)
	return fmt.Sprintf("%s: (%s)", "Cassandra Cluster", segmentsStr)
}

func (id CassandraDatacenterId) ID() string {
	fmtString := "/subscriptions/{%s}/resourceGroups/%s/providers/Microsoft.DocumentDB/cassandraClusters/%s/dataCenters/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ClusterName, id.DatacenterName)
}

func CassandraDatacenterID(input string) (*CassandraDatacenterId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := CassandraDatacenterId{
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

	if resourceId.DatacenterName, err = id.PopSegment("dataCenters"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
