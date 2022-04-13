package entities

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
)

// GetResourceID returns the Resource ID for the given Entity
// This can be useful when, for example, you're using this as a unique identifier
func (client Client) GetResourceID(accountName, tableName, partitionKey, rowKey string) string {
	domain := endpoints.GetTableEndpoint(client.BaseURI, accountName)
	return fmt.Sprintf("%s/%s(PartitionKey='%s',RowKey='%s')", domain, tableName, partitionKey, rowKey)
}

type ResourceID struct {
	AccountName  string
	TableName    string
	PartitionKey string
	RowKey       string
}

// ParseResourceID parses the specified Resource ID and returns an object which
// can be used to look up the specified Entity within the specified Table
func ParseResourceID(id string) (*ResourceID, error) {
	// example: https://account1.table.core.chinacloudapi.cn/table1(PartitionKey='partition1',RowKey='row1')
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
	if !strings.Contains(uri.Path, "(") || !strings.HasSuffix(uri.Path, ")") {
		return nil, fmt.Errorf("Expected the Table Name to be in the format `tables(PartitionKey='',RowKey='')` but got %q", path)
	}

	// NOTE: honestly this could probably be a RegEx, but this seemed like the simplest way to
	// allow these two fields to be specified in either order
	indexOfBracket := strings.IndexByte(path, '(')
	tableName := path[0:indexOfBracket]

	// trim off the brackets
	temp := strings.TrimPrefix(path, fmt.Sprintf("%s(", tableName))
	temp = strings.TrimSuffix(temp, ")")

	dictionary := strings.Split(temp, ",")
	partitionKey := ""
	rowKey := ""
	for _, v := range dictionary {
		split := strings.Split(v, "=")
		if len(split) != 2 {
			return nil, fmt.Errorf("Expected 2 segments but got %d for %q", len(split), v)
		}

		key := split[0]
		value := strings.TrimSuffix(strings.TrimPrefix(split[1], "'"), "'")
		if strings.EqualFold(key, "PartitionKey") {
			partitionKey = value
		} else if strings.EqualFold(key, "RowKey") {
			rowKey = value
		} else {
			return nil, fmt.Errorf("Unexpected Key %q", key)
		}
	}

	if partitionKey == "" {
		return nil, fmt.Errorf("Expected a PartitionKey but didn't get one")
	}
	if rowKey == "" {
		return nil, fmt.Errorf("Expected a RowKey but didn't get one")
	}

	return &ResourceID{
		AccountName:  *accountName,
		TableName:    tableName,
		PartitionKey: partitionKey,
		RowKey:       rowKey,
	}, nil
}
