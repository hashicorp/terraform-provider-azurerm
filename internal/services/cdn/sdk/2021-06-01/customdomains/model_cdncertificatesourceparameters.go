package customdomains

type CdnCertificateSourceParameters struct {
	CertificateType CertificateType `json:"certificateType"`
	TypeName        TypeName        `json:"typeName"`
}
