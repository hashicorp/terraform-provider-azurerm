// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package email_communication_service_domain_sender_username

type EmailCommunicationServiceDomainSenderUsernameModel struct {
	Name                 string `tfschema:"name"`
	EMailServiceDomainID string `tfschema:"email_service_domain_id"`
	DisplayName          string `tfschema:"display_name"`
}

func (Resource) ModelObject() any {
	return &EmailCommunicationServiceDomainSenderUsernameModel{}
}
