// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventgrid_namespace_topic_id_association

type EventGridNamespaceTopicIdAssociationModel struct {
	EventGridNamespaceId      string `tfschema:"eventgrid_namespace_id"`
	EventGridNamespaceTopicId string `tfschema:"eventgrid_namespace_topic_id"`
}

func (r Resource) ModelObject() any {
	return &EventGridNamespaceTopicIdAssociationModel{}
}
