package api_center_api_definition

type ApiCenterApiDefinitionModel struct {
	ApiVersionId  string                                 `tfschema:"api_center_api_version_id"`
	Name          string                                 `tfschema:"name"`
	Description   string                                 `tfschema:"description"`
	Specification []ApiDefinitionPropertiesSpecification `tfschema:"specification"`
	Title         string                                 `tfschema:"display_name"`
}

type ApiDefinitionPropertiesSpecification struct {
	Name    string `tfschema:"name"`
	Version string `tfschema:"version"`
}
