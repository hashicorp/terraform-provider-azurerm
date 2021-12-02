package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type FlexibleDatabaseId struct {
	SubscriptionId     string
	ResourceGroup      string
	FlexibleServerName string
	DatabaseName       string
}

func NewFlexibleDatabaseID(subscriptionId, resourceGroup, flexibleServerName, databaseName string) FlexibleDatabaseId {
	return FlexibleDatabaseId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		FlexibleServerName: flexibleServerName,
		DatabaseName:       databaseName,
	}
}

func (id FlexibleDatabaseId) String() string {
	segments := []string{
		fmt.Sprintf("Database Name %q", id.DatabaseName),
		fmt.Sprintf("Flexible Server Name %q", id.FlexibleServerName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Flexible Database", segmentsStr)
}

func (id FlexibleDatabaseId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DBforMySQL/flexibleServers/%s/databases/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.FlexibleServerName, id.DatabaseName)
}

// FlexibleDatabaseID parses a FlexibleDatabase ID into an FlexibleDatabaseId struct
func FlexibleDatabaseID(input string) (*FlexibleDatabaseId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FlexibleDatabaseId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.FlexibleServerName, err = id.PopSegment("flexibleServers"); err != nil {
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
