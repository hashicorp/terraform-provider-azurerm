// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package desktopvirtualization

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2024-04-03/hostpool"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/appattachpackage"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/msiximage"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/method"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name virtual_desktop_app_attach_package -service-package-name desktopvirtualization -properties "name,resource_group_name" -known-values "subscription_id:data.Subscriptions.Primary"

type VirtualDesktopAppAttachPackageResource struct{}

var (
	_ sdk.Resource                  = VirtualDesktopAppAttachPackageResource{}
	_ sdk.ResourceWithUpdate        = VirtualDesktopAppAttachPackageResource{}
	_ sdk.ResourceWithIdentity      = VirtualDesktopAppAttachPackageResource{}
	_ sdk.ResourceWithCustomizeDiff = VirtualDesktopAppAttachPackageResource{}
)

func (r VirtualDesktopAppAttachPackageResource) Identity() resourceids.ResourceId {
	return &appattachpackage.AppAttachPackageId{}
}

type VirtualDesktopAppAttachPackageModel struct {
	Name                       string                        `tfschema:"name"`
	ResourceGroupName          string                        `tfschema:"resource_group_name"`
	Location                   string                        `tfschema:"location"`
	DisplayName                string                        `tfschema:"display_name"`
	HostPoolIds                []string                      `tfschema:"host_pool_ids"`
	MsixPackageName            string                        `tfschema:"msix_package_name"`
	StorageShareFileId         string                        `tfschema:"storage_share_file_id"`
	HealthCheckStatusOnFailure string                        `tfschema:"health_check_status_on_failure"`
	RegisterAtLogOnEnabled     bool                          `tfschema:"register_at_log_on_enabled"`
	StateEnabled               bool                          `tfschema:"state_enabled"`
	Tags                       map[string]string             `tfschema:"tags"`
	LastUpdated                string                        `tfschema:"last_updated"`
	PackageApplications        []MsixPackageApplicationModel `tfschema:"package_applications"`
	PackageFamilyName          string                        `tfschema:"package_family_name"`
	PackageName                string                        `tfschema:"package_name"`
	PackageRelativePath        string                        `tfschema:"package_relative_path"`
	Version                    string                        `tfschema:"version"`
}

type MsixPackageApplicationModel struct {
	AppId          string `tfschema:"app_id"`
	AppUserModelID string `tfschema:"app_user_model_id"`
	Description    string `tfschema:"description"`
	FriendlyName   string `tfschema:"friendly_name"`
	IconImageName  string `tfschema:"icon_image_name"`
	RawIcon        string `tfschema:"raw_icon"`
	RawPng         string `tfschema:"raw_png"`
}

func (r VirtualDesktopAppAttachPackageResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			// `ValidateFunc` coded according to portal
			ValidateFunc: validation.StringDoesNotContainAny("\\/+?&"),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"host_pool_ids": {
			// `TypeSet` because order is not guaranteed
			Type:     pluginsdk.TypeSet,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: hostpool.ValidateHostPoolID,
			},
		},

		"msix_package_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"storage_share_file_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},

		"health_check_status_on_failure": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      appattachpackage.FailHealthCheckOnStagingFailureNeedsAssistance,
			ValidateFunc: validation.StringInSlice(appattachpackage.PossibleValuesForFailHealthCheckOnStagingFailure(), false),
		},

		"register_at_log_on_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"state_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"tags": commonschema.Tags(),
	}
}

func (r VirtualDesktopAppAttachPackageResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"last_updated": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"package_applications": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"app_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"app_user_model_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"description": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"friendly_name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"icon_image_name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"raw_icon": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"raw_png": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"package_family_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"package_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"package_relative_path": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r VirtualDesktopAppAttachPackageResource) ModelObject() interface{} {
	return &VirtualDesktopAppAttachPackageModel{}
}

func (r VirtualDesktopAppAttachPackageResource) ResourceType() string {
	return "azurerm_virtual_desktop_app_attach_package"
}

func (r VirtualDesktopAppAttachPackageResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.AppAttachPackagesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model VirtualDesktopAppAttachPackageModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := appattachpackage.NewAppAttachPackageID(subscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			// Replicate portal behavior where MSIX image properties required for app attach package creation are retrieved with API
			msixImageProperties, err := getMsixImageProperties(ctx, metadata, model.HostPoolIds, model.StorageShareFileId, model.MsixPackageName)
			if err != nil {
				return fmt.Errorf("retrieving MSIX image properties: %+v", err)
			}

			param := appattachpackage.AppAttachPackage{
				Location: location.Normalize(model.Location),
				Properties: appattachpackage.AppAttachPackageProperties{
					Image:                           r.expandVirtualDesktopAppAttachPackageImage(model, msixImageProperties),
					HostPoolReferences:              pointer.To(model.HostPoolIds),
					FailHealthCheckOnStagingFailure: pointer.ToEnum[appattachpackage.FailHealthCheckOnStagingFailure](model.HealthCheckStatusOnFailure),
				},
			}

			if len(model.Tags) > 0 {
				param.Tags = pointer.To(model.Tags)
			}

			if _, err := client.CreateOrUpdate(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}

			return nil
		},
	}
}

func (r VirtualDesktopAppAttachPackageResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.AppAttachPackagesClient

			id, err := appattachpackage.ParseAppAttachPackageID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model VirtualDesktopAppAttachPackageModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			param := *existing.Model

			// Update `param.Properties.Image` if any of the child properties has changed
			if metadata.ResourceData.HasChanges("display_name", "msix_package_name", "storage_share_file_id", "register_at_log_on_enabled", "state_enabled") {
				msixImageProperties, err := getMsixImageProperties(ctx, metadata, model.HostPoolIds, model.StorageShareFileId, model.MsixPackageName)
				if err != nil {
					return fmt.Errorf("retrieving MSIX image properties: %+v", err)
				}

				param.Properties.Image = r.expandVirtualDesktopAppAttachPackageImage(model, msixImageProperties)
			}

			if metadata.ResourceData.HasChange("host_pool_ids") {
				param.Properties.HostPoolReferences = pointer.To(model.HostPoolIds)
			}

			if metadata.ResourceData.HasChange("health_check_status_on_failure") {
				param.Properties.FailHealthCheckOnStagingFailure = pointer.ToEnum[appattachpackage.FailHealthCheckOnStagingFailure](model.HealthCheckStatusOnFailure)
			}

			if metadata.ResourceData.HasChange("tags") {
				param.Tags = pointer.To(model.Tags)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, param); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r VirtualDesktopAppAttachPackageResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.AppAttachPackagesClient

			id, err := appattachpackage.ParseAppAttachPackageID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			return r.flatten(metadata, id, existing.Model)
		},
	}
}

func (r VirtualDesktopAppAttachPackageResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.AppAttachPackagesClient

			id, err := appattachpackage.ParseAppAttachPackageID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r VirtualDesktopAppAttachPackageResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return appattachpackage.ValidateAppAttachPackageID
}

func (r VirtualDesktopAppAttachPackageResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			oldHostPoolIds, newHostPoolIds := metadata.ResourceDiff.GetChange("host_pool_ids")

			// If number of `host_pool_ids` elements has changed, actual elements will be updated incorrectly to be different from Terraform configurations
			if oldHostPoolIds.(*pluginsdk.Set).Len() != newHostPoolIds.(*pluginsdk.Set).Len() {
				metadata.ResourceDiff.ForceNew("host_pool_ids")
			}

			return nil
		},
	}
}

func (r VirtualDesktopAppAttachPackageResource) flatten(metadata sdk.ResourceMetaData, id *appattachpackage.AppAttachPackageId, model *appattachpackage.AppAttachPackage) error {
	state := VirtualDesktopAppAttachPackageModel{
		Name:              id.AppAttachPackageName,
		ResourceGroupName: id.ResourceGroupName,
	}

	if model != nil {
		state.Location = location.Normalize(model.Location)

		props := model.Properties
		if image := props.Image; image != nil {
			state.DisplayName = pointer.From(image.DisplayName)
			state.MsixPackageName = pointer.From(image.PackageFullName)

			if image.ImagePath != nil {
				// Convert from UNC format (`ImagePath`) to URL format (`StorageShareFileId`)
				storageShareFileId := strings.ReplaceAll(pointer.From(image.ImagePath), "\\", "/")
				storageShareFileId = fmt.Sprintf("https:%s", storageShareFileId)
				state.StorageShareFileId = storageShareFileId
			}

			if image.IsRegularRegistration != nil {
				state.RegisterAtLogOnEnabled = pointer.From(image.IsRegularRegistration)
			}

			if image.IsActive != nil {
				state.StateEnabled = pointer.From(image.IsActive)
			}

			state.LastUpdated = pointer.From(image.LastUpdated)
			state.PackageApplications = r.flattenVirtualDesktopAppAttachPackageApplications(image.PackageApplications)
			state.PackageFamilyName = pointer.From(image.PackageFamilyName)
			state.PackageName = pointer.From(image.PackageName)
			state.PackageRelativePath = pointer.From(image.PackageRelativePath)
			state.Version = pointer.From(image.Version)
		}

		state.HostPoolIds = pointer.From(props.HostPoolReferences)
		for i, hostPoolId := range state.HostPoolIds {
			parsedHostPoolId, err := hostpool.ParseHostPoolIDInsensitively(hostPoolId)
			if err != nil {
				return fmt.Errorf("parsing host pool ID %s: %+v", hostPoolId, err)
			}
			state.HostPoolIds[i] = parsedHostPoolId.ID()
		}

		if props.FailHealthCheckOnStagingFailure != nil {
			state.HealthCheckStatusOnFailure = pointer.FromEnum(props.FailHealthCheckOnStagingFailure)
		}

		state.Tags = pointer.From(model.Tags)
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}
	return metadata.Encode(&state)
}

func (r VirtualDesktopAppAttachPackageResource) expandVirtualDesktopAppAttachPackageImage(model VirtualDesktopAppAttachPackageModel, msixImageProperties *msiximage.ExpandMsixImageProperties) *appattachpackage.AppAttachPackageInfoProperties {
	// Convert from URL format (`StorageShareFileId`) to UNC format (`ImagePath`)
	imagePath := strings.ReplaceAll(model.StorageShareFileId, "/", "\\")
	imagePath = strings.TrimPrefix(imagePath, "https:")

	return &appattachpackage.AppAttachPackageInfoProperties{
		DisplayName:           pointer.To(model.DisplayName),
		PackageFullName:       pointer.To(model.MsixPackageName),
		ImagePath:             pointer.To(imagePath),
		IsRegularRegistration: pointer.To(model.RegisterAtLogOnEnabled),
		IsActive:              pointer.To(model.StateEnabled),
		LastUpdated:           msixImageProperties.LastUpdated,
		PackageApplications:   r.expandVirtualDesktopAppAttachPackageApplications(msixImageProperties.PackageApplications),
		PackageFamilyName:     msixImageProperties.PackageFamilyName,
		PackageName:           msixImageProperties.PackageName,
		PackageRelativePath:   msixImageProperties.PackageRelativePath,
		Version:               msixImageProperties.Version,
	}
}

func (r VirtualDesktopAppAttachPackageResource) expandVirtualDesktopAppAttachPackageApplications(inputs *[]msiximage.MsixPackageApplications) *[]appattachpackage.MsixPackageApplications {
	outputs := make([]appattachpackage.MsixPackageApplications, 0)
	if inputs == nil {
		return pointer.To(outputs)
	}

	for _, input := range *inputs {
		outputs = append(outputs, appattachpackage.MsixPackageApplications{
			AppId:          input.AppId,
			AppUserModelID: input.AppUserModelID,
			Description:    input.Description,
			FriendlyName:   input.FriendlyName,
			IconImageName:  input.IconImageName,
			RawIcon:        input.RawIcon,
			RawPng:         input.RawPng,
		})
	}

	return pointer.To(outputs)
}

func (r VirtualDesktopAppAttachPackageResource) flattenVirtualDesktopAppAttachPackageApplications(inputs *[]appattachpackage.MsixPackageApplications) []MsixPackageApplicationModel {
	outputs := make([]MsixPackageApplicationModel, 0)
	if inputs == nil {
		return outputs
	}

	for _, input := range pointer.From(inputs) {
		outputs = append(outputs, MsixPackageApplicationModel{
			AppId:          pointer.From(input.AppId),
			AppUserModelID: pointer.From(input.AppUserModelID),
			Description:    pointer.From(input.Description),
			FriendlyName:   pointer.From(input.FriendlyName),
			IconImageName:  pointer.From(input.IconImageName),
			RawIcon:        pointer.From(input.RawIcon),
			RawPng:         pointer.From(input.RawPng),
		})
	}

	return outputs
}

func getMsixImageProperties(ctx context.Context, metadata sdk.ResourceMetaData, hostPoolIds []string, storageShareFileId string, msixPackageName string) (*msiximage.ExpandMsixImageProperties, error) {
	msixImageUri := msiximage.MSIXImageURI{
		Uri: pointer.To(storageShareFileId),
	}

	availableMsixPackageNames := make(map[string][]string)

	for _, hostPoolReference := range hostPoolIds {
		availableMsixPackageNames[hostPoolReference] = make([]string, 0)
		hostPoolId, err := hostpool.ParseHostPoolIDInsensitively(hostPoolReference)
		if err != nil {
			return nil, fmt.Errorf("parsing host pool ID %s: %+v", hostPoolReference, err)
		}

		msixImageHostPoolId := msiximage.NewHostPoolID(hostPoolId.SubscriptionId, hostPoolId.ResourceGroupName, hostPoolId.HostPoolName)

		// Replicate portal behavior where MSIX image properties required for app attach package creation are retrieved with API
		// `Expand` method from imported `msiximage` package is not used as it lacks codes to marshal request body
		result, err := method.ExpandCompleteMsixImage(ctx, metadata, msixImageHostPoolId, msixImageUri)
		if err != nil {
			// Continue to check next host pool if there is error with expanding MSIX image of current host pool
			continue
		}

		msixImages := result.Items
		for _, msixImage := range msixImages {
			if properties := msixImage.Properties; properties != nil {
				// Return MSIX image with matched `PackageFullName`
				if properties.PackageFullName != nil && strings.EqualFold(pointer.From(properties.PackageFullName), msixPackageName) {
					return properties, nil
				}
			}
		}

		// Store available `PackageFullName` for user reference if no matched MSIX image is found after checking all host pools
		for _, msixImage := range msixImages {
			if properties := msixImage.Properties; properties != nil && properties.PackageFullName != nil {
				availableMsixPackageNames[hostPoolReference] = append(availableMsixPackageNames[hostPoolReference], pointer.From(properties.PackageFullName))
			}
		}
	}

	concatenatedAvailableMsixPackageNames := ""
	for hostPoolReference, packageFullNames := range availableMsixPackageNames {
		concatenatedAvailableMsixPackageNames += fmt.Sprintf("%v from %s, ", packageFullNames, hostPoolReference)
	}

	concatenatedAvailableMsixPackageNames = strings.TrimSuffix(concatenatedAvailableMsixPackageNames, ", ")

	return nil, fmt.Errorf("no matched MSIX image with package name %s was found. The available package names are %s", msixPackageName, concatenatedAvailableMsixPackageNames)
}
