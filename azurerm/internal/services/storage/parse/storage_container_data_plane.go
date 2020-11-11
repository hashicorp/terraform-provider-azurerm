package parse

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/blob/containers"
)

type StorageContainerDataPlaneId struct {
	AccountName  string
	DomainSuffix string
	Name         string
}

// only present to comply with the interface
func (id StorageContainerDataPlaneId) ID(_ string) string {
	return fmt.Sprintf("https://%s.%s/%s", id.AccountName, id.DomainSuffix, id.Name)
}

func NewStorageContainerDataPlaneId(accountName, domainSuffix, name string) StorageContainerDataPlaneId {
	return StorageContainerDataPlaneId{
		AccountName:  accountName,
		DomainSuffix: domainSuffix,
		Name:         name,
	}
}

// ParseResourceID parses the Resource ID and returns an object which can be used
// to interact with the Container Resource
func ParseStorageContainerDataPlaneID(id string) (*StorageContainerDataPlaneId, error) {
	parsed, err := containers.ParseResourceID(id)
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
	domainNameSuffix := strings.TrimPrefix(host, fmt.Sprintf("%s.", hostSegments[0]))

	return &StorageContainerDataPlaneId{
		AccountName:  parsed.AccountName,
		DomainSuffix: domainNameSuffix,
		Name:         parsed.ContainerName,
	}, nil
}
