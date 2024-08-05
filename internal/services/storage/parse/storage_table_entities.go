// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// TODO: tests for this
var _ resourceids.Id = StorageTableEntitiesId{}

// StorageTableEntitiesId is used by the plural data source azurerm_storage_table_entities
type StorageTableEntitiesId struct {
	AccountName  string
	DomainSuffix string
	TableName    string
	Filter       string
}

func (id StorageTableEntitiesId) String() string {
	components := []string{
		fmt.Sprintf("Account Name %q", id.AccountName),
		fmt.Sprintf("Domain Suffix %q", id.DomainSuffix),
		fmt.Sprintf("TableName %q", id.TableName),
		fmt.Sprintf("Filter %q", id.Filter),
	}
	return fmt.Sprintf("Storage Table %s", strings.Join(components, " / "))
}

func (id StorageTableEntitiesId) ID() string {
	return fmt.Sprintf("https://%s.table.%s/%s(%s)", id.AccountName, id.DomainSuffix, id.TableName, id.Filter)
}

func NewStorageTableEntitiesId(accountName, domainSuffix, tableName, filter string) StorageTableEntitiesId {
	s := utils.Base64EncodeIfNot(filter)
	sha := sha1.Sum([]byte(s))
	filterHash := hex.EncodeToString(sha[:])
	return StorageTableEntitiesId{
		AccountName:  accountName,
		DomainSuffix: domainSuffix,
		TableName:    tableName,
		Filter:       filterHash,
	}
}
