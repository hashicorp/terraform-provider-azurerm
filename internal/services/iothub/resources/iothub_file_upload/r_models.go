// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_file_upload

type IotHubFileUploadResourceModel struct {
	AuthenticationType   string `tfschema:"authentication_type"`
	ConnectionString     string `tfschema:"connection_string"`
	ContainerName        string `tfschema:"container_name"`
	DefaultTTL           string `tfschema:"default_ttl"`
	IdentityId           string `tfschema:"identity_id"`
	IotHubId             string `tfschema:"iothub_id"`
	LockDuration         string `tfschema:"lock_duration"`
	MaxDeliveryCount     int64  `tfschema:"max_delivery_count"`
	NotificationsEnabled bool   `tfschema:"notifications_enabled"`
	SasTTL               string `tfschema:"sas_ttl"`
}

func (r Resource) ModelObject() any {
	return &IotHubFileUploadResourceModel{}
}
