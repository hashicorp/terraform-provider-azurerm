// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_account_table_properties

type AccountTablePropertiesModel struct {
	StorageAccountId string                                `json:"storage_account_id" tfschema:"storage_account_id"`
	CorsRule         []AccountTablePropertiesCorsRule      `tfschema:"cors_rule"`
	HourMetrics      []AccountTablePropertiesHourMetrics   `tfschema:"hour_metrics"`
	MinuteMetrics    []AccountTablePropertiesMinuteMetrics `tfschema:"minute_metrics"`
	Logging          []AccountTablePropertiesLogging       `tfschema:"logging"`
}

func (s Resource) ModelObject() any {
	return &AccountTablePropertiesModel{}
}
