package shares

import (
	"fmt"
)

// GetResourceManagerResourceID returns the Resource Manager specific
// ResourceID for a specific Storage Share
func (c Client) GetResourceManagerResourceID(subscriptionID, resourceGroup, accountName, shareName string) string {
	fmtStr := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/fileServices/default/shares/%s"
	return fmt.Sprintf(fmtStr, subscriptionID, resourceGroup, accountName, shareName)
}
