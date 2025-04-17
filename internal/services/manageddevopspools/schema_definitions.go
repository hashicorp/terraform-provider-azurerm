package manageddevopspools

import (
	"regexp"

	"github.com/hashicorp/go-azure-sdk/resource-manager/devopsinfrastructure/2025-01-21/pools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func AgentProfileSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"kind": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string("Stateless"),
						string("Stateful"),
					}, false),
				},
				"grace_period_time_span": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"max_agent_lifetime": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"resource_predictions":         ResourcePredictionsSchema(),
				"resource_predictions_profile": ResourcePredictionsProfileSchema(),
			},
		},
	}
}

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

func FabricProfileSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"images": ImagesSchema(),
				"kind": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string("Vmss"),
					}, false),
				},
				"network_profile": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"subnet_id": {
								Type:     pluginsdk.TypeString,
								Required: true,
								ValidateFunc: validation.StringMatch(
									regexp.MustCompile(`^/subscriptions/[0-9a-fA-F-]{36}/resourceGroups/[-\w._()]+/providers/Microsoft\.Network/virtualNetworks/[-\w._()]+/subnets/[-\w._()]+$`),
									"Subnet ID must match the format '/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/virtualNetworks/{vnetName}/subnets/{subnetName}'.",
								),
							},
						},
					},
				},
				"os_profile": OsProfileSchema(),
				"sku": {
					Type:     pluginsdk.TypeList,
					Required: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"name": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
				"storage_profile": StorageProfileSchema(),
			},
		},
	}
}

func ImagesSchema() *pluginsdk.Schema {
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
				},
				"well_known_image_name": {
					Type:     pluginsdk.TypeString,
					Optional: true,
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
					Required:     true,
					ValidateFunc: validation.StringInSlice(pools.PossibleValuesForLogonType(), false),
				},
				"secrets_management_settings": SecretsManagementSettingsSchema(),
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
				"key_exportable": {
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
				"data_disks": {
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

func OrganizationProfileSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"kind": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string("AzureDevOps"),
					}, false),
				},
				"organizations": {
					Type:     pluginsdk.TypeList,
					Required: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"parallelism": {
								Type:     pluginsdk.TypeInt,
								Optional: true,
							},
							"projects": {
								Type:     pluginsdk.TypeSet,
								Optional: true,
								Elem: &pluginsdk.Schema{
									Type: pluginsdk.TypeString,
								},
							},
							"url": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.IsURLWithHTTPS,
							},
						},
					},
				},
				"permission_profile": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"groups": {
								Type:     pluginsdk.TypeSet,
								Optional: true,
								Elem: &pluginsdk.Schema{
									Type: pluginsdk.TypeString,
								},
							},
							"kind": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringInSlice(pools.PossibleValuesForAzureDevOpsPermissionType(), false),
							},
							"users": {
								Type:     pluginsdk.TypeSet,
								Optional: true,
								Elem: &pluginsdk.Schema{
									Type: pluginsdk.TypeString,
								},
							},
						},
					},
				},
			},
		},
	}
}

func AgentProfileComputedSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"kind": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"grace_period_time_span": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"max_agent_lifetime": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"resource_predictions":         ResourcePredictionsSchema(),
				"resource_predictions_profile": ResourcePredictionsProfileSchema(),
			},
		},
	}
}

func FabricProfileComputedSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"images": ImagesSchema(),
				"kind": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string("Vmss"),
					}, false),
				},
				"network_profile": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"subnet_id": {
								Type:     pluginsdk.TypeString,
								Required: true,
							},
						},
					},
				},
				"os_profile": OsProfileSchema(),
				"sku": {
					Type:     pluginsdk.TypeList,
					Required: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"name": {
								Type:     pluginsdk.TypeString,
								Required: true,
							},
						},
					},
				},
				"storage_profile": StorageProfileSchema(),
			},
		},
	}
}

func OrganizationProfileComputedSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"kind": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string("AzureDevOps"),
					}, false),
				},
				"organizations": {
					Type:     pluginsdk.TypeList,
					Required: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"parallelism": {
								Type:     pluginsdk.TypeInt,
								Optional: true,
							},
							"projects": {
								Type:     pluginsdk.TypeSet,
								Optional: true,
								Elem: &pluginsdk.Schema{
									Type: pluginsdk.TypeString,
								},
							},
							"url": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.IsURLWithHTTPS,
							},
						},
					},
				},
				"permission_profile": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"groups": {
								Type:     pluginsdk.TypeSet,
								Optional: true,
								Elem: &pluginsdk.Schema{
									Type: pluginsdk.TypeString,
								},
							},
							"kind": {
								Type:     pluginsdk.TypeString,
								Required: true,
							},
							"users": {
								Type:     pluginsdk.TypeSet,
								Optional: true,
								Elem: &pluginsdk.Schema{
									Type: pluginsdk.TypeString,
								},
							},
						},
					},
				},
			},
		},
	}
}
