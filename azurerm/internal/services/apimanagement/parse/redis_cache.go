package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type RedisCacheId struct {
	SubscriptionId string
	ResourceGroup  string
	ServiceName    string
	CacheName      string
}

func NewRedisCacheID(subscriptionId, resourceGroup, serviceName, cacheName string) RedisCacheId {
	return RedisCacheId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ServiceName:    serviceName,
		CacheName:      cacheName,
	}
}

func (id RedisCacheId) String() string {
	segments := []string{
		fmt.Sprintf("Cache Name %q", id.CacheName),
		fmt.Sprintf("Service Name %q", id.ServiceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Redis Cache", segmentsStr)
}

func (id RedisCacheId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/caches/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ServiceName, id.CacheName)
}

// RedisCacheID parses a RedisCache ID into an RedisCacheId struct
func RedisCacheID(input string) (*RedisCacheId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := RedisCacheId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ServiceName, err = id.PopSegment("service"); err != nil {
		return nil, err
	}
	if resourceId.CacheName, err = id.PopSegment("caches"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
