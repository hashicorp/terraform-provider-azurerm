// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cost_anomaly_alert

type AnomalyAlertModel struct {
	Name              string   `tfschema:"name"`
	DisplayName       string   `tfschema:"display_name"`
	SubscriptionId    string   `tfschema:"subscription_id"`
	NotificationEmail string   `tfschema:"notification_email"`
	EmailSubject      string   `tfschema:"email_subject"`
	EmailAddresses    []string `tfschema:"email_addresses"`
	Message           string   `tfschema:"message"`
}

func (Resource) ModelObject() any {
	return &AnomalyAlertModel{}
}
