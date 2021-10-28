package vaults

type Permissions struct {
	Certificates *[]CertificatePermissions `json:"certificates,omitempty"`
	Keys         *[]KeyPermissions         `json:"keys,omitempty"`
	Secrets      *[]SecretPermissions      `json:"secrets,omitempty"`
	Storage      *[]StoragePermissions     `json:"storage,omitempty"`
}
