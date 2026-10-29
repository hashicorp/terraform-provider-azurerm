// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_dataset_json

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataFactoryDatasetJSONCreateUpdate,
		Read:   resourceDataFactoryDatasetJSONRead,
		Update: resourceDataFactoryDatasetJSONCreateUpdate,
		Delete: resourceDataFactoryDatasetJSONDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.DataSetID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: dataFactoryDatasetJsonSchema(),

		CustomizeDiff: func(ctx context.Context, d *pluginsdk.ResourceDiff, i any) error {
			if _, hasCompression := d.GetOk("compression"); hasCompression {
				supportedCompressionTypes := []string{
					helpers.TypeBasicDatasetCompressionTypeBZip2,
					helpers.TypeBasicDatasetCompressionTypeDeflate,
					helpers.TypeBasicDatasetCompressionTypeGZip,
					helpers.TypeBasicDatasetCompressionTypeTar,
					helpers.TypeBasicDatasetCompressionTypeTarGZip,
					helpers.TypeBasicDatasetCompressionTypeZipDeflate,
					helpers.TypeBasicDatasetCompressionTypeSnappy,
					helpers.TypeBasicDatasetCompressionTypeLZ4,
				}

				supportedCompressionLevels := []string{
					"Fastest",
					"Optimal",
				}

				dynamicLevelEnabled := d.Get("compression.0.dynamic_level_enabled")
				dynamicTypeEnabled := d.Get("compression.0.dynamic_type_enabled")
				compressionLevel, hasCompressionLevel := d.GetOk("compression.0.level")
				compressionType, hasCompressionType := d.GetOk("compression.0.type")

				if hasCompressionType && !dynamicTypeEnabled.(bool) && !slices.Contains(supportedCompressionTypes, compressionType.(string)) {
					return fmt.Errorf("compression type must be a supported type when `dynamic_type_enabled` is false, supported types are: %v", supportedCompressionTypes)
				}

				if hasCompressionLevel && !dynamicLevelEnabled.(bool) && !slices.Contains(supportedCompressionLevels, compressionLevel.(string)) {
					return fmt.Errorf("compression level must be a supported level when `dynamic_level_enabled` is false, supported levels are: %v", supportedCompressionLevels)
				}
			}
			return nil
		},
	}
}
