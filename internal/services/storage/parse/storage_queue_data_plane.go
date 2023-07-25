// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/queue/queues"
)

var _ resourceids.Id = StorageQueueDataPlaneId{}

type StorageQueueDataPlaneId struct {
	AccountName  string
	DomainSuffix string
	Name         string
}

func (id StorageQueueDataPlaneId) String() string {
	components := []string{
		fmt.Sprintf("Account Name %q", id.AccountName),
		fmt.Sprintf("Domain Suffix %q", id.DomainSuffix),
		fmt.Sprintf("Name %q", id.Name),
	}
	return fmt.Sprintf("Storage Queue %s", strings.Join(components, " / "))
}

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
