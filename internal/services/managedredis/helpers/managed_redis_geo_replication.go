// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/databases"
)

func FlattenLinkedDatabases(dbs *[]databases.LinkedDatabase) []string {
	if dbs == nil {
		return []string{}
	}

	result := make([]string, 0, len(*dbs))
	for _, db := range *dbs {
		if db.Id != nil {
			result = append(result, pointer.From(db.Id))
		}
	}
	return result
}
