// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventgrid_partner_namespace_channel_topic

type EventGridPartnerNamespaceChannelTopicResourceModel struct {
	Name                         string                       `tfschema:"name"`
	PartnerNamespaceId           string                       `tfschema:"partner_namespace_id"`
	PartnerTopic                 []ChannelPartnerTopicModel   `tfschema:"partner_topic"`
	ActivationMessage            string                       `tfschema:"activation_message"`
	EventTypes                   []ChannelTopicEventTypeModel `tfschema:"event_type"`
	ExpirationTimeIfNotActivated string                       `tfschema:"expiration_time_if_not_activated"`
	ReadinessState               string                       `tfschema:"readiness_state"`
}

type ChannelPartnerTopicModel struct {
	Name              string `tfschema:"name"`
	ResourceGroupName string `tfschema:"resource_group_name"`
	SubscriptionId    string `tfschema:"subscription_id"`
	Source            string `tfschema:"source"`
}

type ChannelTopicEventTypeModel struct {
	Name             string `tfschema:"name"`
	DataSchemaUrl    string `tfschema:"data_schema_url"`
	Description      string `tfschema:"description"`
	DisplayName      string `tfschema:"display_name"`
	DocumentationUrl string `tfschema:"documentation_url"`
}

func (r Resource) ModelObject() any {
	return &EventGridPartnerNamespaceChannelTopicResourceModel{}
}
