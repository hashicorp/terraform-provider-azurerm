package frontdoors

type KeyVaultCertificateSourceParameters struct {
	SecretName    *string                                   `json:"secretName,omitempty"`
	SecretVersion *string                                   `json:"secretVersion,omitempty"`
	Vault         *KeyVaultCertificateSourceParametersVault `json:"vault,omitempty"`
}
