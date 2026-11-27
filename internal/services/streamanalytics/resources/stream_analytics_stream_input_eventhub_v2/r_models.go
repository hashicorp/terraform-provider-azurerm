// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_stream_input_eventhub_v2

type StreamInputEventHubV2ResourceModel struct {
	Name                      string `tfschema:"name"`
	StreamAnalyticsJobId      string `tfschema:"stream_analytics_job_id"`
	ServiceBusNamespace       string `tfschema:"servicebus_namespace"`
	EventHubName              string `tfschema:"eventhub_name"`
	EventHubConsumerGroupName string `tfschema:"eventhub_consumer_group_name"`

	SharedAccessPolicyKey  string          `tfschema:"shared_access_policy_key"`
	SharedAccessPolicyName string          `tfschema:"shared_access_policy_name"`
	PartitionKey           string          `tfschema:"partition_key"`
	AuthenticationMode     string          `tfschema:"authentication_mode"`
	Serialization          []Serialization `tfschema:"serialization"`
}

func (r Resource) ModelObject() any {
	return &StreamInputEventHubV2ResourceModel{}
}
