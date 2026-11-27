// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_pubsub_custom_domain

type CustomDomainWebPubsubModel struct {
	Name                         string `tfschema:"name"`
	WebPubsubId                  string `tfschema:"web_pubsub_id"`
	DomainName                   string `tfschema:"domain_name"`
	WebPubsubCustomCertificateId string `tfschema:"web_pubsub_custom_certificate_id"`
}

func (r Resource) ModelObject() any {
	return &CustomDomainWebPubsubModel{}
}
