package manageddevopspools

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

const (
	AgentProfileKindStateless = "Stateless"
	AgentProfileKindStateful  = "Stateful"
)

type ManagedDevOpsPoolResourceModel struct {
	AgentProfile               AgentProfileModel                          `tfschema:"agent_profile"`
	DevCenterProjectResourceId string                                     `tfschema:"dev_center_project_id"`
	FabricProfile              FabricProfileModel                         `tfschema:"fabric_profile"`
	Identity                   []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Location                   string                                     `tfschema:"location"`
	MaximumConcurrency         int64                                      `tfschema:"maximum_concurrency"`
	Name                       string                                     `tfschema:"name"`
	OrganizationProfile        OrganizationProfileModel                   `tfschema:"organization_profile"`
	Tags                       map[string]string                          `tfschema:"tags"`
	Type                       string                                     `tfschema:"type"`
}

type AgentProfileModel struct {
	GracePeriodTimeSpan        *string                         `tfschema:"grace_period_time_span"`
	Kind                       string                          `tfschema:"kind"`
	MaxAgentLifetime           *string                         `tfschema:"max_agent_lifetime"`
	PredictionPreference       *string                         `tfschema:"prediction_preference"`
	ResourcePredictions        *interface{}                    `tfschema:"resource_predictions"`
	ResourcePredictionsProfile ResourcePredictionsProfileModel `tfschema:"resource_predictions_profile"`
}

type ResourcePredictionsProfileModel struct {
	Kind                 string  `tfschema:"kind"`
	PredictionPreference *string `tfschema:"prediction_preference"`
}

type FabricProfileModel struct {
	Images         []ImageModel         `tfschema:"images"`
	Kind           string               `tfschema:"kind"`
	NetworkProfile *NetworkProfileModel `tfschema:"network_profile"`
	OsProfile      *OsProfileModel      `tfschema:"os_profile"`
	Sku            DevOpsAzureSkuModel  `tfschema:"sku"`
	StorageProfile *StorageProfileModel `tfschema:"storage_profile"`
}

type ImageModel struct {
	Aliases            *[]string `tfschema:"aliases"`
	Buffer             *string   `tfschema:"buffer"`
	ResourceId         *string   `tfschema:"resource_id"`
	WellKnownImageName *string   `tfschema:"well_known_image_name"`
}
type OsProfileModel struct {
	LogonType                 string                          `tfschema:"logon_type"`
	SecretsManagementSettings *SecretsManagementSettingsModel `tfschema:"secrets_management_settings"`
}

type SecretsManagementSettingsModel struct {
	CertificateStoreLocation *string  `tfschema:"certificate_store_location"`
	KeyExportable            bool     `tfschema:"key_exportable"`
	ObservedCertificates     []string `tfschema:"observed_certificates"`
}

type NetworkProfileModel struct {
	SubnetId string `tfschema:"subnet_id"`
}

type DevOpsAzureSkuModel struct {
	Name string `tfschema:"name"`
}

type StorageProfileModel struct {
	DataDisks                *[]DataDiskModel `tfschema:"data_disks"`
	OsDiskStorageAccountType string           `tfschema:"os_disk_storage_account_type"`
}

type DataDiskModel struct {
	Caching            string  `tfschema:"caching"`
	DiskSizeGiB        *int64  `tfschema:"disk_size"`
	DriveLetter        *string `tfschema:"drive_letter"`
	StorageAccountType string  `tfschema:"storage_account_type"`
}

type OrganizationProfileModel struct {
	Organizations     []OrganizationModel    `tfschema:"organizations"`
	PermissionProfile PermissionProfileModel `tfschema:"permission_profile"`
	Kind              string                 `tfschema:"kind"`
}

type OrganizationModel struct {
	Parallelism  int64     `tfschema:"parallelism"`
	Projects     []string  `tfschema:"projects"`
	Repositories *[]string `tfschema:"repositories"`
	Url          string    `tfschema:"url"`
}

type PermissionProfileModel struct {
	Groups []string `tfschema:"groups"`
	Kind   string   `tfschema:"kind"`
	Users  []string `tfschema:"users"`
}
