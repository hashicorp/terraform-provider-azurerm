// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package email_communication_service_domain

type EmailCommunicationServiceDomainResourceModel struct {
	Name                          string            `tfschema:"name"`
	EMailServiceID                string            `tfschema:"email_service_id"`
	DomainManagement              string            `tfschema:"domain_management"`
	UserEngagementTrackingEnabled bool              `tfschema:"user_engagement_tracking_enabled"`
	Tags                          map[string]string `tfschema:"tags"`

	FromSenderDomain     string                           `tfschema:"from_sender_domain"`
	MailFromSenderDomain string                           `tfschema:"mail_from_sender_domain"`
	VerificationRecords  []EmailDomainVerificationRecords `tfschema:"verification_records"`
}

func (Resource) ModelObject() any {
	return &EmailCommunicationServiceDomainResourceModel{}
}
