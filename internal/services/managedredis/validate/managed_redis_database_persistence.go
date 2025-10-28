// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-04-01/redisenterprise"

const (
	DbPersistenceMethodAOF = "AOF"
	DbPersistenceMethodRDB = "RDB"
)

func PossibleValuesForPersistenceMethod() []string {
	return []string{
		DbPersistenceMethodAOF,
		DbPersistenceMethodRDB,
	}
}

func PossibleValuesForPersistenceBackupFrequency() []string {
	res := make([]string, 0, len(redisenterprise.PossibleValuesForRdbFrequency())+len(redisenterprise.PossibleValuesForAofFrequency()))
	res = append(res, redisenterprise.PossibleValuesForRdbFrequency()...)
	res = append(res, redisenterprise.PossibleValuesForAofFrequency()...)
	return res
}
