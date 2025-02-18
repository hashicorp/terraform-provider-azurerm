package desktopvirtualization

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2022-02-10-preview/msixpackage"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = VirtualDesktopMSIXPackageResource{}

type VirtualDesktopMSIXPackageResource struct{}

type PackageApplicationModel struct {
	AppId          string `tfschema:"app_id"`
	AppUserModelId string `tfschema:"app_user_model_id"`
	Description    string `tfschema:"description"`
	FriendlyName   string `tfschema:"friendly_name"`
	IconImageName  string `tfschema:"icon_image_name"`
	RawIcon        string `tfschema:"raw_icon"`
	RawPng         string `tfschema:"raw_png"`
}

type PackageDependencyModel struct {
	DependencyName string `tfschema:"dependency_name"`
	MinVersion     string `tfschema:"min_version"`
	Publisher      string `tfschema:"publisher"`
}

type VirtualDesktopMSIXPackageResourceModel struct {
	// ID fields

	Name              string `tfschema:"name"`
	HostPoolName      string `tfschema:"host_pool_name"`
	ResourceGroupName string `tfschema:"resource_group_name"`

	// Required fields

	ImagePath           string                    `tfschema:"image_path"`
	LastUpdatedInUTC    string                    `tfschema:"last_updated_in_utc"`
	PackageApplications []PackageApplicationModel `tfschema:"package_application"`
	PackageFamilyName   string                    `tfschema:"package_family_name"`
	PackageName         string                    `tfschema:"package_name"`
	PackageRelativePath string                    `tfschema:"package_relative_path"`
	Version             string                    `tfschema:"version"`

	// Optional fields

	DisplayName                string                   `tfschema:"display_name"`
	Enabled                    bool                     `tfschema:"enabled"`
	PackageDependencies        []PackageDependencyModel `tfschema:"package_dependency"`
	RegularRegistrationEnabled bool                     `tfschema:"regular_registration_enabled"`
}

func requiredNotEmptyStringForceNewSchema() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validation.StringIsNotEmpty,
	}
}

func (VirtualDesktopMSIXPackageResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": requiredNotEmptyStringForceNewSchema(),

		"host_pool_name": requiredNotEmptyStringForceNewSchema(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"image_path": requiredNotEmptyStringForceNewSchema(),

		"last_updated_in_utc": requiredNotEmptyStringForceNewSchema(),

		"package_application": {
			Type:     pluginsdk.TypeList,
			MinItems: 1,
			Required: true,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"app_id": requiredNotEmptyStringForceNewSchema(),

					"app_user_model_id": requiredNotEmptyStringForceNewSchema(),

					"description": requiredNotEmptyStringForceNewSchema(),

					"friendly_name": requiredNotEmptyStringForceNewSchema(),

					"icon_image_name": requiredNotEmptyStringForceNewSchema(),

					"raw_icon": requiredNotEmptyStringForceNewSchema(),

					"raw_png": requiredNotEmptyStringForceNewSchema(),
				},
			},
		},

		"package_family_name": requiredNotEmptyStringForceNewSchema(),

		"package_name": requiredNotEmptyStringForceNewSchema(),

		"package_relative_path": requiredNotEmptyStringForceNewSchema(),

		"version": requiredNotEmptyStringForceNewSchema(),

		"display_name": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},

		"package_dependency": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"dependency_name": requiredNotEmptyStringForceNewSchema(),

					"min_version": requiredNotEmptyStringForceNewSchema(),

					"publisher": requiredNotEmptyStringForceNewSchema(),
				},
			},
		},

		"regular_registration_enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
	}
}

func (VirtualDesktopMSIXPackageResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (VirtualDesktopMSIXPackageResource) ModelObject() interface{} {
	return &VirtualDesktopMSIXPackageResourceModel{}
}

func (VirtualDesktopMSIXPackageResource) ResourceType() string {
	return "azurerm_virtual_desktop_msix_package"
}

func (r VirtualDesktopMSIXPackageResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.MSIXPackageClient

			subscriptionId := metadata.Client.Account.SubscriptionId

			var config VirtualDesktopMSIXPackageResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := msixpackage.NewMsixPackageID(subscriptionId, config.ResourceGroupName, config.HostPoolName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := msixpackage.MSIXPackage{
				Properties: msixpackage.MSIXPackageProperties{
					DisplayName:           pointer.To(config.DisplayName),
					ImagePath:             pointer.To(config.ImagePath),
					IsActive:              pointer.To(config.Enabled),
					IsRegularRegistration: pointer.To(config.RegularRegistrationEnabled),
					LastUpdated:           pointer.To(config.LastUpdatedInUTC),
					PackageApplications:   expandPackageApplications(config.PackageApplications),
					PackageDependencies:   expandPackageDependencies(config.PackageDependencies),
					PackageFamilyName:     pointer.To(config.PackageFamilyName),
					PackageName:           pointer.To(config.PackageName),
					PackageRelativePath:   pointer.To(config.PackageRelativePath),
					Version:               pointer.To(config.Version),
				},
			}

			if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func expandPackageApplications(packageApplicationModel []PackageApplicationModel) *[]msixpackage.MsixPackageApplications {
	packageApplications := make([]msixpackage.MsixPackageApplications, 0)
	for _, model := range packageApplicationModel {
		packageApplications = append(packageApplications, msixpackage.MsixPackageApplications{
			AppId:          pointer.To(model.AppId),
			AppUserModelID: pointer.To(model.AppUserModelId),
			Description:    pointer.To(model.Description),
			FriendlyName:   pointer.To(model.FriendlyName),
			IconImageName:  pointer.To(model.IconImageName),
			RawIcon:        pointer.To(model.RawIcon),
			RawPng:         pointer.To(model.RawPng),
		})
	}
	return &packageApplications
}

func expandPackageDependencies(packageDependencyModel []PackageDependencyModel) *[]msixpackage.MsixPackageDependencies {
	packageDependencies := make([]msixpackage.MsixPackageDependencies, 0)
	for _, model := range packageDependencyModel {
		packageDependencies = append(packageDependencies, msixpackage.MsixPackageDependencies{
			DependencyName: pointer.To(model.DependencyName),
			MinVersion:     pointer.To(model.MinVersion),
			Publisher:      pointer.To(model.Publisher),
		})
	}
	return &packageDependencies
}

func (r VirtualDesktopMSIXPackageResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.MSIXPackageClient

			id, err := msixpackage.ParseMsixPackageID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config VirtualDesktopMSIXPackageResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			payload := existing.Model

			if metadata.ResourceData.HasChange("display_name") {
				payload.Properties.DisplayName = pointer.To(config.DisplayName)
			}

			if metadata.ResourceData.HasChange("enabled") {
				payload.Properties.IsActive = pointer.To(config.Enabled)
			}

			if metadata.ResourceData.HasChange("regular_registration_enabled") {
				payload.Properties.IsRegularRegistration = pointer.To(config.RegularRegistrationEnabled)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (VirtualDesktopMSIXPackageResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.MSIXPackageClient

			id, err := msixpackage.ParseMsixPackageID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			props := resp.Model.Properties

			state := VirtualDesktopMSIXPackageResourceModel{
				Name:                       id.MsixPackageName,
				HostPoolName:               id.HostPoolName,
				ResourceGroupName:          id.ResourceGroupName,
				ImagePath:                  pointer.From(props.ImagePath),
				LastUpdatedInUTC:           pointer.From(props.LastUpdated),
				PackageApplications:        *flattenPackageApplications(props.PackageApplications),
				PackageFamilyName:          pointer.From(props.PackageFamilyName),
				PackageName:                pointer.From(props.PackageName),
				PackageRelativePath:        pointer.From(props.PackageRelativePath),
				Version:                    pointer.From(props.Version),
				DisplayName:                pointer.From(props.DisplayName),
				Enabled:                    pointer.From(props.IsActive),
				PackageDependencies:        *flattenPackageDependencies(props.PackageDependencies),
				RegularRegistrationEnabled: pointer.From(props.IsRegularRegistration),
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenPackageApplications(msixPackageApplications *[]msixpackage.MsixPackageApplications) *[]PackageApplicationModel {
	models := make([]PackageApplicationModel, 0)
	for _, msixPackageApplication := range pointer.From(msixPackageApplications) {
		models = append(models, PackageApplicationModel{
			AppId:          pointer.From(msixPackageApplication.AppId),
			AppUserModelId: pointer.From(msixPackageApplication.AppUserModelID),
			Description:    pointer.From(msixPackageApplication.Description),
			FriendlyName:   pointer.From(msixPackageApplication.FriendlyName),
			IconImageName:  pointer.From(msixPackageApplication.IconImageName),
			RawIcon:        pointer.From(msixPackageApplication.RawIcon),
			RawPng:         pointer.From(msixPackageApplication.RawPng),
		})
	}
	return &models
}

func flattenPackageDependencies(msixPackageDependencies *[]msixpackage.MsixPackageDependencies) *[]PackageDependencyModel {
	models := make([]PackageDependencyModel, 0)
	for _, msixPackageDependency := range pointer.From(msixPackageDependencies) {
		models = append(models, PackageDependencyModel{
			DependencyName: pointer.From(msixPackageDependency.DependencyName),
			MinVersion:     pointer.From(msixPackageDependency.MinVersion),
			Publisher:      pointer.From(msixPackageDependency.Publisher),
		})
	}
	return &models
}

func (VirtualDesktopMSIXPackageResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.MSIXPackageClient

			id, err := msixpackage.ParseMsixPackageID(metadata.ResourceData.Id())
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

func (VirtualDesktopMSIXPackageResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return msixpackage.ValidateMsixPackageID
}
