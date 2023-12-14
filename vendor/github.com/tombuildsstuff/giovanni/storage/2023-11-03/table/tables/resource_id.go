package tables

import "fmt"

// GetResourceManagerResourceID returns the Resource ID for the given Table
// This can be useful when, for example, you're using this as a unique identifier
func (c Client) GetResourceManagerResourceID(subscriptionID, resourceGroup, accountName, tableName string) string {
	fmtStr := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/tableServices/default/tables/%s"
	return fmt.Sprintf(fmtStr, subscriptionID, resourceGroup, accountName, tableName)
}
