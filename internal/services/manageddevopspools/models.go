// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package manageddevopspools

const (
	AgentProfileKindStateless = "Stateless"
	AgentProfileKindStateful  = "Stateful"
)

type StatefulAgentModel struct {
	GracePeriodTimeSpan         string                             `tfschema:"grace_period_time_span"`
	MaxAgentLifetime            string                             `tfschema:"max_agent_lifetime"`
	ManualResourcePrediction    []ManualResourcePredictionModel    `tfschema:"manual_resource_prediction"`
	AutomaticResourcePrediction []AutomaticResourcePredictionModel `tfschema:"automatic_resource_prediction"`
}

type StatelessAgentModel struct {
	ManualResourcePrediction    []ManualResourcePredictionModel    `tfschema:"manual_resource_prediction"`
	AutomaticResourcePrediction []AutomaticResourcePredictionModel `tfschema:"automatic_resource_prediction"`
}

type ResourcePredictionsSdkModel struct {
	TimeZone string             `tfschema:"time_zone"`
	DaysData []map[string]int64 `tfschema:"days_data"`
}

type ResourcePredictionsModel struct {
	Kind                 string `tfschema:"kind"`
	PredictionPreference string `tfschema:"prediction_preference"`
}

type ManualResourcePredictionModel struct {
	TimeZoneName      string             `tfschema:"time_zone_name"`
	AllWeekSchedule   int64              `tfschema:"all_week_schedule"`
	SundaySchedule    []DayScheduleModel `tfschema:"sunday_schedule"`
	MondaySchedule    []DayScheduleModel `tfschema:"monday_schedule"`
	TuesdaySchedule   []DayScheduleModel `tfschema:"tuesday_schedule"`
	WednesdaySchedule []DayScheduleModel `tfschema:"wednesday_schedule"`
	ThursdaySchedule  []DayScheduleModel `tfschema:"thursday_schedule"`
	FridaySchedule    []DayScheduleModel `tfschema:"friday_schedule"`
	SaturdaySchedule  []DayScheduleModel `tfschema:"saturday_schedule"`
}

type DayScheduleModel struct {
	Time  string `tfschema:"time"`
	Count int64  `tfschema:"count"`
}

type AutomaticResourcePredictionModel struct {
	PredictionPreference string `tfschema:"prediction_preference"`
}

type VmssFabricModel struct {
	Images                   []ImageModel    `tfschema:"image"`
	SubnetId                 string          `tfschema:"subnet_id"`
	Security                 []SecurityModel `tfschema:"security"`
	SkuName                  string          `tfschema:"sku_name"`
	OsDiskStorageAccountType string          `tfschema:"os_disk_storage_account_type"`
	Storage                  []StorageModel  `tfschema:"storage"`
}

type ImageModel struct {
	Aliases            []string `tfschema:"aliases"`
	Buffer             string   `tfschema:"buffer"`
	Id                 string   `tfschema:"id"`
	WellKnownImageName string   `tfschema:"well_known_image_name"`
}

type SecurityModel struct {
	InteractiveLogonEnabled    bool                              `tfschema:"interactive_logon_enabled"`
	KeyVaultManagementSettings []KeyVaultManagementSettingsModel `tfschema:"key_vault_management"`
}

type KeyVaultManagementSettingsModel struct {
	CertificateStoreLocation string   `tfschema:"certificate_store_location"`
	CertificateStoreName     string   `tfschema:"certificate_store_name"`
	KeyExportable            bool     `tfschema:"key_export_enabled"`
	KeyVaultCertificateIds   []string `tfschema:"key_vault_certificate_ids"`
}

type StorageModel struct {
	Caching            string `tfschema:"caching"`
	DiskSizeInGB       int64  `tfschema:"disk_size_in_gb"`
	DriveLetter        string `tfschema:"drive_letter"`
	StorageAccountType string `tfschema:"storage_account_type"`
}

type AzureDevOpsOrganizationModel struct {
	Organizations []OrganizationModel          `tfschema:"organization"`
	Permission    []AzureDevOpsPermissionModel `tfschema:"permission"`
}

type OrganizationModel struct {
	Parallelism int64    `tfschema:"parallelism"`
	Projects    []string `tfschema:"projects"`
	Url         string   `tfschema:"url"`
}

type AzureDevOpsPermissionModel struct {
	Kind                  string                                  `tfschema:"kind"`
	AdministratorAccounts []AzureDevOpsAdministratorAccountsModel `tfschema:"administrator_account"`
}

type AzureDevOpsAdministratorAccountsModel struct {
	Groups []string `tfschema:"groups"`
	Users  []string `tfschema:"users"`
}
