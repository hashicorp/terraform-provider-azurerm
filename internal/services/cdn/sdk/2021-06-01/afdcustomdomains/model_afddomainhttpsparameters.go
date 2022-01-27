package afdcustomdomains

type AFDDomainHttpsParameters struct {
	CertificateType   AfdCertificateType    `json:"certificateType"`
	MinimumTlsVersion *AfdMinimumTlsVersion `json:"minimumTlsVersion,omitempty"`
	Secret            *ResourceReference    `json:"secret,omitempty"`
}
