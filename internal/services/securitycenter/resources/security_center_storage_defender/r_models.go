// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package security_center_storage_defender

type StorageDefenderModel struct {
	StorageAccountId                 string `tfschema:"storage_account_id"`
	OverrideSubscriptionSettings     bool   `tfschema:"override_subscription_settings_enabled"`
	MalwareScanningOnUploadEnabled   bool   `tfschema:"malware_scanning_on_upload_enabled"`
	MalwareScanningOnUploadCapPerMon int64  `tfschema:"malware_scanning_on_upload_cap_gb_per_month"`
	SensitiveDataDiscoveryEnabled    bool   `tfschema:"sensitive_data_discovery_enabled"`
	ScanResultsEventGridTopicId      string `tfschema:"scan_results_event_grid_topic_id"`
}

func (s Resource) ModelObject() any {
	return &StorageDefenderModel{}
}
