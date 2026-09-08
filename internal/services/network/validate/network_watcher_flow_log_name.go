// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func NetworkWatcherFlowLogName(v interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[^\W_]$|^[^\W_][\w.\-]{0,78}[\w]$`), "the name can be up to 80 characters long. It must begin with a word character, and it must end with a word character or with '_'. The name may contain word characters or '.', '-', '_'")(v, k)
}
