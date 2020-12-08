package parse

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/file/shares"
)

// TODO: tests for this

type StorageShareDataPlaneId struct {
	AccountName  string
	DomainSuffix string
	Name         string
}

// only present to comply with the interface
func (id StorageShareDataPlaneId) ID(_ string) string {
	return fmt.Sprintf("https://%s.file.%s/%s", id.AccountName, id.DomainSuffix, id.Name)
}

func NewStorageShareDataPlaneId(accountName, domainSuffix, name string) StorageShareDataPlaneId {
	return StorageShareDataPlaneId{
		AccountName:  accountName,
		DomainSuffix: domainSuffix,
		Name:         name,
	}
}

func StorageShareDataPlaneID(id string) (*StorageShareDataPlaneId, error) {
	parsed, err := shares.ParseResourceID(id)
	if err != nil {
		return nil, err
	}

	uri, err := url.Parse(id)
	if err != nil {
		return nil, err
	}

	host := uri.Host
	hostSegments := strings.Split(host, ".")
	if len(hostSegments) == 0 {
		return nil, fmt.Errorf("expected multiple host segments but got 0")
	}
	domainNameSuffix := strings.TrimPrefix(host, fmt.Sprintf("%s.file.", hostSegments[0]))

	return &StorageShareDataPlaneId{
		AccountName:  parsed.AccountName,
		DomainSuffix: domainNameSuffix,
		Name:         parsed.ShareName,
	}, nil
}
