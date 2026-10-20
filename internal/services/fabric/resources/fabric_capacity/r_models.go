// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package fabric_capacity

type FabricCapacityResourceModel struct {
	Name                  string            `tfschema:"name"`
	ResourceGroupName     string            `tfschema:"resource_group_name"`
	AdministrationMembers []string          `tfschema:"administration_members"`
	Location              string            `tfschema:"location"`
	Sku                   []SkuModel        `tfschema:"sku"`
	Tags                  map[string]string `tfschema:"tags"`
}

type SkuModel struct {
	Name string `tfschema:"name"`
	Tier string `tfschema:"tier"`
}

func (r Resource) ModelObject() any {
	return &Resource{}
}
