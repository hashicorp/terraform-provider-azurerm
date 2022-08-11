package helper

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

//DataFileCount     *int64   `json:"dataFileCount,omitempty"`
//DataFileSize      *int64   `json:"dataFileSize,omitempty"`
//DataGrowth        *int64   `json:"dataGrowth,omitempty"`
//DefaultFilePath   *string  `json:"defaultFilePath,omitempty"`
//LogFileSize       *int64   `json:"logFileSize,omitempty"`
//LogGrowth         *int64   `json:"logGrowth,omitempty"`
//Luns              *[]int64 `json:"luns,omitempty"`
//PersistFolder     *bool    `json:"persistFolder,omitempty"`
//PersistFolderPath *string  `json:"persistFolderPath,omitempty"`

func TempDbSettings() *pluginsdk.Schema {
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
