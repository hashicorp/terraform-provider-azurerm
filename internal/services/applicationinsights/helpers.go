// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package applicationinsights

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandApplicationInsightsAPIKeyLinkedProperties(v *pluginsdk.Set, appInsightsId string) *[]string {
	if v == nil {
		return &[]string{}
	}

	result := make([]string, v.Len())
	for i, prop := range v.List() {
		result[i] = fmt.Sprintf("%s/%s", appInsightsId, prop)
	}
	return &result
}

func flattenApplicationInsightsAPIKeyLinkedProperties(props *[]string) *[]string {
	if props == nil {
		return &[]string{}
	}

	result := make([]string, len(*props))
	for i, prop := range *props {
		elems := strings.Split(prop, "/")
		result[i] = elems[len(elems)-1]
	}
	return &result
}
