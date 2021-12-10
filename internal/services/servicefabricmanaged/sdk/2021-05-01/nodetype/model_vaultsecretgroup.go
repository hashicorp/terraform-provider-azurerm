package nodetype

type VaultSecretGroup struct {
	SourceVault       SubResource        `json:"sourceVault"`
	VaultCertificates []VaultCertificate `json:"vaultCertificates"`
}
