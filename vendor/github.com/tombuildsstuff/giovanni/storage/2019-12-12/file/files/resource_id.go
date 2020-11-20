package files

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
)

// GetResourceID returns the Resource ID for the given File
// This can be useful when, for example, you're using this as a unique identifier
func (client Client) GetResourceID(accountName, shareName, directoryName, filePath string) string {
	domain := endpoints.GetFileEndpoint(client.BaseURI, accountName)
	return fmt.Sprintf("%s/%s/%s/%s", domain, shareName, directoryName, filePath)
}

type ResourceID struct {
	AccountName   string
	DirectoryName string
	FileName      string
	ShareName     string
}

// ParseResourceID parses the specified Resource ID and returns an object
// which can be used to interact with Files within a Storage Share.
func ParseResourceID(id string) (*ResourceID, error) {
	// example: https://account1.file.core.chinacloudapi.cn/share1/directory1/file1.txt
	// example: https://account1.file.core.chinacloudapi.cn/share1/directory1/directory2/file1.txt

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
	fileName := segments[len(segments)-1]

	directoryName := strings.TrimPrefix(path, shareName)
	directoryName = strings.TrimPrefix(directoryName, "/")
	directoryName = strings.TrimSuffix(directoryName, fileName)
	directoryName = strings.TrimSuffix(directoryName, "/")
	return &ResourceID{
		AccountName:   *accountName,
		ShareName:     shareName,
		DirectoryName: directoryName,
		FileName:      fileName,
	}, nil
}
