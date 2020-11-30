package parse

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/table/tables"
)

// TODO: tests for this

type StorageTableDataPlaneId struct {
	AccountName  string
	DomainSuffix string
	Name         string
}

// only present to comply with the interface
func (id StorageTableDataPlaneId) ID(_ string) string {
	return fmt.Sprintf("https://%s.table.%s/Tables('%s')", id.AccountName, id.DomainSuffix, id.Name)
}

func NewStorageTableDataPlaneId(accountName, domainSuffix, name string) StorageTableDataPlaneId {
	return StorageTableDataPlaneId{
		AccountName:  accountName,
		DomainSuffix: domainSuffix,
		Name:         name,
	}
}

func StorageTableDataPlaneID(input string) (*StorageTableDataPlaneId, error) {
	parsed, err := tables.ParseResourceID(input)
	if err != nil {
		return nil, err
	}

	uri, err := url.Parse(input)
	if err != nil {
		return nil, err
	}

	host := uri.Host
	hostSegments := strings.Split(host, ".")
	if len(hostSegments) == 0 {
		return nil, fmt.Errorf("expected multiple host segments but got 0")
	}
	domainNameSuffix := strings.TrimPrefix(host, fmt.Sprintf("%s.table.", hostSegments[0]))

	return &StorageTableDataPlaneId{
		AccountName:  parsed.AccountName,
		DomainSuffix: domainNameSuffix,
		Name:         parsed.TableName,
	}, nil
}
