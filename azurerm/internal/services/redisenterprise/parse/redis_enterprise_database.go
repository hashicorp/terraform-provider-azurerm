package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type RedisEnterpriseDatabaseId struct {
	SubscriptionId string
	ResourceGroup  string
	ClusterName    string
	Name           string
}

func NewRedisEnterpriseDatabaseID(subscriptionId string, resourcegroup string, clustername string, name string) RedisEnterpriseDatabaseId {
	return RedisEnterpriseDatabaseId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourcegroup,
		ClusterName:    clustername,
		Name:           name,
	}
}

func (id RedisEnterpriseDatabaseId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cache/redisEnterprise/%s/databases/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ClusterName, id.Name)
}

func RedisEnterpriseDatabaseID(input string) (*RedisEnterpriseDatabaseId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Redis Enterprise Database ID %q: %+v", input, err)
	}

	redisenterpriseDatabase := RedisEnterpriseDatabaseId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if redisenterpriseDatabase.ClusterName, err = id.PopSegment("redisEnterprise"); err != nil {
		return nil, err
	}

	if redisenterpriseDatabase.Name, err = id.PopSegment("databases"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &redisenterpriseDatabase, nil
}
