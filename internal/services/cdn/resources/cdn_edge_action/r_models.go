// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_edge_action

type EdgeActionModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	SkuName           string            `tfschema:"sku_name"`
	Tags              map[string]string `tfschema:"tags"`
}

func (Resource) ModelObject() any {
	return &EdgeActionModel{}
}
