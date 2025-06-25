package manageddevopspools

import (
	"regexp"

	"github.com/hashicorp/go-azure-sdk/resource-manager/devopsinfrastructure/2025-01-21/pools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ResourcePredictionsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"time_zone": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
				"days_data": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsJSON,
				},
			},
		},
	}
}

func ResourcePredictionsProfileSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"kind": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						"Manual",
						"Automatic",
					}, false),
				},
				"prediction_preference": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringInSlice(pools.PossibleValuesForPredictionPreference(), false),
				},
			},
		},
	}
}

func ImageSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"aliases": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},
				"buffer": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringMatch(
						regexp.MustCompile(`^(?:\*|[1-9][0-9]?|100)$`),
						`Buffer must be "*" or value between 1 and 100.`,
					),
				},
				"resource_id": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ExactlyOneOf: []string{"resource_id", "well_known_image_name"},
				},
				"well_known_image_name": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ExactlyOneOf: []string{"resource_id", "well_known_image_name"},
				},
			},
		},
	}
}

func OsProfileSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"logon_type": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringInSlice(pools.PossibleValuesForLogonType(), false),
				},
				"secrets_management": SecretsManagementSettingsSchema(),
			},
		},
	}
}

func SecretsManagementSettingsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"certificate_store_location": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
				"certificate_store_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringInSlice(pools.PossibleValuesForCertificateStoreNameOption(), false),
				},
				"key_export_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},
				"observed_certificates": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},
			},
		},
	}
}

func StorageProfileSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		MaxItems: 1,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"data_disk": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"caching": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringInSlice(pools.PossibleValuesForCachingType(), false),
							},
							"disk_size_gb": {
								Type:     pluginsdk.TypeInt,
								Optional: true,
							},
							"drive_letter": {
								Type:     pluginsdk.TypeString,
								Optional: true,
							},
							"storage_account_type": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringInSlice(pools.PossibleValuesForStorageAccountType(), false),
							},
						},
					},
				},
				"os_disk_storage_account_type": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringInSlice(pools.PossibleValuesForOsDiskStorageAccountType(), false),
				},
			},
		},
	}
}
