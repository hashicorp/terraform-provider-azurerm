// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package active_directory_domain_service_trust

type DomainServiceTrustModel struct {
	Name                string   `tfschema:"name"`
	DomainServiceId     string   `tfschema:"domain_service_id"`
	TrustedDomainFqdn   string   `tfschema:"trusted_domain_fqdn"`
	TrustedDomainDnsIPs []string `tfschema:"trusted_domain_dns_ips"`
	Password            string   `tfschema:"password"`
}

func (r Resource) ModelObject() any {
	return &DomainServiceTrustModel{}
}
