// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package legacy

func expandZones(v []interface{}) *[]string {
	zones := make([]string, 0)
	for _, zone := range v {
		zones = append(zones, zone.(string))
	}
	if len(zones) > 0 {
		return &zones
	} else {
		return nil
	}
}
