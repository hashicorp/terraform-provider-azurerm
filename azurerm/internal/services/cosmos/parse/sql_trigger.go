package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SqlTriggerId struct {
	SubscriptionId      string
	ResourceGroup       string
	DatabaseAccountName string
	SqlDatabaseName     string
	ContainerName       string
	TriggerName         string
}

func NewSqlTriggerID(subscriptionId, resourceGroup, databaseAccountName, sqlDatabaseName, containerName, triggerName string) SqlTriggerId {
	return SqlTriggerId{
		SubscriptionId:      subscriptionId,
		ResourceGroup:       resourceGroup,
		DatabaseAccountName: databaseAccountName,
		SqlDatabaseName:     sqlDatabaseName,
		ContainerName:       containerName,
		TriggerName:         triggerName,
	}
}

func (id SqlTriggerId) String() string {
	segments := []string{
		fmt.Sprintf("Trigger Name %q", id.TriggerName),
		fmt.Sprintf("Container Name %q", id.ContainerName),
		fmt.Sprintf("Sql Database Name %q", id.SqlDatabaseName),
		fmt.Sprintf("Database Account Name %q", id.DatabaseAccountName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Sql Trigger", segmentsStr)
}

func (id SqlTriggerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/databaseAccounts/%s/sqlDatabases/%s/containers/%s/triggers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.DatabaseAccountName, id.SqlDatabaseName, id.ContainerName, id.TriggerName)
}

// SqlTriggerID parses a SqlTrigger ID into an SqlTriggerId struct
func SqlTriggerID(input string) (*SqlTriggerId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SqlTriggerId{
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
	if resourceId.SqlDatabaseName, err = id.PopSegment("sqlDatabases"); err != nil {
		return nil, err
	}
	if resourceId.ContainerName, err = id.PopSegment("containers"); err != nil {
		return nil, err
	}
	if resourceId.TriggerName, err = id.PopSegment("triggers"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
