// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package elastic_san_volume

func (r Resource) ModelObject() any {
	return &ElasticSANVolumeResourceModel{}
}

type ElasticSANVolumeResourceModel struct {
	CreateSource         []ElasticSANVolumeCreateSource `tfschema:"create_source"`
	Name                 string                         `tfschema:"name"`
	SizeInGiB            int64                          `tfschema:"size_in_gib"`
	TargetIqn            string                         `tfschema:"target_iqn"`
	TargetPortalHostname string                         `tfschema:"target_portal_hostname"`
	TargetPortalPort     int64                          `tfschema:"target_portal_port"`
	VolumeGroupId        string                         `tfschema:"volume_group_id"`
	VolumeId             string                         `tfschema:"volume_id"`
}
