package parse

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceid"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/blob/containers"
)

var _ resourceid.Formatter = StorageContainerDataPlaneId{}

type StorageContainerDataPlaneId struct {
	AccountName  string
	DomainSuffix string
	Name         string
}

func (id StorageContainerDataPlaneId) String() string {
	segments := []string{
		fmt.Sprintf("Account Name %q", id.AccountName),
		fmt.Sprintf("Domain Suffix %q", id.DomainSuffix),
		fmt.Sprintf("Name %q", id.Name),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Storage Container", segmentsStr)
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

	prefix := fmt.Sprintf("%s.blob.", hostSegments[0])
	if !strings.HasPrefix(host, prefix) {
		return nil, fmt.Errorf("expected blob host segment")
	}
	domainNameSuffix := strings.TrimPrefix(host, prefix)

	if len(parsed.ContainerName) == 0 {
		return nil, fmt.Errorf("expected container name")
	}

	return &StorageContainerDataPlaneId{
		AccountName:  parsed.AccountName,
		DomainSuffix: domainNameSuffix,
		Name:         parsed.ContainerName,
	}, nil
}
