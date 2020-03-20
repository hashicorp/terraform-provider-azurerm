package shares

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
)

// GetResourceID returns the Resource ID for the given File Share
// This can be useful when, for example, you're using this as a unique identifier
func (client Client) GetResourceID(accountName, shareName string) string {
	domain := endpoints.GetFileEndpoint(client.BaseURI, accountName)
	return fmt.Sprintf("%s/%s", domain, shareName)
}

// GetResourceManagerResourceID returns the Resource Manager specific
// ResourceID for a specific Storage Share
func (client Client) GetResourceManagerResourceID(subscriptionID, resourceGroup, accountName, shareName string) string {
	fmtStr := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/fileServices/default/shares/%s"
	return fmt.Sprintf(fmtStr, subscriptionID, resourceGroup, accountName, shareName)
}

type ResourceID struct {
	AccountName string
	ShareName   string
}

// ParseResourceID parses the specified Resource ID and returns an object
// which can be used to interact with the Storage Shares SDK
func ParseResourceID(id string) (*ResourceID, error) {
	// example: https://foo.file.core.windows.net/Bar
	if id == "" {
		return nil, fmt.Errorf("`id` was empty")
	}

	uri, err := url.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("Error parsing ID as a URL: %s", err)
	}

	accountName, err := endpoints.GetAccountNameFromEndpoint(uri.Host)
	if err != nil {
		return nil, fmt.Errorf("Error parsing Account Name: %s", err)
	}

	shareName := strings.TrimPrefix(uri.Path, "/")
	return &ResourceID{
		AccountName: *accountName,
		ShareName:   shareName,
	}, nil
}
