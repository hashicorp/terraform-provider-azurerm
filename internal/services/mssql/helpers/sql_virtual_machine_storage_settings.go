// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

<<<<<<<< HEAD:internal/services/mssql/helpers/sql_virtual_machine_storage_settings.go
package helpers
========
package mssql_virtual_machine
>>>>>>>> bb31be9448 (re run mssql for conflicts resolution):internal/services/mssql/resources/mssql_virtual_machine/sql_virtual_machine_storage_settings.go

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
					Type:     pluginsdk.TypeList,
					Required: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeInt,
					},
				},
				"default_file_path": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
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
					Type:     pluginsdk.TypeInt,
					Optional: true,
					Default:  8,
				},
				"data_file_size_mb": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
					Default:  256,
				},
				"data_file_growth_in_mb": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
					Default:  512,
				},
				"default_file_path": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"log_file_size_mb": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
					Default:  256,
				},
				"log_file_growth_mb": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
					Default:  512,
				},
				"luns": {
					Type:     pluginsdk.TypeList,
					Required: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeInt,
					},
				},
			},
		},
	}
}
