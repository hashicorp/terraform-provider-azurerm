// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func MonitorDiagnosticSettingName(v interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringDoesNotMatch(regexp.MustCompile(`[<>*%&:\\?+\/]+`), `characters <, >, *, %, &, :, \, ?, +, / are not allowed`),
		validation.StringLenBetween(1, 260),
	)(v, k)
}
