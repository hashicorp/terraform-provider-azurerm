// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package suppress

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ListOrder TODO remove for 4.0
func ListOrder(key, old, new string, d *schema.ResourceData) bool {
	// Taken from https://github.com/hashicorp/terraform-plugin-sdk/issues/477#issuecomment-1238807249
	// For a list, the key is path to the element, rather than the list.
	// E.g. "node_groups.2.ips.0"
	lastDotIndex := strings.LastIndex(key, ".")
	if lastDotIndex != -1 {
		key = key[:lastDotIndex]
	}

	oldData, newData := d.GetChange(key)
	if oldData == nil || newData == nil {
		return false
	}

	sOld := make([]string, len(oldData.([]interface{})))
	sNew := make([]string, len(newData.([]interface{})))

	for i, v := range oldData.([]interface{}) {
		sOld[i] = fmt.Sprint(v)
	}

	for i, v := range newData.([]interface{}) {
		sNew[i] = fmt.Sprint(v)
	}

	return stringSlicesAreEqual(sOld, sNew)
}

func stringSlicesAreEqual(a []string, b []string) bool {
	sort.Strings(a)
	sort.Strings(b)

	return reflect.DeepEqual(a, b)
}
