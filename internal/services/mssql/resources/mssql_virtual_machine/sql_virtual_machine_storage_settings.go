// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_virtual_machine

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func StorageSettingSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"luns": {
					Description: "A list of Logical Unit Numbers for the disks.",
					Type:        pluginsdk.TypeList,
					Required:    true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeInt,
					},
				},
				"default_file_path": {
					Description:  "The SQL Server default path.",
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"use_storage_pool": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  true,
				},
			},
		},
	}
}

func SQLTempDBStorageSettingSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"data_file_count": {
					Description: "The SQL Server default file count.",
					Type:        pluginsdk.TypeInt,
					Optional:    true,
					Default:     8,
				},
				"data_file_size_mb": {
					Description: "The SQL Server default file size.",
					Type:        pluginsdk.TypeInt,
					Optional:    true,
					Default:     256,
				},
				"data_file_growth_in_mb": {
					Description: "The SQL Server default file size.",
					Type:        pluginsdk.TypeInt,
					Optional:    true,
					Default:     512,
				},
				"default_file_path": {
					Description:  "The SQL Server default path.",
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"log_file_size_mb": {
					Description: "The SQL Server default file size.",
					Type:        pluginsdk.TypeInt,
					Optional:    true,
					Default:     256,
				},
				"log_file_growth_mb": {
					Description: "The SQL Server default file size.",
					Type:        pluginsdk.TypeInt,
					Optional:    true,
					Default:     512,
				},
				"luns": {
					Description: "A list of Logical Unit Numbers for the disks.",
					Type:        pluginsdk.TypeList,
					Required:    true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeInt,
					},
				},
				"use_storage_pool": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
	}
}
