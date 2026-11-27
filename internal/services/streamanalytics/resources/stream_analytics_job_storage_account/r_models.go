// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_job_storage_account

type JobStorageAccountModel struct {
	JobId              string `tfschema:"stream_analytics_job_id"`
	AuthenticationMode string `tfschema:"authentication_mode"`
	StorageAccountKey  string `tfschema:"storage_account_key"`
	StorageAccountName string `tfschema:"storage_account_name"`
}

func (r Resource) ModelObject() any {
	return &JobStorageAccountModel{}
}
