// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_resource_anchor

type ResourceAnchorResourceModel struct {
	Location            string            `tfschema:"location"`
	Name                string            `tfschema:"name"`
	ResourceGroupName   string            `tfschema:"resource_group_name"`
	LinkedCompartmentID string            `tfschema:"linked_compartment_id"`
	Tags                map[string]string `tfschema:"tags"`
}

func (Resource) ModelObject() any {
	return &ResourceAnchorResourceModel{}
}
