package manageddevopspools

import (
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
				"vm_size": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"os_type": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"os_sku": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"image": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"image_version": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
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
				"image": ImageSchema(),
				"kind": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string("Vmss"),
					}, false),
				},
				"network_profile": {
					Type: pluginsdk.TypeList,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"subnet_id": {
								Type:     pluginsdk.TypeString,
								Optional: true,
							},
						},
					},
				},
				"os_profile": OsProfileSchema(),
				"sku": {
					Type: pluginsdk.TypeList,
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

func ImageSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type: pluginsdk.TypeList,
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
		Type: pluginsdk.TypeList,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"logon_type": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
				"secrets_management_settings": SecretsManagementSettingsSchema(),
			},
		},
	}
}

func SecretsManagementSettingsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type: pluginsdk.TypeList,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"certificate_store_location": {
					Type:     pluginsdk.TypeString,
					Optional: true,
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
		Type: pluginsdk.TypeList,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"data_disks": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"disk_size_gb": {
								Type:     pluginsdk.TypeInt,
								Optional: true,
							},
							"disk_type": {
								Type:     pluginsdk.TypeString,
								Optional: true,
							},
							"lun": {
								Type:     pluginsdk.TypeInt,
								Optional: true,
							},
						},
					},
				},
				"os_disk_storage_account_type": {
					Type:     pluginsdk.TypeString,
					Optional: true,
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
						string("GitHub"),
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
							"repositories": {
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
