package workspaces

type Encryption struct {
	KeyName     *string    `json:"KeyName,omitempty"`
	KeySource   *KeySource `json:"keySource,omitempty"`
	Keyvaulturi *string    `json:"keyvaulturi,omitempty"`
	Keyversion  *string    `json:"keyversion,omitempty"`
}
