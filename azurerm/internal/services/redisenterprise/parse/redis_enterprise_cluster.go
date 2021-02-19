package parse

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type RedisEnterpriseClusterId struct {
	SubscriptionId string
	ResourceGroup  string
	RediName       string
	Name           string
}

func NewRedisEnterpriseClusterID(subscriptionId, resourceGroup, name string) RedisEnterpriseClusterId {
	return RedisEnterpriseClusterId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Name:           name,
	}
}

func (id RedisEnterpriseClusterId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Redis Enterprise Cluster", segmentsStr)
}

func (id RedisEnterpriseClusterId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cache/redisEnterprise/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

// RedisEnterpriseClusterID parses a RedisEnterpriseCluster ID into an RedisEnterpriseClusterId struct
func RedisEnterpriseClusterID(input string) (*RedisEnterpriseClusterId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := RedisEnterpriseClusterId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.Name, err = id.PopSegment("redisEnterprise"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
