package tables

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
)

// GetResourceID returns the Resource ID for the given Table
// This can be useful when, for example, you're using this as a unique identifier
func (client Client) GetResourceID(accountName, tableName string) string {
	domain := endpoints.GetTableEndpoint(client.BaseURI, accountName)
	return fmt.Sprintf("%s/Tables('%s')", domain, tableName)
}

type ResourceID struct {
	AccountName string
	TableName   string
}

// ParseResourceID parses the Resource ID and returns an object which
// can be used to interact with the Table within the specified Storage Account
func ParseResourceID(id string) (*ResourceID, error) {
	// example: https://foo.table.core.windows.net/Table('foo')
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

	// assume there a `Table('')`
	path := strings.TrimPrefix(uri.Path, "/")
	if !strings.HasPrefix(path, "Tables('") || !strings.HasSuffix(path, "')") {
		return nil, fmt.Errorf("Expected the Table Name to be in the format `Tables('name')` but got %q", path)
	}

	// strip off the `Table('')`
	tableName := strings.TrimPrefix(uri.Path, "/Tables('")
	tableName = strings.TrimSuffix(tableName, "')")
	return &ResourceID{
		AccountName: *accountName,
		TableName:   tableName,
	}, nil
}
