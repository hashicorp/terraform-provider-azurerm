package frontdoors

type CustomHttpsConfiguration struct {
	CertificateSource                    FrontDoorCertificateSource            `json:"certificateSource"`
	FrontDoorCertificateSourceParameters *FrontDoorCertificateSourceParameters `json:"frontDoorCertificateSourceParameters,omitempty"`
	KeyVaultCertificateSourceParameters  *KeyVaultCertificateSourceParameters  `json:"keyVaultCertificateSourceParameters,omitempty"`
	MinimumTlsVersion                    MinimumTLSVersion                     `json:"minimumTlsVersion"`
	ProtocolType                         FrontDoorTlsProtocolType              `json:"protocolType"`
}
