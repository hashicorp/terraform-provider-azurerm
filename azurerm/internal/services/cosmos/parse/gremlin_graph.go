package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type GremlinGraphId struct {
	SubscriptionId      string
	ResourceGroup       string
	DatabaseAccountName string
	GremlinDatabaseName string
	GraphName           string
}

func NewGremlinGraphID(subscriptionId, resourceGroup, databaseAccountName, gremlinDatabaseName, graphName string) GremlinGraphId {
	return GremlinGraphId{
		SubscriptionId:      subscriptionId,
		ResourceGroup:       resourceGroup,
		DatabaseAccountName: databaseAccountName,
		GremlinDatabaseName: gremlinDatabaseName,
		GraphName:           graphName,
	}
}

func (id GremlinGraphId) String() string {
	segments := []string{
		fmt.Sprintf("Graph Name %q", id.GraphName),
		fmt.Sprintf("Gremlin Database Name %q", id.GremlinDatabaseName),
		fmt.Sprintf("Database Account Name %q", id.DatabaseAccountName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Gremlin Graph", segmentsStr)
}

func (id GremlinGraphId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/databaseAccounts/%s/gremlinDatabases/%s/graphs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.DatabaseAccountName, id.GremlinDatabaseName, id.GraphName)
}

// GremlinGraphID parses a GremlinGraph ID into an GremlinGraphId struct
func GremlinGraphID(input string) (*GremlinGraphId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := GremlinGraphId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.DatabaseAccountName, err = id.PopSegment("databaseAccounts"); err != nil {
		return nil, err
	}
	if resourceId.GremlinDatabaseName, err = id.PopSegment("gremlinDatabases"); err != nil {
		return nil, err
	}
	if resourceId.GraphName, err = id.PopSegment("graphs"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
