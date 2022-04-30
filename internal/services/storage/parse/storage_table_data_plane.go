package parse

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceid"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/table/tables"
)

var _ resourceid.Formatter = StorageTableDataPlaneId{}

type StorageTableDataPlaneId struct {
	AccountName  string
	DomainSuffix string
	Name         string
}

func (id StorageTableDataPlaneId) String() string {
	segments := []string{
		fmt.Sprintf("Account Name %q", id.AccountName),
		fmt.Sprintf("Domain Suffix %q", id.DomainSuffix),
		fmt.Sprintf("Name %q", id.Name),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Storage Table", segmentsStr)
}

// only present to comply with the interface
func (id StorageTableDataPlaneId) ID() string {
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

	prefix := fmt.Sprintf("%s.table.", hostSegments[0])
	if !strings.HasPrefix(host, prefix) {
		return nil, fmt.Errorf("expected table host segment")
	}
	domainNameSuffix := strings.TrimPrefix(host, prefix)

	if len(parsed.TableName) == 0 {
		return nil, fmt.Errorf("expected table name")
	}

	return &StorageTableDataPlaneId{
		AccountName:  parsed.AccountName,
		DomainSuffix: domainNameSuffix,
		Name:         parsed.TableName,
	}, nil
}
