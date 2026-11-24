// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package relay_namespace

type RelayNamespaceResourceModel struct {
	Name                      string         `tfschema:"name"`
	Location                  string         `tfschema:"location"`
	ResourceGroupName         string         `tfschema:"resource_group_name"`
	SkuName                   string         `tfschema:"sku_name"`
	Tags                      map[string]any `tfschema:"tags"`
	MetricID                  string         `tfschema:"metric_id"`
	PrimaryConnectionString   string         `tfschema:"primary_connection_string"`
	SecondaryConnectionString string         `tfschema:"secondary_connection_string"`
	PrimaryKey                string         `tfschema:"primary_key"`
	SecondaryKey              string         `tfschema:"secondary_key"`
}

func (r Resource) ModelObject() any {
	return &RelayNamespaceResourceModel{}
}
