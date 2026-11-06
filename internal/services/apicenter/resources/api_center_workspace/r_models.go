package api_center_workspace

type ApiCenterWorkspaceModel struct {
	Name               string `tfschema:"name"`
	ApiCenterServiceID string `tfschema:"api_center_service_id"`
	Description        string `tfschema:"description"`
	Title              string `tfschema:"title"`
}
