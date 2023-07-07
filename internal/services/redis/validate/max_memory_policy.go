// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func MaxMemoryPolicy(v interface{}, k string) (warnings []string, errors []error) {
	return validation.StringInSlice([]string{
		"allkeys-lfu",
		"allkeys-lru",
		"allkeys-random",
		"noeviction",
		"volatile-lru",
		"volatile-lfu",
		"volatile-random",
		"volatile-ttl",
	}, false)(v, k)
}
