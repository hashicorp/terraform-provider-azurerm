// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/file/shares"
)

// TODO: tests for this
var _ resourceids.Id = StorageShareDataPlaneId{}

type StorageShareDataPlaneId struct {
	AccountName  string
	DomainSuffix string
	Name         string
}

func (id StorageShareDataPlaneId) String() string {
	components := []string{
		fmt.Sprintf("Account Name %q", id.AccountName),
		fmt.Sprintf("Domain Suffix %q", id.DomainSuffix),
		fmt.Sprintf("Name %q", id.Name),
	}
	return fmt.Sprintf("Storage Share (%s)", strings.Join(components, " / "))
}

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
	domainNameSuffix := strings.TrimPrefix(host, fmt.Sprintf("%s.file.", hostSegments[0]))

	return &StorageShareDataPlaneId{
		AccountName:  parsed.AccountName,
		DomainSuffix: domainNameSuffix,
		Name:         parsed.ShareName,
	}, nil
}
