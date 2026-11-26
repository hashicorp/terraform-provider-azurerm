// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package validate

func ElasticSanVolumeName(i any, k string) ([]string, []error) {
	return elasticSanResourceName(63)(i, k)
}
