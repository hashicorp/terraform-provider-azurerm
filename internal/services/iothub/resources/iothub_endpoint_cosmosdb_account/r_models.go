// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_endpoint_cosmosdb_account

// azignore:AZG010 - package name does not match its path

type IotHubEndpointCosmosDBAccountModel struct {
	Name                 string `tfschema:"name"`
	ResourceGroupName    string `tfschema:"resource_group_name"`
	AuthenticationType   string `tfschema:"authentication_type"`
	ContainerName        string `tfschema:"container_name"`
	DatabaseName         string `tfschema:"database_name"`
	EndpointUri          string `tfschema:"endpoint_uri"`
	IdentityId           string `tfschema:"identity_id"`
	IothubId             string `tfschema:"iothub_id"`
	PartitionKeyName     string `tfschema:"partition_key_name"`
	PartitionKeyTemplate string `tfschema:"partition_key_template"`
	PrimaryKey           string `tfschema:"primary_key"`
	SecondaryKey         string `tfschema:"secondary_key"`
	SubscriptionId       string `tfschema:"subscription_id"`
}

func (r Resource) ModelObject() any {
	return &Resource{}
}
