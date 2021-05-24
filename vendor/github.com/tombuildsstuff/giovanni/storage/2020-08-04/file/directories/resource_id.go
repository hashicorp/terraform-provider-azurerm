package directories

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
)

// GetResourceID returns the Resource ID for the given Directory
// This can be useful when, for example, you're using this as a unique identifier
func (client Client) GetResourceID(accountName, shareName, directoryName string) string {
	domain := endpoints.GetFileEndpoint(client.BaseURI, accountName)
	return fmt.Sprintf("%s/%s/%s", domain, shareName, directoryName)
}

type ResourceID struct {
	AccountName   string
	DirectoryName string
	ShareName     string
}

// ParseResourceID parses the Resource ID into an Object
// which can be used to interact with the Directory within the File Share
func ParseResourceID(id string) (*ResourceID, error) {
	// example: https://foo.file.core.windows.net/Bar/Folder
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

	path := strings.TrimPrefix(uri.Path, "/")
	segments := strings.Split(path, "/")
	if len(segments) == 0 {
		return nil, fmt.Errorf("Expected the path to contain segments but got none")
	}

	shareName := segments[0]
	directoryName := strings.TrimPrefix(path, shareName)
	directoryName = strings.TrimPrefix(directoryName, "/")
	return &ResourceID{
		AccountName:   *accountName,
		ShareName:     shareName,
		DirectoryName: directoryName,
	}, nil
}
