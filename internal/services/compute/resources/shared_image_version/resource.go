// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package shared_image_version

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-07-03/galleryimageversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Properties struct {
	Output Output `tfschema:"output"`
}

type Source struct {
	MediaLink                string `tfschema:"media_link"`
	DefaultConfigurationLink string `tfschema:"default_configuration_link"`
}

type TargetRegion struct {
	Name                 string `tfschema:"name"`
	RegionalReplicaCount int64  `tfschema:"regional_replica_count"`
	ExcludeFromLatest    bool   `tfschema:"exclude_from_latest"`
	StorageAccountType   string `tfschema:"storage_account_type"`
}

type Output struct {
	AccessSAS string `tfschema:"accessSAS"`
}

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSharedImageVersionCreate,
		Read:   resourceSharedImageVersionRead,
		Update: resourceSharedImageVersionUpdate,
		Delete: resourceSharedImageVersionDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := galleryimageversions.ParseImageVersionID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: sharedImageVersionSchema(),

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.ForceNewIfChange("end_of_life_date", func(ctx context.Context, old, new, meta any) bool {
				return old.(string) != "" && new.(string) == ""
			}),
		),
	}
}
