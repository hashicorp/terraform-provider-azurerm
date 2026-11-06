package api_center_deployment

type ApiCenterDeploymentModel struct {
	ApiId         string             `tfschema:"api_center_api_id"`
	Name          string             `tfschema:"name"`
	DefinitionId  string             `tfschema:"definition_id"`
	Description   string             `tfschema:"description"`
	EnvironmentId string             `tfschema:"environment_id"`
	Server        []DeploymentServer `tfschema:"server"`
	State         string             `tfschema:"state"`
	Title         string             `tfschema:"display_name"`
}

type DeploymentServer struct {
	RuntimeUri []string `tfschema:"runtime_uri"`
}
