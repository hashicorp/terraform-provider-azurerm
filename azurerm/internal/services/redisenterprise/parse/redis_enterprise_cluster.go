package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type RedisEnterpriseClusterId struct {
	SubscriptionId      string
	ResourceGroup       string
	RedisEnterpriseName string
}

func NewRedisEnterpriseClusterID(subscriptionId, resourceGroup, redisEnterpriseName string) RedisEnterpriseClusterId {
	return RedisEnterpriseClusterId{
		SubscriptionId:      subscriptionId,
		ResourceGroup:       resourceGroup,
		RedisEnterpriseName: redisEnterpriseName,
	}
}

func (id RedisEnterpriseClusterId) String() string {
	segments := []string{
		fmt.Sprintf("Redis Enterprise Name %q", id.RedisEnterpriseName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Redis Enterprise Cluster", segmentsStr)
}

func (id RedisEnterpriseClusterId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cache/redisEnterprise/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.RedisEnterpriseName)
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

	if resourceId.RedisEnterpriseName, err = id.PopSegment("redisEnterprise"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
