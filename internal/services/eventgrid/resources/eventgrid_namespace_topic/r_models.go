// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventgrid_namespace_topic

type EventGridNamespaceTopicResourceModel struct {
	Name                 string `tfschema:"name"`
	EventgridNamespaceId string `tfschema:"eventgrid_namespace_id"`
	EventRetentionInDays int64  `tfschema:"event_retention_in_days"`
}

func (r Resource) ModelObject() any {
	return &EventGridNamespaceTopicResourceModel{}
}
