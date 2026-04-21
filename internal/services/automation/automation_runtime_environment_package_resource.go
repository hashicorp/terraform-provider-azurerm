// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/packageresource"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = AutomationRuntimeEnvironmentPackageResource{}

type AutomationRuntimeEnvironmentPackageModel struct {
	Name                           string `tfschema:"name"`
	AutomationRuntimeEnvironmentId string `tfschema:"automation_runtime_environment_id"`
	ContentUri                     string `tfschema:"content_uri"`
	ContentVersion                 string `tfschema:"content_version"`
	HashAlgorithm                  string `tfschema:"hash_algorithm"`
	HashValue                      string `tfschema:"hash_value"`
	SizeInBytes                    int64  `tfschema:"size_in_bytes"`
	Version                        string `tfschema:"version"`
	Default                        bool   `tfschema:"default"`
}

type AutomationRuntimeEnvironmentPackageResource struct{}

func (r AutomationRuntimeEnvironmentPackageResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"automation_runtime_environment_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: packageresource.ValidateRuntimeEnvironmentID,
		},

		"content_uri": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsURLWithHTTPS,
		},

		"content_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[0-9]+\.[0-9]+(\.[0-9]+){0,2}$`),
				"The `content_version` must be with 2 to 4 segments (e.g. `1.0`, `1.0.0`, or `1.0.0.0`).",
			),
		},

		"hash_algorithm": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{"hash_value"},
		},

		"hash_value": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{"hash_algorithm"},
		},
	}
}

func (r AutomationRuntimeEnvironmentPackageResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"default": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"size_in_bytes": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r AutomationRuntimeEnvironmentPackageResource) ModelObject() interface{} {
	return &AutomationRuntimeEnvironmentPackageModel{}
}

func (r AutomationRuntimeEnvironmentPackageResource) ResourceType() string {
	return "azurerm_automation_runtime_environment_package"
}

func (r AutomationRuntimeEnvironmentPackageResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return packageresource.ValidatePackageID
}

func (r AutomationRuntimeEnvironmentPackageResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automation.PackageResource

			var model AutomationRuntimeEnvironmentPackageModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			envId, err := packageresource.ParseRuntimeEnvironmentID(model.AutomationRuntimeEnvironmentId)
			if err != nil {
				return err
			}

			id := packageresource.NewPackageID(envId.SubscriptionId, envId.ResourceGroupName, envId.AutomationAccountName, envId.RuntimeEnvironmentName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := packageresource.PackageCreateOrUpdateParameters{
				Properties: packageresource.PackageCreateOrUpdateProperties{
					ContentLink: packageresource.ContentLink{
						Uri: pointer.To(model.ContentUri),
					},
				},
			}

			if model.ContentVersion != "" {
				parameters.Properties.ContentLink.Version = pointer.To(model.ContentVersion)
			}

			if model.HashAlgorithm != "" {
				parameters.Properties.ContentLink.ContentHash = &packageresource.ContentHash{
					Algorithm: model.HashAlgorithm,
					Value:     model.HashValue,
				}
			}

			if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			// custom poller is required until https://github.com/Azure/azure-rest-api-specs/issues/41641 is resolved
			pollerType := custompollers.NewAutomationRuntimeEnvironmentPackagePoller(client, id)
			poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
			if err := poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for %s to finish provisioning: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r AutomationRuntimeEnvironmentPackageResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automation.PackageResource

			id, err := packageresource.ParsePackageID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			var stateModel AutomationRuntimeEnvironmentPackageModel
			if err = metadata.Decode(&stateModel); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			runtimeEnvId := packageresource.NewRuntimeEnvironmentID(id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName, id.RuntimeEnvironmentName)

			output := AutomationRuntimeEnvironmentPackageModel{
				Name:                           id.PackageName,
				AutomationRuntimeEnvironmentId: runtimeEnvId.ID(),

				// the fields below don't return by the API, remove it when issue fixed
				// https://github.com/Azure/azure-rest-api-specs/issues/41604
				ContentVersion: stateModel.ContentVersion,
				ContentUri:     stateModel.ContentUri,
				HashValue:      stateModel.HashValue,
				HashAlgorithm:  stateModel.HashAlgorithm,
			}

			if props := resp.Model.Properties; props != nil {
				if content := props.ContentLink; content != nil {
					output.ContentUri = pointer.From(content.Uri)
					output.ContentVersion = pointer.From(content.Version)

					if hash := content.ContentHash; hash != nil {
						output.HashAlgorithm = hash.Algorithm
						output.HashValue = hash.Value
					}
				}

				output.SizeInBytes = pointer.From(props.SizeInBytes)
				output.Version = pointer.From(props.Version)
				output.Default = pointer.From(props.Default)
			}

			return metadata.Encode(&output)
		},
	}
}

func (r AutomationRuntimeEnvironmentPackageResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automation.PackageResource

			id, err := packageresource.ParsePackageID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Delete(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return nil
				}
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
