// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryapplications"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryapplicationversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type GalleryApplicationVersionResource struct{}

var (
	_ sdk.ResourceWithUpdate        = GalleryApplicationVersionResource{}
	_ sdk.ResourceWithCustomizeDiff = GalleryApplicationVersionResource{}
)

type GalleryApplicationVersionModel struct {
	Name                 string            `tfschema:"name"`
	GalleryApplicationId string            `tfschema:"gallery_application_id"`
	Location             string            `tfschema:"location"`
	ConfigFile           string            `tfschema:"config_file"`
	EnableHealthCheck    bool              `tfschema:"enable_health_check"`
	EndOfLifeDate        string            `tfschema:"end_of_life_date"`
	ExcludeFromLatest    bool              `tfschema:"exclude_from_latest"`
	ManageAction         []ManageAction    `tfschema:"manage_action"`
	PackageFile          string            `tfschema:"package_file"`
	Source               []Source          `tfschema:"source"`
	TargetRegion         []TargetRegion    `tfschema:"target_region"`
	Tags                 map[string]string `tfschema:"tags"`
}

type Source struct {
	MediaLink                string `tfschema:"media_link"`
	DefaultConfigurationLink string `tfschema:"default_configuration_link"`
}

type ManageAction struct {
	Install string `tfschema:"install"`
	Remove  string `tfschema:"remove"`
	Update  string `tfschema:"update"`
}

type TargetRegion struct {
	Name                 string `tfschema:"name"`
	RegionalReplicaCount int64  `tfschema:"regional_replica_count"`
	ExcludeFromLatest    bool   `tfschema:"exclude_from_latest"`
	StorageAccountType   string `tfschema:"storage_account_type"`
}

func (r GalleryApplicationVersionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.GalleryApplicationVersionName,
		},

		"gallery_application_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: galleryapplications.ValidateApplicationID,
		},

		"location": commonschema.Location(),

		"config_file": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"enable_health_check": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"end_of_life_date": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			DiffSuppressFunc: suppress.RFC3339Time,
			ValidateFunc:     validation.IsRFC3339Time,
		},

		"exclude_from_latest": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"manage_action": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"install": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringLenBetween(1, 4096),
					},

					"remove": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringLenBetween(1, 4096),
					},

					"update": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringLenBetween(1, 4096),
					},
				},
			},
		},

		"package_file": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"source": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"media_link": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
					},

					"default_configuration_link": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
					},
				},
			},
		},

		"target_region": {
			// This needs to be a `TypeList` due to the `StateFunc` on the nested property `name`
			// See: https://github.com/hashicorp/terraform-plugin-sdk/issues/160
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": commonschema.LocationWithoutForceNew(),

					"regional_replica_count": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(1, 10),
					},

					"exclude_from_latest": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"storage_account_type": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(galleryapplicationversions.PossibleValuesForStorageAccountType(), false),
						Default:      string(galleryapplicationversions.StorageAccountTypeStandardLRS),
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (r GalleryApplicationVersionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r GalleryApplicationVersionResource) ResourceType() string {
	return "azurerm_gallery_application_version"
}

func (r GalleryApplicationVersionResource) ModelObject() interface{} {
	return &GalleryApplicationVersionModel{}
}

func (r GalleryApplicationVersionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return galleryapplicationversions.ValidateApplicationVersionID
}

func (r GalleryApplicationVersionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.GalleryApplicationVersionsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state GalleryApplicationVersionModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			galleryApplicationId, err := galleryapplications.ParseApplicationID(state.GalleryApplicationId)
			if err != nil {
				return err
			}

			id := galleryapplicationversions.NewApplicationVersionID(subscriptionId, galleryApplicationId.ResourceGroupName, galleryApplicationId.GalleryName, galleryApplicationId.ApplicationName, state.Name)
			existing, err := client.Get(ctx, id, galleryapplicationversions.DefaultGetOperationOptions())
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for the presence of existing %q: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := galleryapplicationversions.GalleryApplicationVersion{
				Location: location.Normalize(state.Location),
				Properties: &galleryapplicationversions.GalleryApplicationVersionProperties{
					PublishingProfile: galleryapplicationversions.GalleryApplicationVersionPublishingProfile{
						EnableHealthCheck: utils.Bool(state.EnableHealthCheck),
						ExcludeFromLatest: utils.Bool(state.ExcludeFromLatest),
						ManageActions:     expandGalleryApplicationVersionManageAction(state.ManageAction),
						Source:            expandGalleryApplicationVersionSource(state.Source),
						TargetRegions:     expandGalleryApplicationVersionTargetRegion(state.TargetRegion),
					},
					SafetyProfile: &galleryapplicationversions.GalleryArtifactSafetyProfileBase{
						AllowDeletionOfReplicatedLocations: pointer.To(true),
					},
				},
				Tags: pointer.To(state.Tags),
			}

			if state.ConfigFile != "" {
				if payload.Properties.PublishingProfile.Settings == nil {
					payload.Properties.PublishingProfile.Settings = &galleryapplicationversions.UserArtifactSettings{}
				}

				payload.Properties.PublishingProfile.Settings.ConfigFileName = &state.ConfigFile
			}

			if state.EndOfLifeDate != "" {
				endOfLifeDate, _ := time.Parse(time.RFC3339, state.EndOfLifeDate)
				payload.Properties.PublishingProfile.SetEndOfLifeDateAsTime(endOfLifeDate)
			}

			if state.PackageFile != "" {
				if payload.Properties.PublishingProfile.Settings == nil {
					payload.Properties.PublishingProfile.Settings = &galleryapplicationversions.UserArtifactSettings{}
				}

				payload.Properties.PublishingProfile.Settings.PackageFileName = &state.PackageFile
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r GalleryApplicationVersionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.GalleryApplicationVersionsClient
			id, err := galleryapplicationversions.ParseApplicationVersionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id, galleryapplicationversions.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					metadata.Logger.Infof("%s was not found - removing from state!", *id)
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := &GalleryApplicationVersionModel{
				Name:                 id.VersionName,
				GalleryApplicationId: galleryapplications.NewApplicationID(id.SubscriptionId, id.ResourceGroupName, id.GalleryName, id.ApplicationName).ID(),
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				if model.Tags != nil {
					state.Tags = *model.Tags
				}

				if props := model.Properties; props != nil {
					if props.PublishingProfile.EnableHealthCheck != nil {
						state.EnableHealthCheck = *props.PublishingProfile.EnableHealthCheck
					}

					if props.PublishingProfile.EndOfLifeDate != nil {
						d, err := props.PublishingProfile.GetEndOfLifeDateAsTime()
						if err != nil {
							return fmt.Errorf("parsing API response for `end_of_life_date`: %+v", err)
						}
						state.EndOfLifeDate = d.Format(time.RFC3339)
					}

					excludeFromLatest := false
					if props.PublishingProfile.ExcludeFromLatest != nil {
						excludeFromLatest = *props.PublishingProfile.ExcludeFromLatest
					}
					state.ExcludeFromLatest = excludeFromLatest

					state.ConfigFile = ""
					state.PackageFile = ""
					if props.PublishingProfile.Settings != nil {
						state.ConfigFile = pointer.From(props.PublishingProfile.Settings.ConfigFileName)
						state.PackageFile = pointer.From(props.PublishingProfile.Settings.PackageFileName)
					}

					state.ManageAction = flattenGalleryApplicationVersionManageAction(props.PublishingProfile.ManageActions)
					state.Source = flattenGalleryApplicationVersionSource(props.PublishingProfile.Source)
					state.TargetRegion = flattenGalleryApplicationVersionTargetRegion(props.PublishingProfile.TargetRegions)
				}
			}

			return metadata.Encode(state)
		},
		Timeout: 5 * time.Minute,
	}
}

func (r GalleryApplicationVersionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.GalleryApplicationVersionsClient

			id, err := galleryapplicationversions.ParseApplicationVersionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state GalleryApplicationVersionModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			payload := galleryapplicationversions.GalleryApplicationVersionUpdate{}

			if metadata.ResourceData.HasChanges("enable_health_check", "end_of_life_date", "exclude_from_latest", "manage_actions", "source", "target_region") {
				if payload.Properties == nil {
					payload.Properties = &galleryapplicationversions.GalleryApplicationVersionProperties{}
				}

				if metadata.ResourceData.HasChange("enable_health_check") {
					payload.Properties.PublishingProfile.EnableHealthCheck = utils.Bool(state.EnableHealthCheck)
				}

				if metadata.ResourceData.HasChange("end_of_life_date") {
					endOfLifeDate, _ := time.Parse(time.RFC3339, state.EndOfLifeDate)
					payload.Properties.PublishingProfile.SetEndOfLifeDateAsTime(endOfLifeDate)
				}

				if metadata.ResourceData.HasChange("exclude_from_latest") {
					payload.Properties.PublishingProfile.ExcludeFromLatest = utils.Bool(state.ExcludeFromLatest)
				}

				if metadata.ResourceData.HasChange("manage_actions") {
					payload.Properties.PublishingProfile.ManageActions = expandGalleryApplicationVersionManageAction(state.ManageAction)
				}

				if metadata.ResourceData.HasChange("source") {
					payload.Properties.PublishingProfile.Source = expandGalleryApplicationVersionSource(state.Source)
				}

				if metadata.ResourceData.HasChange("target_region") {
					payload.Properties.PublishingProfile.TargetRegions = expandGalleryApplicationVersionTargetRegion(state.TargetRegion)
				}
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(state.Tags)
			}

			if err := client.UpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r GalleryApplicationVersionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.GalleryApplicationVersionsClient
			id, err := galleryapplicationversions.ParseApplicationVersionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			metadata.Logger.Infof("Waiting for %s to be eventually deleted", *id)
			timeout, _ := ctx.Deadline()
			stateConf := &pluginsdk.StateChangeConf{
				Pending: []string{"Exists"},
				Target:  []string{"NotFound"},
				Refresh: func() (interface{}, string, error) {
					// Whilst the Gallery Application Version is deleted quickly, it appears it's not actually finished replicating at this time
					// so the deletion of the parent Gallery Application fails with "can not delete until nested resources are deleted"
					// ergo we need to poll on this for a bit, see https://github.com/Azure/azure-rest-api-specs/issues/19686
					res, err := client.Get(ctx, *id, galleryapplicationversions.DefaultGetOperationOptions())
					if err != nil {
						if response.WasNotFound(res.HttpResponse) {
							return "NotFound", "NotFound", nil
						}

						return nil, "", fmt.Errorf("polling to check if the %s has been deleted: %+v", *id, err)
					}

					return res, "Exists", nil
				},
				MinTimeout:                10 * time.Second,
				ContinuousTargetOccurence: 10,
				Timeout:                   time.Until(timeout),
			}

			if _, err := stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for %s to be deleted: %+v", *id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r GalleryApplicationVersionResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if oldVal, newVal := metadata.ResourceDiff.GetChange("end_of_life_date"); oldVal.(string) != "" && newVal.(string) == "" {
				if err := metadata.ResourceDiff.ForceNew("end_of_life_date"); err != nil {
					return err
				}
			}

			return nil
		},
	}
}

func expandGalleryApplicationVersionManageAction(input []ManageAction) *galleryapplicationversions.UserArtifactManage {
	if len(input) == 0 {
		return &galleryapplicationversions.UserArtifactManage{}
	}
	v := input[0]
	return &galleryapplicationversions.UserArtifactManage{
		Install: v.Install,
		Remove:  v.Remove,
		Update:  utils.String(v.Update),
	}
}

func flattenGalleryApplicationVersionManageAction(input *galleryapplicationversions.UserArtifactManage) []ManageAction {
	if input == nil {
		return nil
	}

	output := make([]ManageAction, 0)

	obj := ManageAction{
		Install: input.Install,
		Remove:  input.Remove,
	}
	if input.Update != nil {
		obj.Update = *input.Update
	}
	output = append(output, obj)

	return output
}

func expandGalleryApplicationVersionSource(input []Source) galleryapplicationversions.UserArtifactSource {
	if len(input) == 0 {
		return galleryapplicationversions.UserArtifactSource{}
	}
	v := input[0]
	return galleryapplicationversions.UserArtifactSource{
		MediaLink:                v.MediaLink,
		DefaultConfigurationLink: utils.String(v.DefaultConfigurationLink),
	}
}

func flattenGalleryApplicationVersionSource(input galleryapplicationversions.UserArtifactSource) []Source {
	out := Source{
		MediaLink: input.MediaLink,
	}
	if input.DefaultConfigurationLink != nil {
		out.DefaultConfigurationLink = *input.DefaultConfigurationLink
	}
	return []Source{
		out,
	}
}

func expandGalleryApplicationVersionTargetRegion(input []TargetRegion) *[]galleryapplicationversions.TargetRegion {
	results := make([]galleryapplicationversions.TargetRegion, 0)
	for _, item := range input {
		targetRegion := galleryapplicationversions.TargetRegion{
			Name:                 location.Normalize(item.Name),
			RegionalReplicaCount: pointer.To(item.RegionalReplicaCount),
			StorageAccountType:   pointer.To(galleryapplicationversions.StorageAccountType(item.StorageAccountType)),
		}

		if item.ExcludeFromLatest {
			targetRegion.ExcludeFromLatest = &item.ExcludeFromLatest
		}

		results = append(results, targetRegion)
	}

	return &results
}

func flattenGalleryApplicationVersionTargetRegion(input *[]galleryapplicationversions.TargetRegion) []TargetRegion {
	results := make([]TargetRegion, 0)

	for _, item := range *input {
		obj := TargetRegion{
			Name:              location.Normalize(item.Name),
			ExcludeFromLatest: false,
		}

		if item.ExcludeFromLatest != nil {
			obj.ExcludeFromLatest = *item.ExcludeFromLatest
		}

		if item.RegionalReplicaCount != nil {
			obj.RegionalReplicaCount = *item.RegionalReplicaCount
		}

		if item.StorageAccountType != nil {
			obj.StorageAccountType = string(*item.StorageAccountType)
		}

		results = append(results, obj)
	}

	return results
}
