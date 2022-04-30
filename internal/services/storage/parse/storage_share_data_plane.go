package parse

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceid"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/file/shares"
)

var _ resourceid.Formatter = StorageShareDataPlaneId{}

type StorageShareDataPlaneId struct {
	AccountName  string
	DomainSuffix string
	Name         string
}

func (id StorageShareDataPlaneId) String() string {
	segments := []string{
		fmt.Sprintf("Account Name %q", id.AccountName),
		fmt.Sprintf("Domain Suffix %q", id.DomainSuffix),
		fmt.Sprintf("Name %q", id.Name),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Storage Share", segmentsStr)
}

// only present to comply with the interface
func (id StorageShareDataPlaneId) ID() string {
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

	prefix := fmt.Sprintf("%s.file.", hostSegments[0])
	if !strings.HasPrefix(host, prefix) {
		return nil, fmt.Errorf("expected file host segment")
	}
	domainNameSuffix := strings.TrimPrefix(host, prefix)

	if len(parsed.ShareName) == 0 {
		return nil, fmt.Errorf("expected share name")
	}

	return &StorageShareDataPlaneId{
		AccountName:  parsed.AccountName,
		DomainSuffix: domainNameSuffix,
		Name:         parsed.ShareName,
	}, nil
}
