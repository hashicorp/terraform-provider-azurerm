package filesystems

import (
	"fmt"
)

// GetResourceManagerResourceID returns the Resource Manager specific
// ResourceID for a specific dfs filesystem
func (c Client) GetResourceManagerResourceID(subscriptionID, resourceGroup, accountName, fileSystemName string) string {
	fmtStr := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/blobServices/default/containers/%s"
	return fmt.Sprintf(fmtStr, subscriptionID, resourceGroup, accountName, fileSystemName)
}
