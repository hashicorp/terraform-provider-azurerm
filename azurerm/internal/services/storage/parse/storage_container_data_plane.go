package parse

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/blob/containers"
)

// TODO: tests for this

var _ resourceid.Formatter = StorageContainerDataPlaneId{}

type StorageContainerDataPlaneId struct {
	AccountName  string
	DomainSuffix string
	Name         string
}

// only present to comply with the interface
func (id StorageContainerDataPlaneId) ID() string {
	return fmt.Sprintf("https://%s.blob.%s/%s", id.AccountName, id.DomainSuffix, id.Name)
}

func NewStorageContainerDataPlaneId(accountName, domainSuffix, name string) StorageContainerDataPlaneId {
	return StorageContainerDataPlaneId{
		AccountName:  accountName,
		DomainSuffix: domainSuffix,
		Name:         name,
	}
}

func StorageContainerDataPlaneID(id string) (*StorageContainerDataPlaneId, error) {
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
	domainNameSuffix := strings.TrimPrefix(host, fmt.Sprintf("%s.blob.", hostSegments[0]))

	return &StorageContainerDataPlaneId{
		AccountName:  parsed.AccountName,
		DomainSuffix: domainNameSuffix,
		Name:         parsed.ContainerName,
	}, nil
}
