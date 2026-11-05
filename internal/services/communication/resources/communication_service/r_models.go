// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package communication_service

type CommunicationServiceResourceModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	DataLocation      string            `tfschema:"data_location"`
	Tags              map[string]string `tfschema:"tags"`

	PrimaryConnectionString   string `tfschema:"primary_connection_string"`
	SecondaryConnectionString string `tfschema:"secondary_connection_string"`
	PrimaryKey                string `tfschema:"primary_key"`
	SecondaryKey              string `tfschema:"secondary_key"`
	HostName                  string `tfschema:"hostname"`
}

func (Resource) ModelObject() any {
	return &CommunicationServiceResourceModel{}
}
