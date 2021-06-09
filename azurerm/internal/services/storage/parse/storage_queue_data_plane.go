package parse

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/queue/queues"
)

// TODO: tests for this
var _ resourceid.Formatter = StorageQueueDataPlaneId{}

type StorageQueueDataPlaneId struct {
	AccountName  string
	DomainSuffix string
	Name         string
}

// only present to comply with the interface
func (id StorageQueueDataPlaneId) ID() string {
	return fmt.Sprintf("https://%s.queue.%s/%s", id.AccountName, id.DomainSuffix, id.Name)
}

func NewStorageQueueDataPlaneId(accountName, domainSuffix, name string) StorageQueueDataPlaneId {
	return StorageQueueDataPlaneId{
		AccountName:  accountName,
		DomainSuffix: domainSuffix,
		Name:         name,
	}
}

func StorageQueueDataPlaneID(id string) (*StorageQueueDataPlaneId, error) {
	parsed, err := queues.ParseResourceID(id)
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
	domainNameSuffix := strings.TrimPrefix(host, fmt.Sprintf("%s.queue.", hostSegments[0]))

	return &StorageQueueDataPlaneId{
		AccountName:  parsed.AccountName,
		DomainSuffix: domainNameSuffix,
		Name:         parsed.QueueName,
	}, nil
}
