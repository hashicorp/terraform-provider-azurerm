package workspaces

type WorkspaceProviderAuthorization struct {
	PrincipalId      string `json:"principalId"`
	RoleDefinitionId string `json:"roleDefinitionId"`
}
