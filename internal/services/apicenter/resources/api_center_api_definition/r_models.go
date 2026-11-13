// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_api_definition

type ApiCenterApiDefinitionModel struct {
	ApiVersionId  string                                 `tfschema:"api_version_id"`
	Name          string                                 `tfschema:"name"`
	Description   string                                 `tfschema:"description"`
	Specification []ApiDefinitionPropertiesSpecification `tfschema:"specification"`
	Title         string                                 `tfschema:"display_name"`
}

type ApiDefinitionPropertiesSpecification struct {
	Name    string `tfschema:"name"`
	Version string `tfschema:"version"`
	Url     string `tfschema:"url"`
	Content string `tfschema:"content"`
}
