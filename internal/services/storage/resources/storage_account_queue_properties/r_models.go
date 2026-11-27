// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_account_queue_properties

type AccountQueuePropertiesModel struct {
	StorageAccountId string                                `json:"storage_account_id" tfschema:"storage_account_id"`
	CorsRule         []AccountQueuePropertiesCorsRule      `tfschema:"cors_rule"`
	HourMetrics      []AccountQueuePropertiesHourMetrics   `tfschema:"hour_metrics"`
	MinuteMetrics    []AccountQueuePropertiesMinuteMetrics `tfschema:"minute_metrics"`
	Logging          []AccountQueuePropertiesLogging       `tfschema:"logging"`
}

func (s Resource) ModelObject() any {
	return &AccountQueuePropertiesModel{}
}
