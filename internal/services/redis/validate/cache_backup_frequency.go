// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func CacheBackupFrequency(v interface{}, k string) ([]string, []error) {
	return validation.IntInSlice([]int{15, 30, 60, 360, 720, 1440})(v, k)
}
