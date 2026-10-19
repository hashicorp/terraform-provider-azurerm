// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package email_communication_service

type EmailCommunicationServiceResourceModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	DataLocation      string            `tfschema:"data_location"`
	Tags              map[string]string `tfschema:"tags"`
}

func (Resource) ModelObject() any {
	return &EmailCommunicationServiceResourceModel{}
}
