// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/table/tables"
)

// TODO: tests for this
var _ resourceids.Id = StorageTableDataPlaneId{}

type StorageTableDataPlaneId struct {
	AccountName  string
	DomainSuffix string
	Name         string
}

func (id StorageTableDataPlaneId) String() string {
	components := []string{
		fmt.Sprintf("Account Name %q", id.AccountName),
		fmt.Sprintf("Domain Suffix %q", id.DomainSuffix),
		fmt.Sprintf("Name %q", id.Name),
	}
	return fmt.Sprintf("Storage Table %s", strings.Join(components, " / "))
}

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
	domainNameSuffix := strings.TrimPrefix(host, fmt.Sprintf("%s.table.", hostSegments[0]))

	return &StorageTableDataPlaneId{
		AccountName:  parsed.AccountName,
		DomainSuffix: domainNameSuffix,
		Name:         parsed.TableName,
	}, nil
}
