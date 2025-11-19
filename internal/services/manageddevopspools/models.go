package manageddevopspools

const (
	AgentProfileKindStateless = "Stateless"
	AgentProfileKindStateful  = "Stateful"
)

type StatefulAgentProfileModel struct {
	GracePeriodTimeSpan                 string                                     `tfschema:"grace_period_time_span"`
	MaxAgentLifetime                    string                                     `tfschema:"max_agent_lifetime"`
	ManualResourcePredictionsProfile    []ManualResourcePredictionsProfileModel    `tfschema:"manual_resource_predictions_profile"`
	AutomaticResourcePredictionsProfile []AutomaticResourcePredictionsProfileModel `tfschema:"automatic_resource_predictions_profile"`
}

type StatelessAgentProfileModel struct {
	ManualResourcePredictionsProfile    []ManualResourcePredictionsProfileModel    `tfschema:"manual_resource_predictions_profile"`
	AutomaticResourcePredictionsProfile []AutomaticResourcePredictionsProfileModel `tfschema:"automatic_resource_predictions_profile"`
}

type ResourcePredictionsSdkModel struct {
	TimeZone string             `tfschema:"time_zone"`
	DaysData []map[string]int64 `tfschema:"days_data"`
}

type ResourcePredictionsProfileModel struct {
	Kind                 string `tfschema:"kind"`
	PredictionPreference string `tfschema:"prediction_preference"`
}

type ManualResourcePredictionsProfileModel struct {
	TimeZone          string           `tfschema:"time_zone"`
	AllWeekSchedule   int64            `tfschema:"all_week_schedule"`
	SundaySchedule    map[string]int64 `tfschema:"sunday_schedule"`
	MondaySchedule    map[string]int64 `tfschema:"monday_schedule"`
	TuesdaySchedule   map[string]int64 `tfschema:"tuesday_schedule"`
	WednesdaySchedule map[string]int64 `tfschema:"wednesday_schedule"`
	ThursdaySchedule  map[string]int64 `tfschema:"thursday_schedule"`
	FridaySchedule    map[string]int64 `tfschema:"friday_schedule"`
	SaturdaySchedule  map[string]int64 `tfschema:"saturday_schedule"`
}

type AutomaticResourcePredictionsProfileModel struct {
	PredictionPreference string `tfschema:"prediction_preference"`
}

type VmssFabricProfileModel struct {
	Images         []ImageModel          `tfschema:"image"`
	SubnetId       string                `tfschema:"subnet_id"`
	OsProfile      []OsProfileModel      `tfschema:"os_profile"`
	SkuName        string                `tfschema:"sku_name"`
	StorageProfile []StorageProfileModel `tfschema:"storage_profile"`
}

type ImageModel struct {
	Aliases            []string `tfschema:"aliases"`
	Buffer             string   `tfschema:"buffer"`
	ResourceId         string   `tfschema:"resource_id"`
	WellKnownImageName string   `tfschema:"well_known_image_name"`
}

type OsProfileModel struct {
	LogonType                 string                           `tfschema:"logon_type"`
	SecretsManagementSettings []SecretsManagementSettingsModel `tfschema:"secrets_management"`
}

type SecretsManagementSettingsModel struct {
	CertificateStoreLocation string   `tfschema:"certificate_store_location"`
	CertificateStoreName     string   `tfschema:"certificate_store_name"`
	KeyExportable            bool     `tfschema:"key_export_enabled"`
	ObservedCertificates     []string `tfschema:"observed_certificates"`
}

type StorageProfileModel struct {
	DataDisks                []DataDiskModel `tfschema:"data_disk"`
	OsDiskStorageAccountType string          `tfschema:"os_disk_storage_account_type"`
}

type DataDiskModel struct {
	Caching            string `tfschema:"caching"`
	DiskSizeGB         int64  `tfschema:"disk_size_gb"`
	DriveLetter        string `tfschema:"drive_letter"`
	StorageAccountType string `tfschema:"storage_account_type"`
}

type AzureDevOpsOrganizationProfileModel struct {
	Organizations     []OrganizationModel                 `tfschema:"organization"`
	PermissionProfile []AzureDevOpsPermissionProfileModel `tfschema:"permission_profile"`
}

type OrganizationModel struct {
	Parallelism int64    `tfschema:"parallelism"`
	Projects    []string `tfschema:"projects"`
	Url         string   `tfschema:"url"`
}

type AzureDevOpsPermissionProfileModel struct {
	Kind                  string                                  `tfschema:"kind"`
	AdministratorAccounts []AzureDevOpsAdministratorAccountsModel `tfschema:"administrator_account"`
}

type AzureDevOpsAdministratorAccountsModel struct {
	Groups []string `tfschema:"groups"`
	Users  []string `tfschema:"users"`
}
