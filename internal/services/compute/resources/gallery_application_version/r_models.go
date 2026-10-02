// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package gallery_application_version

type GalleryApplicationVersionModel struct {
	Name                 string            `tfschema:"name"`
	GalleryApplicationId string            `tfschema:"gallery_application_id"`
	Location             string            `tfschema:"location"`
	ConfigFile           string            `tfschema:"config_file"`
	EnableHealthCheck    bool              `tfschema:"enable_health_check"`
	EndOfLifeDate        string            `tfschema:"end_of_life_date"`
	ExcludeFromLatest    bool              `tfschema:"exclude_from_latest"`
	ManageAction         []ManageAction    `tfschema:"manage_action"`
	PackageFile          string            `tfschema:"package_file"`
	Source               []Source          `tfschema:"source"`
	TargetRegion         []TargetRegion    `tfschema:"target_region"`
	Tags                 map[string]string `tfschema:"tags"`
}

func (r Resource) ModelObject() any {
	return &GalleryApplicationVersionModel{}
}
