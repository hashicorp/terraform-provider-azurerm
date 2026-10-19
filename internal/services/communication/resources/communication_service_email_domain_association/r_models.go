// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package communication_service_email_domain_association

type EmailDomainAssociationResourceModel struct {
	CommunicationServiceId string `tfschema:"communication_service_id"`
	EMailServiceDomainId   string `tfschema:"email_service_domain_id"`
}

func (Resource) ModelObject() any {
	return &EmailDomainAssociationResourceModel{}
}
