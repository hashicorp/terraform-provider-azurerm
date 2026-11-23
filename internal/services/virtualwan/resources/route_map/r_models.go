// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package route_map

type RouteMapModel struct {
	Name         string `tfschema:"name"`
	VirtualHubId string `tfschema:"virtual_hub_id"`
	Rules        []Rule `tfschema:"rule"`
}

func (r Resource) ModelObject() any {
	return &RouteMapModel{}
}
