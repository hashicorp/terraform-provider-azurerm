// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/blob/containers"
)

var _ resourceids.Id = StorageContainerDataPlaneId{}

type StorageContainerDataPlaneId struct {
	AccountName  string
	DomainSuffix string
	Name         string
}

func (id StorageContainerDataPlaneId) String() string {
	components := []string{
		fmt.Sprintf("Account Name %q", id.AccountName),
		fmt.Sprintf("Domain Suffix %q", id.DomainSuffix),
		fmt.Sprintf("Name %q", id.Name),
	}
	return fmt.Sprintf("Storage Container %s", strings.Join(components, " / "))
}

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
