// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-04-01/redisenterprise"

const (
	DbPersistenceMethodAOF = "AppendOnlyFile"
	DbPersistenceMethodRDB = "RedisDatabase"
)

func PossibleValuesForPersistenceMethod() []string {
	return []string{
		DbPersistenceMethodAOF,
		DbPersistenceMethodRDB,
	}
}

func PossibleValuesForAofFrequency() []string {
	res := make([]string, 0, len(redisenterprise.PossibleValuesForAofFrequency())-1)
	for _, freq := range redisenterprise.PossibleValuesForAofFrequency() {
		// `always` is not documented as a valid value in MS learn doc / portal despite being available in the API
		// https://learn.microsoft.com/azure/redis/how-to-persistence
		if freq != string(redisenterprise.AofFrequencyAlways) {
			res = append(res, freq)
		}
	}
	return res
}

func PossibleValuesForPersistenceBackupFrequency() []string {
	res := make([]string, 0, len(redisenterprise.PossibleValuesForRdbFrequency())+len(PossibleValuesForAofFrequency()))
	res = append(res, redisenterprise.PossibleValuesForRdbFrequency()...)
	res = append(res, PossibleValuesForAofFrequency()...)
	return res
}
