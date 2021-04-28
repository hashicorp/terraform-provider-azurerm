package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type RedisEnterpriseDatabaseId struct {
	SubscriptionId      string
	ResourceGroup       string
	RedisEnterpriseName string
	DatabaseName        string
}

func NewRedisEnterpriseDatabaseID(subscriptionId, resourceGroup, redisEnterpriseName, databaseName string) RedisEnterpriseDatabaseId {
	return RedisEnterpriseDatabaseId{
		SubscriptionId:      subscriptionId,
		ResourceGroup:       resourceGroup,
		RedisEnterpriseName: redisEnterpriseName,
		DatabaseName:        databaseName,
	}
}

func (id RedisEnterpriseDatabaseId) String() string {
	segments := []string{
		fmt.Sprintf("Database Name %q", id.DatabaseName),
		fmt.Sprintf("Redis Enterprise Name %q", id.RedisEnterpriseName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Redis Enterprise Database", segmentsStr)
}

func (id RedisEnterpriseDatabaseId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cache/redisEnterprise/%s/databases/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.RedisEnterpriseName, id.DatabaseName)
}

// RedisEnterpriseDatabaseID parses a RedisEnterpriseDatabase ID into an RedisEnterpriseDatabaseId struct
func RedisEnterpriseDatabaseID(input string) (*RedisEnterpriseDatabaseId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := RedisEnterpriseDatabaseId{
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
	if resourceId.DatabaseName, err = id.PopSegment("databases"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
