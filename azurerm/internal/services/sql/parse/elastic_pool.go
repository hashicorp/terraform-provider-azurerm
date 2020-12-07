package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ElasticPoolId struct {
	SubscriptionId string
	ResourceGroup  string
	ServerName     string
	Name           string
}

func NewElasticPoolID(subscriptionId, resourceGroup, serverName, name string) ElasticPoolId {
	return ElasticPoolId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ServerName:     serverName,
		Name:           name,
	}
}

func (id ElasticPoolId) String() string {
	segments := []string{
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
		fmt.Sprintf("Server Name %q", id.ServerName),
		fmt.Sprintf("Name %q", id.Name),
	}
	return strings.Join(segments, " / ")
}

func (id ElasticPoolId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql/servers/%s/elasticPools/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ServerName, id.Name)
}

// ElasticPoolID parses a ElasticPool ID into an ElasticPoolId struct
func ElasticPoolID(input string) (*ElasticPoolId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ElasticPoolId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ServerName, err = id.PopSegment("servers"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("elasticPools"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
