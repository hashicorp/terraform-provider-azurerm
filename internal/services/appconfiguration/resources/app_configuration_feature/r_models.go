// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package app_configuration_feature

type FeatureResourceModel struct {
	ConfigurationStoreId string                       `tfschema:"configuration_store_id"`
	Description          string                       `tfschema:"description"`
	Enabled              bool                         `tfschema:"enabled"`
	Key                  string                       `tfschema:"key"`
	Name                 string                       `tfschema:"name"`
	Label                string                       `tfschema:"label"`
	Locked               bool                         `tfschema:"locked"`
	Tags                 map[string]any               `tfschema:"tags"`
	PercentageFilter     float64                      `tfschema:"percentage_filter_value"`
	TimewindowFilters    []TimewindowFilterParameters `tfschema:"timewindow_filter"`
	TargetingFilters     []TargetingFilterAudience    `tfschema:"targeting_filter"`
	CustomFilters        []CustomFilter               `tfschema:"custom_filter"`
}

func (k Resource) ModelObject() any {
	return &FeatureResourceModel{}
}
