// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-04-01/redisenterprise"

func PossibleValuesForAofFrequency() []string {
	res := make([]string, 0, len(redisenterprise.PossibleValuesForAofFrequency())-1)
	for _, freq := range redisenterprise.PossibleValuesForAofFrequency() {
		// `always` has been deprecated but not yet marked on the OpenAPI spec. It is no longer listed in the docs / portal:
		// https://learn.microsoft.com/azure/redis/how-to-persistence
		if freq != string(redisenterprise.AofFrequencyAlways) {
			res = append(res, freq)
		}
	}
	return res
}
