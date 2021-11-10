package managedcluster

type ClientCertificate struct {
	CommonName       *string `json:"commonName,omitempty"`
	IsAdmin          bool    `json:"isAdmin"`
	IssuerThumbprint *string `json:"issuerThumbprint,omitempty"`
	Thumbprint       *string `json:"thumbprint,omitempty"`
}
