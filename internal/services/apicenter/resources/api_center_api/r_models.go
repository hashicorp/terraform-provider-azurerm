package api_center_api

type ApiCenterApiModel struct {
	ServiceId             string                  `tfschema:"api_center_service_id"`
	Name                  string                  `tfschema:"name"`
	Contacts              []Contact               `tfschema:"contact"`
	Description           string                  `tfschema:"description"`
	ExternalDocumentation []ExternalDocumentation `tfschema:"external_documentation"`
	Kind                  string                  `tfschema:"kind"`
	License               []License               `tfschema:"license"`
	LifecycleStage        string                  `tfschema:"lifecycle_stage"`
	Summary               string                  `tfschema:"summary"`
	Title                 string                  `tfschema:"display_name"`
	TermsOfServiceUrl     string                  `tfschema:"terms_of_service_url"`
}

type Contact struct {
	Email string `tfschema:"email"`
	Name  string `tfschema:"name"`
	Url   string `tfschema:"url"`
}

type ExternalDocumentation struct {
	Description string `tfschema:"description"`
	Title       string `tfschema:"title"`
	Url         string `tfschema:"url"`
}

type License struct {
	Identifier string `tfschema:"identifier"`
	Name       string `tfschema:"name"`
	Url        string `tfschema:"url"`
}
