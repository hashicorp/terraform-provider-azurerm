package hybridkubernetes

type ListClusterUserCredentialProperties struct {
	AuthenticationMethod AuthenticationMethod `json:"authenticationMethod"`
	ClientProxy          bool                 `json:"clientProxy"`
}
