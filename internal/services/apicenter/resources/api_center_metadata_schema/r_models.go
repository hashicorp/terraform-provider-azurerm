// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_metadata_schema

type ApiCenterMetadataSchemaModel struct {
	ServiceId          string               `tfschema:"service_id"`
	Name               string               `tfschema:"name"`
	MetadataAssignment []MetadataAssignment `tfschema:"metadata_assignment"`
	Schema             string               `tfschema:"schema"`
}

type MetadataAssignment struct {
	Deprecated bool   `tfschema:"deprecated"`
	Entity     string `tfschema:"entity"`
	Required   bool   `tfschema:"required"`
}
