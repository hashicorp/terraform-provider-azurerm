// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package databricks_access_connector

type AccessConnectorResourceModel struct {
	Name          string            `tfschema:"name"`
	ResourceGroup string            `tfschema:"resource_group_name"`
	Location      string            `tfschema:"location"`
	Tags          map[string]string `tfschema:"tags"`
}

func (r Resource) ModelObject() any {
	return &AccessConnectorResourceModel{}
}
