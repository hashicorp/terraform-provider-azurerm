package endpoints

type UrlSigningKey struct {
	KeyId               string                       `json:"keyId"`
	KeySourceParameters KeyVaultSigningKeyParameters `json:"keySourceParameters"`
}
