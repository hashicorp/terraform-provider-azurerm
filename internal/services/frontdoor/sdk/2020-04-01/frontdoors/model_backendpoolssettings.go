package frontdoors

type BackendPoolsSettings struct {
	EnforceCertificateNameCheck *EnforceCertificateNameCheckEnabledState `json:"enforceCertificateNameCheck,omitempty"`
	SendRecvTimeoutSeconds      *int64                                   `json:"sendRecvTimeoutSeconds,omitempty"`
}
