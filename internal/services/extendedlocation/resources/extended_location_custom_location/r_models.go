// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package extended_location_custom_location

type ExtendedLocationCustomLocationResourceModel struct {
	Name                string      `tfschema:"name"`
	ResourceGroupName   string      `tfschema:"resource_group_name"`
	Location            string      `tfschema:"location"`
	Authentication      []AuthModel `tfschema:"authentication"`
	ClusterExtensionIds []string    `tfschema:"cluster_extension_ids"`
	DisplayName         string      `tfschema:"display_name"`
	HostResourceId      string      `tfschema:"host_resource_id"`
	HostType            string      `tfschema:"host_type"`
	Namespace           string      `tfschema:"namespace"`
}

type AuthModel struct {
	Type  string `tfschema:"type"`
	Value string `tfschema:"value"`
}

func (r Resource) ModelObject() any {
	return &ExtendedLocationCustomLocationResourceModel{}
}
