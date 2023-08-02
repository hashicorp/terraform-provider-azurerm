package filesystems

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
)

// GetResourceID returns the Resource ID for the given Data Lake Storage FileSystem
// This can be useful when, for example, you're using this as a unique identifier
func (client Client) GetResourceID(accountName, shareName string) string {
	domain := endpoints.GetDataLakeStoreEndpoint(client.BaseURI, accountName)
	return fmt.Sprintf("%s/%s", domain, shareName)
}

type ResourceID struct {
	AccountName   string
	DirectoryName string
}

// ParseResourceID parses the specified Resource ID and returns an object
// which can be used to interact with the Data Lake Storage FileSystem API's
func ParseResourceID(id string) (*ResourceID, error) {
	// example: https://foo.dfs.core.windows.net/Bar
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

	directoryName := strings.TrimPrefix(uri.Path, "/")
	return &ResourceID{
		AccountName:   *accountName,
		DirectoryName: directoryName,
	}, nil
}
