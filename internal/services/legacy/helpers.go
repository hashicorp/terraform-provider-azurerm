// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package legacy

func expandZones(v []any) *[]string {
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
