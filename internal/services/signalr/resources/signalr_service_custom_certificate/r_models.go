// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package signalr_service_custom_certificate

type CustomCertSignalrServiceResourceModel struct {
	Name               string `tfschema:"name"`
	SignalRServiceId   string `tfschema:"signalr_service_id"`
	CustomCertId       string `tfschema:"custom_certificate_id"`
	CertificateVersion string `tfschema:"certificate_version"`
}

func (r Resource) ModelObject() any {
	return &CustomCertSignalrServiceResourceModel{}
}
