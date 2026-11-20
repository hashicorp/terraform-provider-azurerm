// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package signalr_service_custom_domain

type CustomDomainSignalrServiceModel struct {
	Name                       string `tfschema:"name"`
	SignalRServiceId           string `tfschema:"signalr_service_id"`
	DomainName                 string `tfschema:"domain_name"`
	SignalrCustomCertificateId string `tfschema:"signalr_custom_certificate_id"`
}

func (r Resource) ModelObject() any {
	return &CustomDomainSignalrServiceModel{}
}
