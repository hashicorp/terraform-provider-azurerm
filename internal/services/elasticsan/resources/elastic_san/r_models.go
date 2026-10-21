// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package elastic_san

func (r Resource) ModelObject() any {
	return &ElasticSANResourceModel{}
}

type ElasticSANResourceModel struct {
	BaseSizeInTiB        int64                        `tfschema:"base_size_in_tib"`
	ExtendedSizeInTiB    int64                        `tfschema:"extended_size_in_tib"`
	Location             string                       `tfschema:"location"`
	Name                 string                       `tfschema:"name"`
	ResourceGroupName    string                       `tfschema:"resource_group_name"`
	Sku                  []ElasticSANResourceSkuModel `tfschema:"sku"`
	Tags                 map[string]any               `tfschema:"tags"`
	TotalIops            int64                        `tfschema:"total_iops"`
	TotalMBps            int64                        `tfschema:"total_mbps"`
	TotalSizeInTiB       int64                        `tfschema:"total_size_in_tib"`
	TotalVolumeSizeInGiB int64                        `tfschema:"total_volume_size_in_gib"`
	VolumeGroupCount     int64                        `tfschema:"volume_group_count"`
	Zones                []string                     `tfschema:"zones"`
}

type ElasticSANResourceSkuModel struct {
	Name string `tfschema:"name"`
	Tier string `tfschema:"tier"`
}
