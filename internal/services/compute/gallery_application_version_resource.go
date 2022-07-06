package compute

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2021-11-01/compute"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
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
	EnableHealthCheck    bool              `tfschema:"enable_health_check"`
	EndOfLifeDate        string            `tfschema:"end_of_life_date"`
	ExcludeFromLatest    bool              `tfschema:"exclude_from_latest"`
	ManageAction         []ManageAction    `tfschema:"manage_action"`
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
	RegionalReplicaCount int    `tfschema:"regional_replica_count"`
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
			ValidateFunc: validate.GalleryApplicationID,
		},

		"location": commonschema.Location(),

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
					"name": {
						Type:             pluginsdk.TypeString,
						Required:         true,
						StateFunc:        location.StateFunc,
						DiffSuppressFunc: location.DiffSuppressFunc,
					},

					"regional_replica_count": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(1, 10),
					},

					"storage_account_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(compute.StorageAccountTypePremiumLRS),
							string(compute.StorageAccountTypeStandardLRS),
							string(compute.StorageAccountTypeStandardZRS),
						}, false),
						Default: string(compute.StorageAccountTypeStandardLRS),
					},
				},
			},
		},

		"tags": tags.Schema(),
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
	return validate.GalleryApplicationVersionID
}

func (r GalleryApplicationVersionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var state GalleryApplicationVersionModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			client := metadata.Client.Compute.GalleryApplicationVersionsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			galleryApplicationId, err := parse.GalleryApplicationID(state.GalleryApplicationId)
			if err != nil {
				return err
			}

			id := parse.NewGalleryApplicationVersionID(subscriptionId, galleryApplicationId.ResourceGroup, galleryApplicationId.GalleryName, galleryApplicationId.ApplicationName, state.Name)
			existing, err := client.Get(ctx, id.ResourceGroup, id.GalleryName, id.ApplicationName, id.VersionName, "")
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for the presence of existing %q: %+v", id, err)
			}
			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			input := compute.GalleryApplicationVersion{
				Location: utils.String(location.Normalize(state.Location)),
				GalleryApplicationVersionProperties: &compute.GalleryApplicationVersionProperties{
					PublishingProfile: &compute.GalleryApplicationVersionPublishingProfile{
						EnableHealthCheck: utils.Bool(state.EnableHealthCheck),
						ExcludeFromLatest: utils.Bool(state.ExcludeFromLatest),
						ManageActions:     expandGalleryApplicationVersionManageAction(state.ManageAction),
						Source:            expandGalleryApplicationVersionSource(state.Source),
						TargetRegions:     expandGalleryApplicationVersionTargetRegion(state.TargetRegion),
					},
				},
				Tags: tags.FromTypedObject(state.Tags),
			}

			if state.EndOfLifeDate != "" {
				endOfLifeDate, _ := time.Parse(time.RFC3339, state.EndOfLifeDate)
				input.GalleryApplicationVersionProperties.PublishingProfile.EndOfLifeDate = &date.Time{
					Time: endOfLifeDate,
				}
			}

			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.GalleryName, id.ApplicationName, id.VersionName, input)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation of %s: %+v", id, err)
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
			id, err := parse.GalleryApplicationVersionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, id.ResourceGroup, id.GalleryName, id.ApplicationName, id.VersionName, "")
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					metadata.Logger.Infof("%q was not found - removing from state!", *id)
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			galleryApplicationId := parse.NewGalleryApplicationID(id.SubscriptionId, id.ResourceGroup, id.GalleryName, id.ApplicationName)

			state := &GalleryApplicationVersionModel{
				Name:                 id.VersionName,
				GalleryApplicationId: galleryApplicationId.ID(),
				Location:             location.NormalizeNilable(resp.Location),
				Tags:                 tags.ToTypedObject(resp.Tags),
			}

			if props := resp.GalleryApplicationVersionProperties; props != nil {
				if publishingProfile := props.PublishingProfile; publishingProfile != nil {
					if publishingProfile.EnableHealthCheck != nil {
						state.EnableHealthCheck = *publishingProfile.EnableHealthCheck
					}

					if publishingProfile.EndOfLifeDate != nil {
						state.EndOfLifeDate = publishingProfile.EndOfLifeDate.Format(time.RFC3339)
					}

					if publishingProfile.ExcludeFromLatest != nil {
						state.ExcludeFromLatest = *publishingProfile.ExcludeFromLatest
					}

					if publishingProfile.ManageActions != nil {
						state.ManageAction = flattenGalleryApplicationVersionManageAction(publishingProfile.ManageActions)
					}

					if publishingProfile.Source != nil {
						state.Source = flattenGalleryApplicationVersionSource(publishingProfile.Source)
					}

					if publishingProfile.TargetRegions != nil {
						state.TargetRegion = flattenGalleryApplicationVersionTargetRegion(publishingProfile.TargetRegions)
					}
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
			id, err := parse.GalleryApplicationVersionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state GalleryApplicationVersionModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			client := metadata.Client.Compute.GalleryApplicationVersionsClient
			existing, err := client.Get(ctx, id.ResourceGroup, id.GalleryName, id.ApplicationName, id.VersionName, "")
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.PublishingProfile == nil {
				existing.PublishingProfile = &compute.GalleryApplicationVersionPublishingProfile{}
			}

			if metadata.ResourceData.HasChange("enable_health_check") {
				existing.PublishingProfile.EnableHealthCheck = utils.Bool(state.EnableHealthCheck)
			}

			if metadata.ResourceData.HasChange("end_of_life_date") {
				endOfLifeDate, _ := time.Parse(time.RFC3339, state.EndOfLifeDate)
				existing.GalleryApplicationVersionProperties.PublishingProfile.EndOfLifeDate = &date.Time{
					Time: endOfLifeDate,
				}
			}

			if metadata.ResourceData.HasChange("exclude_from_latest") {
				existing.PublishingProfile.ExcludeFromLatest = utils.Bool(state.ExcludeFromLatest)
			}

			if metadata.ResourceData.HasChange("manage_actions") {
				existing.GalleryApplicationVersionProperties.PublishingProfile.ManageActions = expandGalleryApplicationVersionManageAction(state.ManageAction)
			}
			if metadata.ResourceData.HasChange("source") {
				existing.GalleryApplicationVersionProperties.PublishingProfile.Source = expandGalleryApplicationVersionSource(state.Source)
			}

			if metadata.ResourceData.HasChange("target_region") {
				existing.GalleryApplicationVersionProperties.PublishingProfile.TargetRegions = expandGalleryApplicationVersionTargetRegion(state.TargetRegion)
			}

			if metadata.ResourceData.HasChange("tags") {
				existing.Tags = tags.FromTypedObject(state.Tags)
			}

			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.GalleryName, id.ApplicationName, id.VersionName, existing)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for update of %s: %+v", id, err)
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
			id, err := parse.GalleryApplicationVersionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			future, err := client.Delete(ctx, id.ResourceGroup, id.GalleryName, id.ApplicationName, id.VersionName)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
			}

			metadata.Logger.Infof("Waiting for %s to be eventually deleted", *id)
			timeout, _ := ctx.Deadline()
			stateConf := &pluginsdk.StateChangeConf{
				Pending:                   []string{"Exists"},
				Target:                    []string{"NotFound"},
				Refresh:                   galleryApplicationVersionDeleteStateRefreshFunc(ctx, client, *id),
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
			rd := metadata.ResourceDiff

			if oldVal, newVal := rd.GetChange("end_of_life_date"); oldVal.(string) != "" && newVal.(string) == "" {
				if err := rd.ForceNew("end_of_life_date"); err != nil {
					return err
				}
			}

			return nil
		},
	}
}

func expandGalleryApplicationVersionManageAction(input []ManageAction) *compute.UserArtifactManage {
	if len(input) == 0 {
		return &compute.UserArtifactManage{}
	}
	v := input[0]
	return &compute.UserArtifactManage{
		Install: utils.String(v.Install),
		Remove:  utils.String(v.Remove),
		Update:  utils.String(v.Update),
	}
}

func flattenGalleryApplicationVersionManageAction(input *compute.UserArtifactManage) []ManageAction {
	if input == nil {
		return nil
	}

	obj := ManageAction{}

	if input.Install != nil {
		obj.Install = *input.Install
	}

	if input.Remove != nil {
		obj.Remove = *input.Remove
	}

	if input.Update != nil {
		obj.Update = *input.Update
	}

	return []ManageAction{obj}
}

func expandGalleryApplicationVersionSource(input []Source) *compute.UserArtifactSource {
	if len(input) == 0 {
		return &compute.UserArtifactSource{}
	}
	v := input[0]
	return &compute.UserArtifactSource{
		MediaLink:                utils.String(v.MediaLink),
		DefaultConfigurationLink: utils.String(v.DefaultConfigurationLink),
	}
}

func flattenGalleryApplicationVersionSource(input *compute.UserArtifactSource) []Source {
	if input == nil {
		return nil
	}

	obj := Source{}

	if input.MediaLink != nil {
		obj.MediaLink = *input.MediaLink
	}

	if input.DefaultConfigurationLink != nil {
		obj.DefaultConfigurationLink = *input.DefaultConfigurationLink
	}

	return []Source{obj}
}

func expandGalleryApplicationVersionTargetRegion(input []TargetRegion) *[]compute.TargetRegion {
	results := make([]compute.TargetRegion, 0)
	for _, item := range input {
		results = append(results, compute.TargetRegion{
			Name:                 utils.String(location.Normalize(item.Name)),
			RegionalReplicaCount: utils.Int32(int32(item.RegionalReplicaCount)),
			StorageAccountType:   compute.StorageAccountType(item.StorageAccountType),
		})
	}
	return &results
}

func flattenGalleryApplicationVersionTargetRegion(input *[]compute.TargetRegion) []TargetRegion {
	if input == nil {
		return nil
	}

	results := make([]TargetRegion, 0)

	for _, item := range *input {
		obj := TargetRegion{}

		if item.Name != nil {
			obj.Name = location.Normalize(*item.Name)
		}

		if item.RegionalReplicaCount != nil {
			obj.RegionalReplicaCount = int(*item.RegionalReplicaCount)
		}

		if item.StorageAccountType != "" {
			obj.StorageAccountType = string(item.StorageAccountType)
		}
		results = append(results, obj)
	}
	return results
}

func galleryApplicationVersionDeleteStateRefreshFunc(ctx context.Context, client *compute.GalleryApplicationVersionsClient, id parse.GalleryApplicationVersionId) pluginsdk.StateRefreshFunc {
	// Whilst the Gallery Application Version is deleted quickly, it appears it's not actually finished replicating at this time
	// so the deletion of the parent Gallery Application fails with "can not delete until nested resources are deleted"
	// ergo we need to poll on this for a bit, see https://github.com/Azure/azure-rest-api-specs/issues/19686
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id.ResourceGroup, id.GalleryName, id.ApplicationName, id.VersionName, "")
		if err != nil {
			if utils.ResponseWasNotFound(res.Response) {
				return "NotFound", "NotFound", nil
			}

			return nil, "", fmt.Errorf("failed to poll to check if the Gallery Application Version has been deleted: %+v", err)
		}

		return res, "Exists", nil
	}
}
