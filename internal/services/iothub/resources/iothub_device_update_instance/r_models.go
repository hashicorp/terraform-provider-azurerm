// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_device_update_instance

type IotHubDeviceUpdateInstanceModel struct {
	Name                     string                          `tfschema:"name"`
	DeviceUpdateAccountId    string                          `tfschema:"device_update_account_id"`
	DiagnosticStorageAccount []DiagnosticStorageAccountModel `tfschema:"diagnostic_storage_account"`
	DiagnosticEnabled        bool                            `tfschema:"diagnostic_enabled"`
	IotHubId                 string                          `tfschema:"iothub_id"`
	Tags                     map[string]string               `tfschema:"tags"`
}

type DiagnosticStorageAccountModel struct {
	ConnectionString          string `tfschema:"connection_string"`
	ConnectionStringWOVersion int64  `tfschema:"connection_string_wo_version"`
	Id                        string `tfschema:"id"`
}

func (r Resource) ModelObject() any {
	return &IotHubDeviceUpdateInstanceModel{}
}
