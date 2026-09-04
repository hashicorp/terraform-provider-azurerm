// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

func ElasticSanSnapshotName(i interface{}, k string) ([]string, []error) {
	return elasticSanResourceName(63)(i, k)
}
