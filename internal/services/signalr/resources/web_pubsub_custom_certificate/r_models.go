// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_pubsub_custom_certificate

type CustomCertWebPubsubModel struct {
	Name               string `tfschema:"name"`
	WebPubsubId        string `tfschema:"web_pubsub_id"`
	CustomCertId       string `tfschema:"custom_certificate_id"`
	CertificateVersion string `tfschema:"certificate_version"`
}

func (r Resource) ModelObject() any {
	return &CustomCertWebPubsubModel{}
}
